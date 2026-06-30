// Package gui defines the scenarios for gui tests.
package gui

import (
	"fmt"
	"image/color"
	"mcs/internal/machine/spectrum/keyboard"
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// expectedMapping defines a test case for KeyMap mapping.
type expectedMapping struct {
	ebitenKey ebiten.Key
	specKey   keyboard.Key
}

// guiScenarios lists the scenario-based tests for the GUI logic.
var guiScenarios = []dsl.Scenario{
	dsl.NewScenario("Verify KeyMap mappings", func(t *testing.T) {
		// Define a table of expected mappings.
		cases := []expectedMapping{
			{ebiten.Key1, keyboard.Key1},
			{ebiten.Key2, keyboard.Key2},
			{ebiten.Key3, keyboard.Key3},
			{ebiten.Key4, keyboard.Key4},
			{ebiten.Key5, keyboard.Key5},
			{ebiten.Key6, keyboard.Key6},
			{ebiten.Key7, keyboard.Key7},
			{ebiten.Key8, keyboard.Key8},
			{ebiten.Key9, keyboard.Key9},
			{ebiten.Key0, keyboard.Key0},
			{ebiten.KeyQ, keyboard.KeyQ},
			{ebiten.KeyW, keyboard.KeyW},
			{ebiten.KeyE, keyboard.KeyE},
			{ebiten.KeyR, keyboard.KeyR},
			{ebiten.KeyT, keyboard.KeyT},
			{ebiten.KeyY, keyboard.KeyY},
			{ebiten.KeyU, keyboard.KeyU},
			{ebiten.KeyI, keyboard.KeyI},
			{ebiten.KeyO, keyboard.KeyO},
			{ebiten.KeyP, keyboard.KeyP},
			{ebiten.KeyA, keyboard.KeyA},
			{ebiten.KeyS, keyboard.KeyS},
			{ebiten.KeyD, keyboard.KeyD},
			{ebiten.KeyF, keyboard.KeyF},
			{ebiten.KeyG, keyboard.KeyG},
			{ebiten.KeyH, keyboard.KeyH},
			{ebiten.KeyJ, keyboard.KeyJ},
			{ebiten.KeyK, keyboard.KeyK},
			{ebiten.KeyL, keyboard.KeyL},
			{ebiten.KeyEnter, keyboard.KeyEnter},
			{ebiten.KeyShiftLeft, keyboard.KeyCapsShift},
			{ebiten.KeyZ, keyboard.KeyZ},
			{ebiten.KeyX, keyboard.KeyX},
			{ebiten.KeyC, keyboard.KeyC},
			{ebiten.KeyV, keyboard.KeyV},
			{ebiten.KeyB, keyboard.KeyB},
			{ebiten.KeyN, keyboard.KeyN},
			{ebiten.KeyM, keyboard.KeyM},
			{ebiten.KeyControlLeft, keyboard.KeySymbolShift},
			{ebiten.KeySpace, keyboard.KeySpace},
		}

		for _, tc := range cases {
			got, exists := KeyMap[tc.ebitenKey]
			assert.True(t, fmt.Sprintf("Key %v exists in KeyMap", tc.ebitenKey), exists)
			assert.Equal(t, fmt.Sprintf("Key %v maps to correct Spectrum Key", tc.ebitenKey), got, tc.specKey)
		}

		assert.Equal(t, "KeyMap size is exactly 40 keys", len(KeyMap), 40)
	}),

	dsl.NewScenario("Verify SmallFont5x7 contents", func(t *testing.T) {
		// Verify some known glyphs exist in the font map.
		runes := []rune{'0', '9', 'A', 'Z', 'a', 'z', ' ', '.', '|', '-'}
		for _, r := range runes {
			_, exists := SmallFont5x7[r]
			assert.True(t, fmt.Sprintf("Rune %q exists in SmallFont5x7", r), exists)
		}

		// Verify glyph height is exactly 7 for all runes in the font.
		for r, glyph := range SmallFont5x7 {
			assert.Equal(t, fmt.Sprintf("Glyph %q row length is 7", r), len(glyph), 7)
		}

		// Space should have zero rows.
		glyphSpace := SmallFont5x7[' ']
		for row := 0; row < 7; row++ {
			assert.Equal(t, fmt.Sprintf("Space glyph row %d is empty", row), glyphSpace[row], uint8(0))
		}
	}),

	dsl.NewScenario("Verify DrawSmallText executes without panic", func(t *testing.T) {
		img := ebiten.NewImage(100, 20)
		clr := color.RGBA{R: 255, G: 0, B: 0, A: 255}

		// 1. Draw standard text.
		DrawSmallText(img, "HelloWorld123", 0, 0, clr)

		// 2. Draw text with unknown characters (should trigger fallback without panic).
		DrawSmallText(img, "Hello#World?", 0, 0, clr)

		// 3. Draw text at different positions.
		DrawSmallText(img, "Test", 10, 10, clr)

		// 4. Draw with different colors.
		DrawSmallText(img, "Color", 0, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})
		DrawSmallText(img, "White", 0, 0, color.White)

		// 5. Draw empty string.
		DrawSmallText(img, "", 0, 0, clr)

		// Since we didn't panic, this test scenario successfully proves the execution path works.
		assert.True(t, "DrawSmallText finished execution without panics", true)
	}),
}
