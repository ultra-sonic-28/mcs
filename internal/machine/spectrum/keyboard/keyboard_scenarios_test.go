// Package keyboard implements the ZX Spectrum keyboard logic.
package keyboard

import (
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var keyboardScenarios = []dsl.Scenario{
	dsl.NewScenario("Keyboard scanning basic", func(t *testing.T) {
		k := NewKeyboard()
		// No keys pressed: bits 0-4 should be 1
		// Mask 0xFE selects row 0
		val := k.Scan(0xFE)
		assert.Equal(t, "Scan 0xFE (no keys)", val, uint8(0x1F))
		
		// Press 'Z' (Bit 1 of row 0xFE)
		k.SetKeyState(KeyZ, true)
		val = k.Scan(0xFE)
		assert.Equal(t, "Scan 0xFE (Z pressed)", val, uint8(0x1D)) // 11101
		
		// Release 'Z'
		k.SetKeyState(KeyZ, false)
		val = k.Scan(0xFE)
		assert.Equal(t, "Scan 0xFE (Z released)", val, uint8(0x1F))
	}),
	dsl.NewScenario("Keyboard multiple rows scanning", func(t *testing.T) {
		k := NewKeyboard()
		k.SetKeyState(KeyA, true) // Bit 0 of row 0xFD
		k.SetKeyState(KeyQ, true) // Bit 0 of row 0xFB
		
		// Scan 0xFD (A pressed)
		assert.Equal(t, "Scan 0xFD", k.Scan(0xFD), uint8(0x1E))
		
		// Scan 0xFB (Q pressed)
		assert.Equal(t, "Scan 0xFB", k.Scan(0xFB), uint8(0x1E))
		
		// Scan with both bits 0: 0xFC (0xF9 & 0xFB?) No, 0xFD is 1111 1101, 0xFB is 1111 1011. 
		// 0xFD & 0xFB = 1111 1001 = 0xF9
		assert.Equal(t, "Scan 0xF9 (both rows)", k.Scan(0xF9), uint8(0x1E))
	}),
}
