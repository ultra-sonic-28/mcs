package spectrum

import (
	"log/slog"
	spectrumrom "mcs/assets/machines/spectrum"
)

// Bus128 implements the memory and I/O bus for the ZX Spectrum 128K.
type Bus128 struct {
	BaseBus
	rom [2][16384]uint8
	ram [8][16384]uint8

	currentRom   uint8
	currentRam   uint8
	pagingLocked bool
	activeScreen uint8 // 5 or 7

	AY *AY38912
}

// NewBus128 creates a new Spectrum 128K Bus with the embedded ROMs loaded.
func NewBus128() *Bus128 {
	slog.Info("Initializing Spectrum 128K Bus")
	b := &Bus128{
		BaseBus: BaseBus{
			Keyboard: NewKeyboard(),
			Display:  NewDisplay(),
			Tape:     NewTape(),
		},
		activeScreen: 5,
		AY:           NewAY38912(),
	}

	// Load ROMs
	if len(spectrumrom.Rom128_0) != 16384 {
		slog.Warn("Spectrum 128K ROM 1 size is unexpected", "expected", 16384, "actual", len(spectrumrom.Rom128_0))
	}
	if len(spectrumrom.Rom128_1) != 16384 {
		slog.Warn("Spectrum 128K ROM 2 size is unexpected", "expected", 16384, "actual", len(spectrumrom.Rom128_1))
	}
	copy(b.rom[0][:], spectrumrom.Rom128_0)
	copy(b.rom[1][:], spectrumrom.Rom128_1)

	return b
}

// Read returns the byte at the specified memory address.
func (b *Bus128) Read(addr uint16) uint8 {
	switch {
	case addr < 0x4000:
		return b.rom[b.currentRom][addr]
	case addr < 0x8000:
		return b.ram[5][addr-0x4000]
	case addr < 0xC000:
		return b.ram[2][addr-0x8000]
	default:
		return b.ram[b.currentRam][addr-0xC000]
	}
}

// Read16 returns the 16-bit word at the specified memory address (little-endian).
func (b *Bus128) Read16(addr uint16) uint16 {
	return uint16(b.Read(addr)) | (uint16(b.Read(addr+1)) << 8)
}

// Write stores a byte at the specified memory address.
func (b *Bus128) Write(addr uint16, val uint8) {
	switch {
	case addr < 0x4000:
		// ROM is read-only
		return
	case addr < 0x8000:
		b.ram[5][addr-0x4000] = val
	case addr < 0xC000:
		b.ram[2][addr-0x8000] = val
	default:
		b.ram[b.currentRam][addr-0xC000] = val
	}
}

// In reads a byte from the specified I/O port.
func (b *Bus128) In(port uint16) uint8 {
	// Port 0xFE: Keyboard and Tape
	if port&0x01 == 0 {
		mask := uint8(port >> 8)
		result := b.Keyboard.Scan(mask)

		if b.TapeInState {
			result |= 0x40
		} else {
			result &= ^uint8(0x40)
		}

		result |= 0xA0
		return result
	}

	// Port 0xFFFD: AY Register Read (bit 1=0, bit 14=1, bit 15=1)
	if (port & 0xC002) == 0xC000 {
		return b.AY.ReadData()
	}

	return 0xFF
}

// Out writes a byte to the specified I/O port.
func (b *Bus128) Out(port uint16, val uint8) {
	// Port 0xFE: Border, Beeper, Mic
	if port&0x01 == 0 {
		b.BorderColor = val & 0x07
		b.MicState = (val & 0x08) != 0
		b.BeeperState = (val & 0x10) != 0
		return
	}

	// Port 0x7FFD: Memory Paging (bit 1=0, bit 15=0)
	if (port & 0x8002) == 0 {
		if b.pagingLocked {
			return
		}

		b.currentRam = val & 0x07
		if (val & 0x08) != 0 {
			b.activeScreen = 7
		} else {
			b.activeScreen = 5
		}
		if (val & 0x10) != 0 {
			b.currentRom = 1
		} else {
			b.currentRom = 0
		}
		if (val & 0x20) != 0 {
			b.pagingLocked = true
		}

		slog.Debug("Spectrum 128K Paging", "ram", b.currentRam, "rom", b.currentRom, "screen", b.activeScreen, "locked", b.pagingLocked)
		return
	}

	// Port 0xFFFD: AY Register Select (bit 1=0, bit 14=1, bit 15=1)
	if (port & 0xC002) == 0xC000 {
		b.AY.WriteAddress(val)
		return
	}

	// Port 0xBFFD: AY Data Write (bit 1=0, bit 14=0, bit 15=1)
	if (port & 0xC002) == 0x8000 {
		b.AY.WriteData(val)
		return
	}
}

// GetDisplayMemory returns the 6912 bytes of memory used for the display.
func (b *Bus128) GetDisplayMemory() []byte {
	return b.ram[b.activeScreen][0:6912]
}

func (b *Bus128) IsRom1Active() bool {
	return b.currentRom == 1
}
