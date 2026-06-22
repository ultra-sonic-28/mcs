// Package components defines the scenarios for toolbar tests.
package components

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var toolbarScenarios = []dsl.Scenario{
	dsl.NewScenario("Create Toolbar with valid values", func(t *testing.T) {
		tb := NewToolbar(25, "#FF0000")
		assert.Equal(t, "Height matches", tb.Height(), 25)

		r, g, b, a := tb.color.RGBA()
		assert.Equal(t, "R matches", r, uint32(65535))
		assert.Equal(t, "G matches", g, uint32(0))
		assert.Equal(t, "B matches", b, uint32(0))
		assert.Equal(t, "A matches", a, uint32(65535))
	}),
	dsl.NewScenario("Create Toolbar with default fallback on invalid color", func(t *testing.T) {
		tb := NewToolbar(15, "invalid")
		assert.Equal(t, "Height matches", tb.Height(), 15)

		r, g, b, a := tb.color.RGBA()
		// #D6CDC9 -> R: 214, G: 205, B: 201
		// 214 * 257 = 54998
		// 205 * 257 = 52685
		// 201 * 257 = 51657
		assert.Equal(t, "Fallback R", r, uint32(54998))
		assert.Equal(t, "Fallback G", g, uint32(52685))
		assert.Equal(t, "Fallback B", b, uint32(51657))
		assert.Equal(t, "Fallback A", a, uint32(65535))
	}),
	dsl.NewScenario("Create Toolbar with negative height clamped to 0", func(t *testing.T) {
		tb := NewToolbar(-5, "#00FF00")
		assert.Equal(t, "Negative toolbar height clamped to 0", tb.Height(), 0)
	}),
}
