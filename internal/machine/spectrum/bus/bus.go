// Package bus implements the ZX Spectrum bus logic.
package bus

import (
	"mcs/internal/machine/spectrum/display"
	"mcs/internal/machine/spectrum/keyboard"
	"mcs/internal/machine/spectrum/tape"
)

// BaseBus contains the common components and state for all Spectrum models.
type BaseBus struct {
	Keyboard *keyboard.Keyboard
	Display  *display.Display
	Tape     *tape.Tape

	// ULA State
	BorderColor uint8
	BeeperState bool
	MicState    bool
	TapeInState bool
}

func (b *BaseBus) GetKeyboard() *keyboard.Keyboard { return b.Keyboard }
func (b *BaseBus) GetTape() *tape.Tape             { return b.Tape }
func (b *BaseBus) GetDisplay() *display.Display     { return b.Display }
func (b *BaseBus) GetBorderColor() uint8           { return b.BorderColor }
func (b *BaseBus) GetTapeInState() bool            { return b.TapeInState }
func (b *BaseBus) SetTapeInState(state bool)       { b.TapeInState = state }
func (b *BaseBus) GetBeeperState() bool            { return b.BeeperState }
