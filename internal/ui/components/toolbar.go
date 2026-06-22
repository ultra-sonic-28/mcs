// Package components implements reusable UI components for the MCS system.
package components

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Toolbar represents a configurable toolbar UI component.
type Toolbar struct {
	height int
	color  color.Color
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
		height: height,
		color:  c,
	}
}

// Height returns the height of the toolbar in logical pixels.
func (t *Toolbar) Height() int {
	return t.height
}

// Draw renders the toolbar on the target image.
func (t *Toolbar) Draw(screen *ebiten.Image) {
	if t.height <= 0 {
		return
	}
	// The toolbar is drawn at the top (y = 0 to y = height)
	w := screen.Bounds().Dx()
	rectImg := ebiten.NewImage(w, t.height)
	rectImg.Fill(t.color)

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(rectImg, op)
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
