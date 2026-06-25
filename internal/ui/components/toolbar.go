// Package components implements reusable UI components for the MCS system.
package components

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// buttonPadding is the horizontal gap between buttons and between the toolbar
// edge and the first button, in logical pixels.
const buttonPadding = 2

// Toolbar represents a configurable toolbar UI component.
// It owns a list of Button components that are laid out left-to-right.
type Toolbar struct {
	height  int
	color   color.Color
	buttons []*Button
}

// NewToolbar creates a new Toolbar component with the given height and color string.
// If the color string is invalid, it falls back to a default gray color (#D6CDC9).
func NewToolbar(height int, hexColor string) *Toolbar {
	fallback := color.RGBA{R: 214, G: 205, B: 201, A: 255}
	c := parseHexColor(hexColor, fallback)
	if height < 0 {
		height = 0
	}
	return &Toolbar{
		height:  height,
		color:   c,
		buttons: make([]*Button, 0),
	}
}

// Height returns the height of the toolbar in logical pixels.
func (t *Toolbar) Height() int {
	return t.height
}

// AddButton appends a button to the toolbar.
// Buttons are drawn left-to-right in the order they are added.
func (t *Toolbar) AddButton(b *Button) {
	t.buttons = append(t.buttons, b)
	t.layoutButtons()
}

// layoutButtons recalculates each button's position so they are centred
// vertically and packed left-to-right with buttonPadding gaps.
func (t *Toolbar) layoutButtons() {
	x := buttonPadding
	for _, b := range t.buttons {
		// Centre vertically within the toolbar.
		y := (t.height - b.Height) / 2
		if y < 0 {
			y = 0
		}
		b.SetPosition(x, y)
		x += b.Width + buttonPadding
	}
}

// Update forwards input events to all buttons in the toolbar.
// Must be called from the Ebitengine Update loop.
func (t *Toolbar) Update() {
	for _, b := range t.buttons {
		b.Update()
	}
}

// Draw renders the toolbar background and all its buttons on the target image.
func (t *Toolbar) Draw(screen *ebiten.Image) {
	if t.height <= 0 {
		return
	}
	// Draw toolbar background spanning the full screen width.
	w := screen.Bounds().Dx()
	rectImg := ebiten.NewImage(w, t.height)
	rectImg.Fill(t.color)

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(rectImg, op)

	// Draw buttons on top of the background.
	for _, b := range t.buttons {
		b.Draw(screen)
	}
}

// parseHexColor parses a hexadecimal color string (e.g., "#D6CDC9") and returns the corresponding color.RGBA.
func parseHexColor(s string, fallback color.RGBA) color.RGBA {
	if s == "" {
		return fallback
	}
	if s[0] == '#' {
		s = s[1:]
	}
	var r, g, b uint8
	var a uint8 = 255
	if len(s) == 6 {
		n, err := fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b)
		if err == nil && n == 3 {
			return color.RGBA{R: r, G: g, B: b, A: a}
		}
	} else if len(s) == 8 {
		n, err := fmt.Sscanf(s, "%02x%02x%02x%02x", &r, &g, &b, &a)
		if err == nil && n == 4 {
			return color.RGBA{R: r, G: g, B: b, A: a}
		}
	}
	return fallback
}
