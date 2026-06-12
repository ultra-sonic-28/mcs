package spectrum

// Bus defines the interface for the memory and I/O bus of a ZX Spectrum.
type Bus interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
	Read16(addr uint16) uint16
	In(port uint16) uint8
	Out(port uint16, val uint8)
	
	// GetDisplayMemory returns the memory currently being used for display.
	GetDisplayMemory() []byte
	
	// Common components access (might be needed by Machine)
	GetKeyboard() *Keyboard
	GetTape() *Tape
	GetDisplay() *Display
	GetBorderColor() uint8
	GetTapeInState() bool
	SetTapeInState(state bool)
	
	// IsRom1Active returns true if the 48K BASIC ROM is currently paged in.
	IsRom1Active() bool
}
