// Package components implements reusable UI components for the MCS system.
package components

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png" // register PNG decoder

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Bitmap is a 2D array of booleans that represents a monochrome icon.
// Each true value corresponds to a foreground pixel; false is transparent.
type Bitmap [][]bool

// Button represents a clickable toolbar button with a bitmap icon.
//
// The button renders its bitmap centered within its bounds and supports
// three visual states: normal, hovered, and pressed.
type Button struct {
	// Width is the total width of the button in logical pixels.
	Width int
	// Height is the total height of the button in logical pixels.
	Height int
	// bitmap is the monochrome icon to draw inside the button.
	bitmap Bitmap
	// iconImage is the original PNG image rendered for image-based buttons.
	iconImage *ebiten.Image

	// x and y are the top-left position of the button on the screen.
	x, y int

	// onClick is the callback triggered when the button is released.
	onClick func()
}

// NewButton creates a new Button with the given dimensions and bitmap icon.
// Width and height must be positive; negative or zero values are clamped to 1.
// bitmap may be nil, in which case the button renders as a plain rectangle.
func NewButton(width, height int, bitmap Bitmap, onClick func()) *Button {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return &Button{
		Width:   width,
		Height:  height,
		bitmap:  bitmap,
		onClick: onClick,
	}
}

// NewButtonFromImageData decodes a PNG from raw bytes and builds a Button whose
// size matches the image dimensions. If the image height exceeds maxHeight, the
// button is scaled proportionally so that its height equals maxHeight.
//
// The image is converted to a monochrome Bitmap: any pixel whose alpha channel
// is >= 128 and luminance < 128 is treated as a foreground (true) pixel.
// Returns nil and an error if the PNG data cannot be decoded.
func NewButtonFromImageData(pngData []byte, maxHeight int, onClick func()) (*Button, error) {
	img, _, err := image.Decode(bytes.NewReader(pngData))
	if err != nil {
		return nil, err
	}
	return newButtonFromImage(img, maxHeight, onClick), nil
}

// newButtonFromImage converts an image.Image to a Button, scaling it if needed.
func newButtonFromImage(img image.Image, maxHeight int, onClick func()) *Button {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	// Scale proportionally when the image is taller than the allowed height.
	dstW, dstH := srcW, srcH
	if maxHeight > 0 && srcH > maxHeight {
		dstH = maxHeight
		dstW = srcW * maxHeight / srcH
		if dstW < 1 {
			dstW = 1
		}
	}

	b := NewButton(dstW, dstH, nil, onClick)
	b.iconImage = ebiten.NewImageFromImage(img)
	return b
}

// SetPosition sets the top-left position of the button on the screen.
func (b *Button) SetPosition(x, y int) {
	b.x = x
	b.y = y
}

// isHovered returns true when the mouse cursor is inside the button's bounds.
func (b *Button) isHovered() bool {
	mx, my := ebiten.CursorPosition()
	return mx >= b.x && mx < b.x+b.Width &&
		my >= b.y && my < b.y+b.Height
}

// Update processes input events and fires the onClick callback when the
// left mouse button is released while the cursor is within the button.
func (b *Button) Update() {
	if b.onClick != nil && b.isHovered() && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		b.onClick()
	}
}

// Draw renders the button onto the provided screen image at the position set
// by SetPosition. The button background changes color on hover and press.
//
// Visual states:
//   - Normal:   dark gray  (#3A3A3A)
//   - Hovered:  medium gray (#5A5A5A)
//   - Pressed:  very dark gray (#2A2A2A)
func (b *Button) Draw(screen *ebiten.Image) {
	// --- Background color by state ---
	bg := color.RGBA{R: 58, G: 58, B: 58, A: 255} // normal
	if b.isHovered() {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			bg = color.RGBA{R: 42, G: 42, B: 42, A: 255} // pressed
		} else {
			bg = color.RGBA{R: 90, G: 90, B: 90, A: 255} // hovered
		}
	}

	// --- Draw background ---
	btnImg := ebiten.NewImage(b.Width, b.Height)
	btnImg.Fill(bg)

	// --- Draw bitmap icon centered ---
	if b.iconImage != nil {
		srcW, srcH := b.iconImage.Bounds().Dx(), b.iconImage.Bounds().Dy()
		scale := 1.0
		if srcW > 0 && srcH > 0 {
			scale = float64(b.Width) / float64(srcW)
			if float64(b.Height)/float64(srcH) < scale {
				scale = float64(b.Height) / float64(srcH)
			}
			if scale > 1 {
				scale = 1
			}
		}
		offX := (b.Width - int(float64(srcW)*scale)) / 2
		offY := (b.Height - int(float64(srcH)*scale)) / 2
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(offX), float64(offY))
		btnImg.DrawImage(b.iconImage, op)
	} else if len(b.bitmap) > 0 {
		iconH := len(b.bitmap)
		iconW := 0
		if iconH > 0 {
			iconW = len(b.bitmap[0])
		}
		offX := (b.Width - iconW) / 2
		offY := (b.Height - iconH) / 2
		fg := color.RGBA{R: 220, G: 220, B: 220, A: 255}
		for row := 0; row < iconH; row++ {
			for col := 0; col < iconW; col++ {
				if col < len(b.bitmap[row]) && b.bitmap[row][col] {
					btnImg.Set(offX+col, offY+row, fg)
				}
			}
		}
	}

	// --- Blit button onto screen at its position ---
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(btnImg, op)
}
