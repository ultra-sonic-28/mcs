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
	dsl.NewScenario("AddButton registers button in toolbar", func(t *testing.T) {
		tb := NewToolbar(20, "#D6CDC9")
		assert.Equal(t, "No buttons initially", len(tb.buttons), 0)
		btn := NewButton(16, 16, nil, nil)
		tb.AddButton(btn)
		assert.Equal(t, "One button after AddButton", len(tb.buttons), 1)
	}),
	dsl.NewScenario("Button is positioned by toolbar layout", func(t *testing.T) {
		tb := NewToolbar(20, "#D6CDC9")
		btn := NewButton(16, 16, nil, nil)
		tb.AddButton(btn)
		// x = buttonPadding = 2; y = (20-16)/2 = 2
		assert.Equal(t, "Button X position", btn.x, buttonPadding)
		assert.Equal(t, "Button Y centred", btn.y, (20-16)/2)
	}),
	dsl.NewScenario("Multiple buttons are laid out left-to-right", func(t *testing.T) {
		tb := NewToolbar(20, "#D6CDC9")
		btn1 := NewButton(16, 16, nil, nil)
		btn2 := NewButton(16, 16, nil, nil)
		tb.AddButton(btn1)
		tb.AddButton(btn2)
		// btn1: x=2; btn2: x=2+16+2=20
		assert.Equal(t, "Button 1 X", btn1.x, buttonPadding)
		assert.Equal(t, "Button 2 X", btn2.x, buttonPadding+16+buttonPadding)
	}),
}
