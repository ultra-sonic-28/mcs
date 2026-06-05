// Package spectrum implements the ZX Spectrum 48K machine logic.
package spectrum

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// KeyMap maps Ebitengine keys to Spectrum keys.
var KeyMap = map[ebiten.Key]Key{
	ebiten.Key1: Key1, ebiten.Key2: Key2, ebiten.Key3: Key3, ebiten.Key4: Key4, ebiten.Key5: Key5,
	ebiten.Key6: Key6, ebiten.Key7: Key7, ebiten.Key8: Key8, ebiten.Key9: Key9, ebiten.Key0: Key0,
	ebiten.KeyQ: KeyQ, ebiten.KeyW: KeyW, ebiten.KeyE: KeyE, ebiten.KeyR: KeyR, ebiten.KeyT: KeyT,
	ebiten.KeyY: KeyY, ebiten.KeyU: KeyU, ebiten.KeyI: KeyI, ebiten.KeyO: KeyO, ebiten.KeyP: KeyP,
	ebiten.KeyA: KeyA, ebiten.KeyS: KeyS, ebiten.KeyD: KeyD, ebiten.KeyF: KeyF, ebiten.KeyG: KeyG,
	ebiten.KeyH: KeyH, ebiten.KeyJ: KeyJ, ebiten.KeyK: KeyK, ebiten.KeyL: KeyL, ebiten.KeyEnter: KeyEnter,
	ebiten.KeyShiftLeft: KeyCapsShift, ebiten.KeyZ: KeyZ, ebiten.KeyX: KeyX, ebiten.KeyC: KeyC, ebiten.KeyV: KeyV,
	ebiten.KeyB: KeyB, ebiten.KeyN: KeyN, ebiten.KeyM: KeyM, ebiten.KeyControlLeft: KeySymbolShift, ebiten.KeySpace: KeySpace,
}

// UpdateKeyboard reads the host keyboard state and updates the Spectrum keyboard.
func (m *Machine) UpdateKeyboard() {
	if m.autoStartEnabled {
		// During auto-start typing, we don't want physical keyboard to interfere.
		return
	}
	for eKey, sKey := range KeyMap {
		m.Bus.Keyboard.SetKeyState(sKey, ebiten.IsKeyPressed(eKey))
	}
}

// Update implements the ebiten.Game interface.
func (m *Machine) Update() error {
	m.UpdateKeyboard()
	// In Ebiten, Update is called 60 times per second by default.
	// Spectrum runs at 50Hz, so we might need some adjustment if we want perfect timing,
	// but for now, running one frame per Update is a good start.
	m.RunFrame()
	return nil
}

// Draw implements the ebiten.Game interface.
func (m *Machine) Draw(screen *ebiten.Image) {
	m.Bus.Display.RenderFrame(m.Bus.GetDisplayMemory())
	screen.WritePixels(m.Bus.Display.FrameBuffer[:])
}

// Layout implements the ebiten.Game interface.
func (m *Machine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
