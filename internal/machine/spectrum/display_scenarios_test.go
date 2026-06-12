// Package spectrum implements the ZX Spectrum machine logic.
package spectrum

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var displayScenarios = []dsl.Scenario{
	dsl.NewScenario("Display pixel memory mapping", func(t *testing.T) {
		bus := NewBus48()
		
		// Set pixel at (0, 0)
		// Address: 010 [00] [000] [000] [00000] -> 0x4000
		bus.Write(0x4000, 0x80) // Bit 7 set
		
		// Set attribute at (0, 0)
		// Address: 0x5800 + (0 * 32) + 0 = 0x5800
		bus.Write(0x5800, 0x42) // Bright=1, Paper=0 (Black), Ink=2 (Red)
		
		bus.Display.RenderFrame(bus.GetDisplayMemory())
		
		// Check pixel (0, 0) - should be Bright Red
		pixelIdx := 0
		inkColor := SpectrumPalette[2|8] // Bright Red
		assert.Equal(t, "Pixel (0,0) Red", bus.Display.FrameBuffer[pixelIdx+0], inkColor.R)
		assert.Equal(t, "Pixel (0,0) Green", bus.Display.FrameBuffer[pixelIdx+1], inkColor.G)
		assert.Equal(t, "Pixel (0,0) Blue", bus.Display.FrameBuffer[pixelIdx+2], inkColor.B)
		
		// Check pixel (1, 0) - should be Bright Black (Paper)
		pixelIdx = 4
		paperColor := SpectrumPalette[0|8] // Bright Black
		assert.Equal(t, "Pixel (1,0) Red", bus.Display.FrameBuffer[pixelIdx+0], paperColor.R)
		assert.Equal(t, "Pixel (1,0) Green", bus.Display.FrameBuffer[pixelIdx+1], paperColor.G)
		assert.Equal(t, "Pixel (1,0) Blue", bus.Display.FrameBuffer[pixelIdx+2], paperColor.B)
	}),
	dsl.NewScenario("Display non-linear mapping (Line 1)", func(t *testing.T) {
		bus := NewBus48()
		
		// Y = 1: [00][00][001][000]
		// memY = (0 << 11) | (1 << 8) | (0 << 5) = 0x0100
		// Address: 0x4000 + 0x0100 = 0x4100
		bus.Write(0x4100, 0x80)
		bus.Write(0x5800, 0x07) // White on Black
		
		bus.Display.RenderFrame(bus.GetDisplayMemory())
		
		// Check pixel (0, 1)
		pixelIdx := (1 * ScreenWidth + 0) * 4
		assert.Equal(t, "Pixel (0,1) R", bus.Display.FrameBuffer[pixelIdx+0], uint8(205)) // White
	}),
	dsl.NewScenario("Display flash toggling", func(t *testing.T) {
		bus := NewBus48()
		
		bus.Write(0x4000, 0x80)
		bus.Write(0x5800, 0x87) // Flash=1, Bright=0, Paper=0, Ink=7 (White on Black)
		
		// FlashState false: Ink is White
		bus.Display.FlashState = false
		bus.Display.RenderFrame(bus.GetDisplayMemory())
		assert.Equal(t, "Pixel (0,0) White", bus.Display.FrameBuffer[0], uint8(205))
		
		// FlashState true: Ink and Paper swap -> Ink is Black
		bus.Display.FlashState = true
		bus.Display.RenderFrame(bus.GetDisplayMemory())
		assert.Equal(t, "Pixel (0,0) Black", bus.Display.FrameBuffer[0], uint8(0))
	}),
}
