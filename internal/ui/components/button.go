// Package components implements reusable UI components for the MCS system.
package components

import (
	"image/color"

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
//   - Pressed:  light gray  (#2A2A2A)
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
	if len(b.bitmap) > 0 {
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
