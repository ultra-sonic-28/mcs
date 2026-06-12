package bus

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var baseBusScenarios = []dsl.Scenario{
	dsl.NewScenario("BaseBus component accessors", func(t *testing.T) {
		bus := NewBus48()
		
		assert.True(t, "Keyboard is not nil", bus.GetKeyboard() != nil)
		assert.True(t, "Tape is not nil", bus.GetTape() != nil)
		assert.True(t, "Display is not nil", bus.GetDisplay() != nil)
	}),
	dsl.NewScenario("BaseBus ULA state management", func(t *testing.T) {
		bus := NewBus48()
		
		// Border color (via Out port 0xFE)
		bus.Out(0x00FE, 0x01) // Blue
		assert.Equal(t, "Border color is Blue", bus.GetBorderColor(), uint8(1))
		
		bus.Out(0x00FE, 0x02) // Red
		assert.Equal(t, "Border color is Red", bus.GetBorderColor(), uint8(2))
		
		// Tape In state
		bus.SetTapeInState(true)
		assert.True(t, "Tape in state is true", bus.GetTapeInState())
		
		bus.SetTapeInState(false)
		assert.True(t, "Tape in state is false", !bus.GetTapeInState())
	}),
	dsl.NewScenario("ULA Out Port 0xFE details", func(t *testing.T) {
		bus := NewBus48()
		// Border 2 (Red), MIC on, Beeper off
		bus.Out(0x00FE, 0x0A) // 0000 1010 -> Border 2, MIC bit 3=1, Beeper bit 4=0
		assert.Equal(t, "Border Color", bus.BorderColor, uint8(2))
		assert.True(t, "Mic State", bus.MicState)
		assert.False(t, "Beeper State", bus.BeeperState)
		
		// Border 5 (Cyan), Beeper on
		bus.Out(0x00FE, 0x15) // 0001 0101 -> Border 5, Beeper bit 4=1
		assert.Equal(t, "Border Color", bus.BorderColor, uint8(5))
		assert.True(t, "Beeper State", bus.BeeperState)
	}),
	dsl.NewScenario("Tape input bit 6 via In port", func(t *testing.T) {
		bus := NewBus48()
		bus.TapeInState = true
		assert.Equal(t, "Tape bit 6 high", bus.In(0xFEFE) & 0x40, uint8(0x40))
		
		bus.TapeInState = false
		assert.Equal(t, "Tape bit 6 low", bus.In(0xFEFE) & 0x40, uint8(0x00))
	}),
}
