// Package spectrum implements the ZX Spectrum machine logic.
package spectrum

// BaseBus contains the common components and state for all Spectrum models.
type BaseBus struct {
	Keyboard *Keyboard
	Display  *Display
	Tape     *Tape

	// ULA State
	BorderColor uint8
	BeeperState bool
	MicState    bool
	TapeInState bool
}

func (b *BaseBus) GetKeyboard() *Keyboard      { return b.Keyboard }
func (b *BaseBus) GetTape() *Tape              { return b.Tape }
func (b *BaseBus) GetDisplay() *Display        { return b.Display }
func (b *BaseBus) GetBorderColor() uint8       { return b.BorderColor }
func (b *BaseBus) GetTapeInState() bool        { return b.TapeInState }
func (b *BaseBus) SetTapeInState(state bool)   { b.TapeInState = state }
