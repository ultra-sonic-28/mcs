// Package components defines the scenarios for button tests.
package components

import (
	toolbarassets "mcs/assets/ui/toolbar"
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
	dsl.NewScenario("NewButtonFromImageData decodes valid PNG without error", func(t *testing.T) {
		btn, err := NewButtonFromImageData(toolbarassets.QuitApp, 20, nil)
		assert.Equal(t, "No decode error", err, nil)
		assert.True(t, "Button is not nil", btn != nil)
		// The quit-app.png is 16x16; maxHeight=20 means no scaling, so width=16 height=16.
		assert.True(t, "Width > 0", btn.Width > 0)
		assert.True(t, "Height <= 20", btn.Height <= 20)
	}),
	dsl.NewScenario("NewButtonFromImageData scales button when image taller than maxHeight", func(t *testing.T) {
		// quit-app.png is 16×16; maxHeight=8 should scale it to 8×8.
		btn, err := NewButtonFromImageData(toolbarassets.QuitApp, 8, nil)
		assert.Equal(t, "No decode error", err, nil)
		assert.True(t, "Button is not nil", btn != nil)
		assert.Equal(t, "Height capped at maxHeight", btn.Height, 8)
	}),
	dsl.NewScenario("NewButtonFromImageData returns error on invalid PNG data", func(t *testing.T) {
		btn, err := NewButtonFromImageData([]byte("not a png"), 20, nil)
		assert.True(t, "Error is returned", err != nil)
		assert.True(t, "Button is nil on error", btn == nil)
	}),
}
