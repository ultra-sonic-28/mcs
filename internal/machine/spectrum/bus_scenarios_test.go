package spectrum

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
}
