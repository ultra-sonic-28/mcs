package spectrum

import (
	"log/slog"
)

// Bus128 implements the memory and I/O bus for the ZX Spectrum 128K.
type Bus128 struct {
	BaseBus
	Memory *Memory128
	AY     *AY38912
}

// NewBus128 creates a new Spectrum 128K Bus with the embedded ROMs loaded.
func NewBus128() *Bus128 {
	slog.Info("Initializing Spectrum 128K Bus")
	b := &Bus128{
		BaseBus: BaseBus{
			Keyboard: NewKeyboard(),
			Display:  NewDisplay(),
			Tape:     NewTape(),
		},
		Memory: NewMemory128(),
		AY:     NewAY38912(),
	}

	return b
}

// Read returns the byte at the specified memory address.
func (b *Bus128) Read(addr uint16) uint8 {
	return b.Memory.Read(addr)
}

// Read16 returns the 16-bit word at the specified memory address (little-endian).
func (b *Bus128) Read16(addr uint16) uint16 {
	return uint16(b.Read(addr)) | (uint16(b.Read(addr+1)) << 8)
}

// Write stores a byte at the specified memory address.
func (b *Bus128) Write(addr uint16, val uint8) {
	b.Memory.Write(addr, val)
}

// In reads a byte from the specified I/O port.
func (b *Bus128) In(port uint16) uint8 {
	// Port 0xFE: Keyboard and Tape
	if port&0x01 == 0 {
		mask := uint8(port >> 8)
		result := b.Keyboard.Scan(mask)

		if b.TapeInState {
			result |= 0x40
		} else {
			result &= ^uint8(0x40)
		}

		result |= 0xA0
		return result
	}

	// Port 0xFFFD: AY Register Read (bit 1=0, bit 14=1, bit 15=1)
	if (port & 0xC002) == 0xC000 {
		return b.AY.ReadData()
	}

	return 0xFF
}

// Out writes a byte to the specified I/O port.
func (b *Bus128) Out(port uint16, val uint8) {
	// Port 0xFE: Border, Beeper, Mic
	if port&0x01 == 0 {
		b.BorderColor = val & 0x07
		b.MicState = (val & 0x08) != 0
		b.BeeperState = (val & 0x10) != 0
		return
	}

	// Port 0x7FFD: Memory Paging (bit 1=0, bit 15=0)
	if (port & 0x8002) == 0 {
		b.Memory.Page(val)
		return
	}

	// Port 0xFFFD: AY Register Select (bit 1=0, bit 14=1, bit 15=1)
	if (port & 0xC002) == 0xC000 {
		b.AY.WriteAddress(val)
		return
	}

	// Port 0xBFFD: AY Data Write (bit 1=0, bit 14=0, bit 15=1)
	if (port & 0xC002) == 0x8000 {
		b.AY.WriteData(val)
		return
	}
}

// GetDisplayMemory returns the 6912 bytes of memory used for the display.
func (b *Bus128) GetDisplayMemory() []byte {
	return b.Memory.GetDisplayMemory()
}

func (b *Bus128) IsRom1Active() bool {
	return b.Memory.IsRom1Active()
}
