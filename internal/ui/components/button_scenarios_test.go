// Package components defines the scenarios for button tests.
package components

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

// simpleBitmap is a small 3×3 icon used across button scenarios.
var simpleBitmap = Bitmap{
	{true, false, true},
	{false, true, false},
	{true, false, true},
}

var buttonScenarios = []dsl.Scenario{
	dsl.NewScenario("Create button with valid dimensions", func(t *testing.T) {
		btn := NewButton(20, 20, simpleBitmap, nil)
		assert.Equal(t, "Width matches", btn.Width, 20)
		assert.Equal(t, "Height matches", btn.Height, 20)
	}),
	dsl.NewScenario("Create button with zero width is clamped to 1", func(t *testing.T) {
		btn := NewButton(0, 10, nil, nil)
		assert.Equal(t, "Width clamped to 1", btn.Width, 1)
		assert.Equal(t, "Height unchanged", btn.Height, 10)
	}),
	dsl.NewScenario("Create button with zero height is clamped to 1", func(t *testing.T) {
		btn := NewButton(10, 0, nil, nil)
		assert.Equal(t, "Width unchanged", btn.Width, 10)
		assert.Equal(t, "Height clamped to 1", btn.Height, 1)
	}),
	dsl.NewScenario("Create button with negative dimensions are clamped to 1", func(t *testing.T) {
		btn := NewButton(-5, -3, nil, nil)
		assert.Equal(t, "Width clamped to 1", btn.Width, 1)
		assert.Equal(t, "Height clamped to 1", btn.Height, 1)
	}),
	dsl.NewScenario("SetPosition updates button position", func(t *testing.T) {
		btn := NewButton(16, 16, simpleBitmap, nil)
		btn.SetPosition(10, 20)
		assert.Equal(t, "X position", btn.x, 10)
		assert.Equal(t, "Y position", btn.y, 20)
	}),
	dsl.NewScenario("Create button with nil bitmap does not panic", func(t *testing.T) {
		btn := NewButton(16, 16, nil, nil)
		assert.Equal(t, "Width matches", btn.Width, 16)
		assert.Equal(t, "bitmap is nil", len(btn.bitmap), 0)
	}),
	dsl.NewScenario("onClick callback is stored", func(t *testing.T) {
		called := false
		btn := NewButton(16, 16, simpleBitmap, func() { called = true })
		assert.Equal(t, "callback not yet called", called, false)
		// Verify the function reference is stored (not nil)
		assert.True(t, "callback is not nil", btn.onClick != nil)
	}),
}
