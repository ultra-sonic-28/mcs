// Package spectrum implements the ZX Spectrum machine logic.
package display

import (
	"image/color"
)

const (
	ScreenWidth  = 256
	ScreenHeight = 192
)

// SpectrumPalette defines the 16 colors of the ZX Spectrum (8 basic colors in 2 brightness levels).
var SpectrumPalette = [16]color.RGBA{
	{0, 0, 0, 255},       // 0: Black
	{0, 0, 205, 255},     // 1: Blue
	{205, 0, 0, 255},     // 2: Red
	{205, 0, 205, 255},   // 3: Magenta
	{0, 205, 0, 255},     // 4: Green
	{0, 205, 205, 255},   // 5: Cyan
	{205, 205, 0, 255},   // 6: Yellow
	{205, 205, 205, 255}, // 7: White
	{0, 0, 0, 255},       // 8: Bright Black (same as black)
	{0, 0, 255, 255},     // 9: Bright Blue
	{255, 0, 0, 255},     // 10: Bright Red
	{255, 0, 255, 255},   // 11: Bright Magenta
	{0, 255, 0, 255},     // 12: Bright Green
	{0, 255, 255, 255},   // 13: Bright Cyan
	{255, 255, 0, 255},   // 14: Bright Yellow
	{255, 255, 255, 255}, // 15: Bright White
}

// Display handles the rendering of the Spectrum's display memory.
type Display struct {
	// FrameBuffer stores the RGBA pixels for the 256x192 screen.
	FrameBuffer [ScreenWidth * ScreenHeight * 4]uint8
	// FlashState toggles every 16 or 32 frames to handle flashing attributes.
	FlashState bool
}

// NewDisplay creates a new Display instance.
func NewDisplay() *Display {
	return &Display{}
}

// RenderFrame translates the memory provided into the FrameBuffer.
// displayMem is expected to be a slice of 6912 bytes (6144 pixels + 768 attributes).
func (d *Display) RenderFrame(displayMem []byte) {
	if len(displayMem) < 6912 {
		return
	}

	pixelBase := 0
	attrBase := 6144

	for y := 0; y < ScreenHeight; y++ {
		// Spectrum non-linear memory calculation:
		// The 192 lines are split into 3 blocks of 64 lines.
		// Within each block, lines are not sequential but interleaved by 8.
		// Y = [00][Y7 Y6] [Y2 Y1 Y0] [Y5 Y4 Y3]
		// In memory address bits: 010 [Y7 Y6] [Y5 Y4 Y3] [Y2 Y1 Y0] [X4 X3 X2 X1 X0]
		
		block := y / 64
		lineInBlock := y % 64
		rowInBlock := lineInBlock / 8
		lineInRow := lineInBlock % 8
		
		memY := (block << 11) | (lineInRow << 8) | (rowInBlock << 5)
		attrY := (block << 8) | (rowInBlock << 5)

		for xByte := 0; xByte < 32; xByte++ {
			pixelByte := displayMem[pixelBase+memY+xByte]
			attrByte := displayMem[attrBase+attrY+xByte]

			// Attribute bits: [Flash][Bright][Paper2-0][Ink2-0]
			inkIdx := attrByte & 0x07
			paperIdx := (attrByte >> 3) & 0x07
			bright := (attrByte >> 6) & 0x01
			flash := (attrByte >> 7) & 0x01

			if flash != 0 && d.FlashState {
				inkIdx, paperIdx = paperIdx, inkIdx
			}

			inkColor := SpectrumPalette[inkIdx|(bright<<3)]
			paperColor := SpectrumPalette[paperIdx|(bright<<3)]

			for bit := 0; bit < 8; bit++ {
				x := (xByte << 3) + bit
				pixelIdx := (y*ScreenWidth + x) * 4
				
				// Bit 7 is the leftmost pixel
				if (pixelByte & (0x80 >> bit)) != 0 {
					d.FrameBuffer[pixelIdx+0] = inkColor.R
					d.FrameBuffer[pixelIdx+1] = inkColor.G
					d.FrameBuffer[pixelIdx+2] = inkColor.B
					d.FrameBuffer[pixelIdx+3] = inkColor.A
				} else {
					d.FrameBuffer[pixelIdx+0] = paperColor.R
					d.FrameBuffer[pixelIdx+1] = paperColor.G
					d.FrameBuffer[pixelIdx+2] = paperColor.B
					d.FrameBuffer[pixelIdx+3] = paperColor.A
				}
			}
		}
	}
}
