package sound

import (
	"log/slog"
	"sync"
)

// SampleRate is the target audio sampling rate.
const SampleRate = 44100

// AYClockRate is the clock rate of the AY chip in a Spectrum 128K.
const AYClockRate = 1773447

// volumeTable provides the 32-level logarithmic volume mapping for the AY-3-8912.
// These values are based on the ESPectrum dynamic logarithmic model (Hacker KAY).
var volumeTable = []uint16{
	0, 5, 8, 11, 16, 22, 32, 45,
	64, 90, 127, 180, 255, 360, 510, 720,
	1020, 1440, 2040, 2880, 4080, 5760, 8160, 11520,
	16320, 23040, 32640, 46080, 65280, 65535, 65535, 65535,
}

// rampaAYTable maps 4-bit amplitude registers (0-15) to 32-level volume table indices.
var rampaAYTable = []uint8{
	0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30,
}

// Channel represents one of the three sound channels (A, B, C).
type Channel struct {
	Period    uint16
	Counter   uint16
	Output    bool
	Amplitude uint8
	Envelope  bool
}

// EnvelopeGenerator manages the automated volume changes.
type EnvelopeGenerator struct {
	Period  uint32
	Counter uint32
	Step    int8
	Shape   uint8
	Hold    bool
	Repeat  bool
	Attack  bool
	Alt     bool
	Done    bool
}

// NoiseGenerator manages the pseudo-random noise source.
type NoiseGenerator struct {
	Period  uint8
	Counter uint8
	Seed    uint32
	Output  bool
}

// AY38912 represents the AY-3-8912 Programmable Sound Generator.
type AY38912 struct {
	Registers        [16]uint8
	SelectedRegister uint8

	Channels [3]Channel
	Noise    NoiseGenerator
	Envelope EnvelopeGenerator

	// Internal clock divider
	masterCounter uint32

	// Accumulation for sampling (Stereo)
	sampleSumL  int32
	sampleSumR  int32
	sampleCount uint32

	// Buffer for Ebitengine (interleaved L/R int16)
	audioBuffer []byte
	mutex       sync.Mutex
}

// NewAY38912 creates a new AY-3-8912 instance.
func NewAY38912() *AY38912 {
	ay := &AY38912{}
	ay.Noise.Seed = 1
	ay.Noise.Output = true
	return ay
}

// WriteAddress selects the register for subsequent read/write operations.
func (ay *AY38912) WriteAddress(val uint8) {
	ay.SelectedRegister = val & 0x0F
}

// WriteData writes a value to the currently selected register.
func (ay *AY38912) WriteData(val uint8) {
	reg := ay.SelectedRegister
	ay.Registers[reg] = val

	switch reg {
	case 0, 1: // Channel A Period
		ay.Channels[0].Period = (uint16(ay.Registers[1]&0x0F) << 8) | uint16(ay.Registers[0])
	case 2, 3: // Channel B Period
		ay.Channels[1].Period = (uint16(ay.Registers[3]&0x0F) << 8) | uint16(ay.Registers[2])
	case 4, 5: // Channel C Period
		ay.Channels[2].Period = (uint16(ay.Registers[5]&0x0F) << 8) | uint16(ay.Registers[4])
	case 6: // Noise Period
		ay.Noise.Period = val & 0x1F
	case 8, 9, 10: // Amplitude
		ch := reg - 8
		ay.Channels[ch].Amplitude = val & 0x0F
		ay.Channels[ch].Envelope = (val & 0x10) != 0
	case 11, 12: // Envelope Period
		ay.Envelope.Period = (uint32(ay.Registers[12]) << 8) | uint32(ay.Registers[11])
	case 13: // Envelope Shape
		ay.Envelope.Shape = val & 0x0F
		ay.Envelope.Counter = 0
		ay.Envelope.Done = false
		// Bits: Continue, Attack, Alternate, Hold
		cont := (val & 0x08) != 0
		ay.Envelope.Attack = (val & 0x04) != 0
		ay.Envelope.Alt = (val & 0x02) != 0
		ay.Envelope.Hold = (val & 0x01) != 0

		if !cont {
			ay.Envelope.Repeat = false
			if ay.Envelope.Attack {
				ay.Envelope.Alt = false
			} else {
				ay.Envelope.Alt = true
			}
		} else {
			ay.Envelope.Repeat = !ay.Envelope.Hold
		}

		if ay.Envelope.Attack {
			ay.Envelope.Step = 0
		} else {
			ay.Envelope.Step = 31
		}
	}

	slog.Debug("AY-3-8912 register write", "reg", reg, "val", val)
}

// ReadData returns the value of the currently selected register.
func (ay *AY38912) ReadData() uint8 {
	return ay.Registers[ay.SelectedRegister]
}

// Tick advances the internal state of the AY chip by one AY clock cycle (1.77MHz).
// beeper: state of the ULA beeper (Port 0xFE bit 4)
func (ay *AY38912) Tick(beeper bool) {
	ay.masterCounter++

	// 1. Tones (1.77MHz / 8 = 221kHz)
	if ay.masterCounter&0x07 == 0 {
		for i := 0; i < 3; i++ {
			ch := &ay.Channels[i]
			if ch.Period > 0 {
				ch.Counter++
				if ch.Counter >= ch.Period {
					ch.Counter = 0
					ch.Output = !ch.Output
				}
			} else {
				ch.Output = true
			}
		}
	}

	// 2. Noise (1.77MHz / 16 = 110kHz)
	if ay.masterCounter&0x0F == 0 {
		ay.Noise.Counter++
		noisePeriod := ay.Noise.Period
		if noisePeriod == 0 {
			noisePeriod = 1
		}
		if ay.Noise.Counter >= noisePeriod {
			ay.Noise.Counter = 0
			bit16 := (ay.Noise.Seed >> 16) & 1
			bit13 := (ay.Noise.Seed >> 13) & 1
			ay.Noise.Seed = ((ay.Noise.Seed*2 + 1) ^ (bit16 ^ bit13)) & 0x1FFFF
			ay.Noise.Output = (ay.Noise.Seed & 1) != 0
		}
	}

	// 3. Envelope (1.77MHz / 8 = 221kHz)
	if ay.masterCounter&0x07 == 0 && !ay.Envelope.Done {
		ay.Envelope.Counter++
		envPeriod := ay.Envelope.Period
		if envPeriod == 0 {
			envPeriod = 1
		}
		if ay.Envelope.Counter >= uint32(envPeriod) {
			ay.Envelope.Counter = 0

			if ay.Envelope.Attack {
				ay.Envelope.Step++
				if ay.Envelope.Step > 31 {
					if ay.Envelope.Repeat {
						if ay.Envelope.Alt {
							ay.Envelope.Attack = false
							ay.Envelope.Step = 31
						} else {
							ay.Envelope.Step = 0
						}
					} else {
						if ay.Envelope.Hold {
							if ay.Envelope.Alt {
								ay.Envelope.Step = 0
							} else {
								ay.Envelope.Step = 31
							}
						} else {
							ay.Envelope.Step = 0
						}
						ay.Envelope.Done = true
					}
				}
			} else {
				ay.Envelope.Step--
				if ay.Envelope.Step < 0 {
					if ay.Envelope.Repeat {
						if ay.Envelope.Alt {
							ay.Envelope.Attack = true
							ay.Envelope.Step = 0
						} else {
							ay.Envelope.Step = 31
						}
					} else {
						if ay.Envelope.Hold {
							if ay.Envelope.Alt {
								ay.Envelope.Step = 31
							} else {
								ay.Envelope.Step = 0
							}
						} else {
							ay.Envelope.Step = 0
						}
						ay.Envelope.Done = true
					}
				}
			}
		}
	}

	// 4. Mix channels
	ml, mr := ay.Mix()

	// 5. Add Beeper (ULA Sound)
	var l int32 = int32(ml)
	var r int32 = int32(mr)
	if beeper {
		// Beeper value. Max AY is 65535 per channel.
		// Total max can be (2*A + C) = 3*65535 = 196605.
		// We use a high beeper volume (~1/4 of total max).
		beeperVal := int32(49152)
		l += beeperVal
		r += beeperVal
	}

	ay.sampleSumL += l
	ay.sampleSumR += r
	ay.sampleCount++
}

// Mix combines the outputs of the three channels with ACB stereo panning.
func (ay *AY38912) Mix() (uint16, uint16) {
	enableReg := ay.Registers[7]
	var v [3]uint32

	for i := 0; i < 3; i++ {
		ch := ay.Channels[i]

		// Tone and Noise enable bits (active low in register)
		toneEnabled := (enableReg & (1 << i)) == 0
		noiseEnabled := (enableReg & (8 << i)) == 0

		out := true
		if toneEnabled {
			out = out && ch.Output
		}
		if noiseEnabled {
			out = out && ay.Noise.Output
		}

		if out {
			volIdx := uint8(0)
			if ch.Envelope {
				volIdx = uint8(ay.Envelope.Step)
			} else {
				volIdx = rampaAYTable[ch.Amplitude&0x0F]
			}
			v[i] = uint32(volumeTable[volIdx])
		} else {
			v[i] = 0
		}
	}

	// ACB Panning (Standard Spectrum 128k):
	// L = 2*A + C, R = 2*B + C
	// Max value for each is 3 * 65535 = 196605.
	l := v[0]*2 + v[2]
	r := v[1]*2 + v[2]

	// Normalize to uint16 range (0-65535)
	return uint16(l / 3), uint16(r / 3)
}

// GetSample returns the current averaged stereo samples and resets the accumulators.
func (ay *AY38912) GetSample() (int16, int16) {
	if ay.sampleCount == 0 {
		return 0, 0
	}
	avgL := ay.sampleSumL / int32(ay.sampleCount)
	avgR := ay.sampleSumR / int32(ay.sampleCount)
	ay.sampleSumL = 0
	ay.sampleSumR = 0
	ay.sampleCount = 0

	// Centering and Scaling
	// Range is 0 to ~114687 (65535 AY + 49152 Beeper).
	// We center it around half-max (~57344) and scale down to fit int16.
	// result = (avg - 57344) * 4 / 7
	resL := (avgL - 57344) * 4 / 7
	resR := (avgR - 57344) * 4 / 7

	// Final clamp
	if resL > 32767 { resL = 32767 }
	if resL < -32768 { resL = -32768 }
	if resR > 32767 { resR = 32767 }
	if resR < -32768 { resR = -32768 }

	return int16(resL), int16(resR)
}

// Read implements io.Reader to provide audio samples.
func (ay *AY38912) Read(p []byte) (n int, err error) {
	ay.mutex.Lock()
	defer ay.mutex.Unlock()

	toCopy := len(p)
	if toCopy > len(ay.audioBuffer) {
		toCopy = len(ay.audioBuffer)
	}

	if toCopy > 0 {
		copy(p, ay.audioBuffer[:toCopy])
		ay.audioBuffer = ay.audioBuffer[toCopy:]
	}

	return toCopy, nil
}

// AddAudioSample adds a new stereo sample to the internal audio buffer.
func (ay *AY38912) AddAudioSample(left, right int16) {
	ay.mutex.Lock()
	defer ay.mutex.Unlock()

	ay.audioBuffer = append(ay.audioBuffer, byte(left), byte(left>>8))
	ay.audioBuffer = append(ay.audioBuffer, byte(right), byte(right>>8))

	if len(ay.audioBuffer) > 16384 {
		ay.audioBuffer = ay.audioBuffer[len(ay.audioBuffer)-16384:]
	}
}
