// Package spectrum implements the ZX Spectrum machine logic.
package spectrum

// Key represents a physical key on the ZX Spectrum keyboard.
type Key int

const (
	// Half-row: 0xFEFE
	KeyCapsShift Key = iota
	KeyZ
	KeyX
	KeyC
	KeyV
	// Half-row: 0xFDFE
	KeyA
	KeyS
	KeyD
	KeyF
	KeyG
	// Half-row: 0xFBFE
	KeyQ
	KeyW
	KeyE
	KeyR
	KeyT
	// Half-row: 0xF7FE
	Key1
	Key2
	Key3
	Key4
	Key5
	// Half-row: 0xEFFE
	Key0
	Key9
	Key8
	Key7
	Key6
	// Half-row: 0xDFFE
	KeyP
	KeyO
	KeyI
	KeyU
	KeyY
	// Half-row: 0xBFFE
	KeyEnter
	KeyL
	KeyK
	KeyJ
	KeyH
	// Half-row: 0x7FFE
	KeySpace
	KeySymbolShift
	KeyM
	KeyN
	KeyB
)

// Keyboard manages the state of the Spectrum keyboard matrix.
type Keyboard struct {
	// keys stores the pressed state of each key (true = pressed).
	keys [40]bool
}

// NewKeyboard creates a new Keyboard instance.
func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

// SetKeyState updates the state of a specific key.
func (k *Keyboard) SetKeyState(key Key, pressed bool) {
	if key >= 0 && key < 40 {
		k.keys[key] = pressed
	}
}

// Scan returns the 5-bit state for a specific half-row mask.
// The mask is the high byte of the port address (A8-A15).
// A '0' bit in the mask selects the corresponding half-row.
// If multiple bits are '0', the results are ANDed together.
func (k *Keyboard) Scan(mask uint8) uint8 {
	result := uint8(0x1F) // Start with all bits 1 (not pressed)

	// Rows are active low. We check each bit of the mask.
	// Bit 0: 0xFE (Caps Shift, Z, X, C, V)
	if mask&0x01 == 0 {
		result &= k.getRowBits(0)
	}
	// Bit 1: 0xFD (A, S, D, F, G)
	if mask&0x02 == 0 {
		result &= k.getRowBits(5)
	}
	// Bit 2: 0xFB (Q, W, E, R, T)
	if mask&0x04 == 0 {
		result &= k.getRowBits(10)
	}
	// Bit 3: 0xF7 (1, 2, 3, 4, 5)
	if mask&0x08 == 0 {
		result &= k.getRowBits(15)
	}
	// Bit 4: 0xEF (0, 9, 8, 7, 6)
	if mask&0x10 == 0 {
		result &= k.getRowBits(20)
	}
	// Bit 5: 0xDF (P, O, I, U, Y)
	if mask&0x20 == 0 {
		result &= k.getRowBits(25)
	}
	// Bit 6: 0xBF (Enter, L, K, J, H)
	if mask&0x40 == 0 {
		result &= k.getRowBits(30)
	}
	// Bit 7: 0x7F (Space, Symbol Shift, M, N, B)
	if mask&0x80 == 0 {
		result &= k.getRowBits(35)
	}

	return result
}

// getRowBits returns the 5-bit state for a set of 5 keys starting at the offset.
// A '0' bit means pressed, '1' means not pressed.
func (k *Keyboard) getRowBits(offset int) uint8 {
	bits := uint8(0x1F)
	if k.keys[offset+0] { bits &= ^uint8(0x01) }
	if k.keys[offset+1] { bits &= ^uint8(0x02) }
	if k.keys[offset+2] { bits &= ^uint8(0x04) }
	if k.keys[offset+3] { bits &= ^uint8(0x08) }
	if k.keys[offset+4] { bits &= ^uint8(0x10) }
	return bits
}
