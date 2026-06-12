package memory

// Memory represents the memory system of a ZX Spectrum.
type Memory interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
	GetDisplayMemory() []byte
	IsRom1Active() bool
}
