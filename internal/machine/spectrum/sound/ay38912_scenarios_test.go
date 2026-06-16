package sound

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var ay38912Scenarios = []dsl.Scenario{
	dsl.NewScenario("AY38912 Register Write/Read", func(t *testing.T) {
		ay := NewAY38912()
		ay.WriteAddress(0)
		ay.WriteData(0x44)
		assert.Equal(t, "Register 0 value", ay.ReadData(), uint8(0x44))
		
		ay.WriteAddress(1)
		ay.WriteData(0x05)
		assert.Equal(t, "Channel A Period", ay.Channels[0].Period, uint16(0x0544))
	}),
	dsl.NewScenario("AY38912 Tone Generation", func(t *testing.T) {
		ay := NewAY38912()
		// Enable Tone A only (R7 bit 0 = 0)
		ay.WriteAddress(7)
		ay.WriteData(0xFE) // 1111 1110
		
		// Set Period A to 2 (R0=2, R1=0)
		ay.WriteAddress(0)
		ay.WriteData(2)
		ay.WriteAddress(1)
		ay.WriteData(0)
		
		// Set Amplitude A to 15 (R8=15)
		ay.WriteAddress(8)
		ay.WriteData(15)
		
		// Initial state
		assert.True(t, "Channel A Output initial", !ay.Channels[0].Output)
		
		// AY internal divider for tones is 8.
		// For Period 2, it flips every 2 * 8 = 16 ticks.
		
		// Tick 1 to 15
		for i := 0; i < 15; i++ {
			ay.Tick(false)
		}
		assert.True(t, "Channel A Output after 15 ticks", !ay.Channels[0].Output)
		
		// Tick 16
		ay.Tick(false)
		assert.True(t, "Channel A Output after 16 ticks", ay.Channels[0].Output)
		
		// Tick 17 to 31
		for i := 0; i < 15; i++ {
			ay.Tick(false)
		}
		assert.True(t, "Channel A Output after 31 ticks", ay.Channels[0].Output)
		
		// Tick 32
		ay.Tick(false)
		assert.True(t, "Channel A Output after 32 ticks", !ay.Channels[0].Output)
	}),
	dsl.NewScenario("AY38912 Envelope Step", func(t *testing.T) {
		ay := NewAY38912()
		// Set Envelope Period to 1 (R11=1, R12=0)
		ay.WriteAddress(11)
		ay.WriteData(1)
		ay.WriteAddress(12)
		ay.WriteData(0)
		
		// Set Shape to Single Decay
		ay.WriteAddress(13)
		ay.WriteData(0)
		
		assert.Equal(t, "Initial Envelope Step", ay.Envelope.Step, int8(31))
		
		// Envelope steps every Period * 8 ticks. For Period 1, it steps every 8 ticks.
		
		for i := 0; i < 7; i++ {
			ay.Tick(false)
		}
		assert.Equal(t, "Envelope Step after 7 ticks", ay.Envelope.Step, int8(31))

		ay.Tick(false)
		assert.Equal(t, "Envelope Step after 8 ticks", ay.Envelope.Step, int8(30))
		
		for i := 0; i < 30*8; i++ {
			ay.Tick(false)
		}
		assert.Equal(t, "Envelope Step after 31 full steps", ay.Envelope.Step, int8(0))
		assert.True(t, "Envelope NOT Done yet", !ay.Envelope.Done)
		
		for i := 0; i < 8; i++ {
			ay.Tick(false)
		}
		assert.True(t, "Envelope Done after 32 full steps", ay.Envelope.Done)
	}),
	dsl.NewScenario("AY38912 Stereo Panning", func(t *testing.T) {
		ay := NewAY38912()
		// Enable Tones only
		ay.WriteAddress(7)
		ay.WriteData(0xF8) // 1111 1000 (A, B, C enabled)
		
		// Ensure outputs are high for mixing test
		for i := 0; i < 3; i++ {
			ay.Channels[i].Output = true
		}

		// ACB Panning (Standard Spectrum 128k):
		// L = (2*A + C) / 3
		// R = (2*B + C) / 3
		
		// Only A on (A=15, B=0, C=0)
		ay.WriteAddress(8); ay.WriteData(15)
		ay.WriteAddress(9); ay.WriteData(0)
		ay.WriteAddress(10); ay.WriteData(0)
		l, r := ay.Mix()
		assert.True(t, "Only A: Left > Right", l > r)
		assert.Equal(t, "Only A: Right is zero", r, uint16(0))
		
		// Only B on (A=0, B=15, C=0)
		ay.WriteAddress(8); ay.WriteData(0)
		ay.WriteAddress(9); ay.WriteData(15)
		ay.WriteAddress(10); ay.WriteData(0)
		l, r = ay.Mix()
		assert.True(t, "Only B: Right > Left", r > l)
		assert.Equal(t, "Only B: Left is zero", l, uint16(0))
		
		// Only C on (A=0, B=0, C=15)
		ay.WriteAddress(8); ay.WriteData(0)
		ay.WriteAddress(9); ay.WriteData(0)
		ay.WriteAddress(10); ay.WriteData(15)
		l, r = ay.Mix()
		assert.Equal(t, "Only C: Left == Right", l, r)
		assert.True(t, "Only C: Both > 0", l > 0)
	}),
}
