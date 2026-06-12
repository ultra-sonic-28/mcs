// Package gui implements the ZX Spectrum GUI logic.
package gui

import (
	"mcs/internal/machine/spectrum/keyboard"
	"github.com/hajimehoshi/ebiten/v2"
)

// KeyMap maps Ebitengine keys to Spectrum keys.
var KeyMap = map[ebiten.Key]keyboard.Key{
	ebiten.Key1: keyboard.Key1, ebiten.Key2: keyboard.Key2, ebiten.Key3: keyboard.Key3, ebiten.Key4: keyboard.Key4, ebiten.Key5: keyboard.Key5,
	ebiten.Key6: keyboard.Key6, ebiten.Key7: keyboard.Key7, ebiten.Key8: keyboard.Key8, ebiten.Key9: keyboard.Key9, ebiten.Key0: keyboard.Key0,
	ebiten.KeyQ: keyboard.KeyQ, ebiten.KeyW: keyboard.KeyW, ebiten.KeyE: keyboard.KeyE, ebiten.KeyR: keyboard.KeyR, ebiten.KeyT: keyboard.KeyT,
	ebiten.KeyY: keyboard.KeyY, ebiten.KeyU: keyboard.KeyU, ebiten.KeyI: keyboard.KeyI, ebiten.KeyO: keyboard.KeyO, ebiten.KeyP: keyboard.KeyP,
	ebiten.KeyA: keyboard.KeyA, ebiten.KeyS: keyboard.KeyS, ebiten.KeyD: keyboard.KeyD, ebiten.KeyF: keyboard.KeyF, ebiten.KeyG: keyboard.KeyG,
	ebiten.KeyH: keyboard.KeyH, ebiten.KeyJ: keyboard.KeyJ, ebiten.KeyK: keyboard.KeyK, ebiten.KeyL: keyboard.KeyL, ebiten.KeyEnter: keyboard.KeyEnter,
	ebiten.KeyShiftLeft: keyboard.KeyCapsShift, ebiten.KeyZ: keyboard.KeyZ, ebiten.KeyX: keyboard.KeyX, ebiten.KeyC: keyboard.KeyC, ebiten.KeyV: keyboard.KeyV,
	ebiten.KeyB: keyboard.KeyB, ebiten.KeyN: keyboard.KeyN, ebiten.KeyM: keyboard.KeyM, ebiten.KeyControlLeft: keyboard.KeySymbolShift, ebiten.KeySpace: keyboard.KeySpace,
}
