// Package display implements the ZX Spectrum display logic.
package display

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var displayScenarios = []dsl.Scenario{
	dsl.NewScenario("Display pixel memory mapping", func(t *testing.T) {
		d := NewDisplay()
		mem := make([]byte, 6912)
		
		// Set pixel at (0, 0)
		// Address: 0x0000 (relative to display memory)
		mem[0x0000] = 0x80 // Bit 7 set
		
		// Set attribute at (0, 0)
		// Address: 6144 + 0 = 6144 (0x1800)
		mem[6144] = 0x42 // Bright=1, Paper=0 (Black), Ink=2 (Red)
		
		d.RenderFrame(mem)
		
		// Check pixel (0, 0) - should be Bright Red
		pixelIdx := 0
		inkColor := SpectrumPalette[2|8] // Bright Red
		assert.Equal(t, "Pixel (0,0) Red", d.FrameBuffer[pixelIdx+0], inkColor.R)
		assert.Equal(t, "Pixel (0,0) Green", d.FrameBuffer[pixelIdx+1], inkColor.G)
		assert.Equal(t, "Pixel (0,0) Blue", d.FrameBuffer[pixelIdx+2], inkColor.B)
		
		// Check pixel (1, 0) - should be Bright Black (Paper)
		pixelIdx = 4
		paperColor := SpectrumPalette[0|8] // Bright Black
		assert.Equal(t, "Pixel (1,0) Red", d.FrameBuffer[pixelIdx+0], paperColor.R)
		assert.Equal(t, "Pixel (1,0) Green", d.FrameBuffer[pixelIdx+1], paperColor.G)
		assert.Equal(t, "Pixel (1,0) Blue", d.FrameBuffer[pixelIdx+2], paperColor.B)
	}),
	dsl.NewScenario("Display non-linear mapping (Line 1)", func(t *testing.T) {
		d := NewDisplay()
		mem := make([]byte, 6912)
		
		// Y = 1: [00][00][001][000]
		// memY = (0 << 11) | (1 << 8) | (0 << 5) = 0x0100
		mem[0x0100] = 0x80
		mem[6144] = 0x07 // White on Black
		
		d.RenderFrame(mem)
		
		// Check pixel (0, 1)
		pixelIdx := (1 * ScreenWidth + 0) * 4
		assert.Equal(t, "Pixel (0,1) R", d.FrameBuffer[pixelIdx+0], uint8(205)) // White
	}),
	dsl.NewScenario("Display flash toggling", func(t *testing.T) {
		d := NewDisplay()
		mem := make([]byte, 6912)
		
		mem[0x0000] = 0x80
		mem[6144] = 0x87 // Flash=1, Bright=0, Paper=0, Ink=7 (White on Black)
		
		// FlashState false: Ink is White
		d.FlashState = false
		d.RenderFrame(mem)
		assert.Equal(t, "Pixel (0,0) White", d.FrameBuffer[0], uint8(205))
		
		// FlashState true: Ink and Paper swap -> Ink is Black
		d.FlashState = true
		d.RenderFrame(mem)
		assert.Equal(t, "Pixel (0,0) Black", d.FrameBuffer[0], uint8(0))
	}),
}
