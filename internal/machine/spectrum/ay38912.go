package spectrum

import (
	"log/slog"
)

// AY38912 represents the AY-3-8912 Programmable Sound Generator.
type AY38912 struct {
	Registers [16]uint8
	SelectedRegister uint8
}

// NewAY38912 creates a new AY-3-8912 instance.
func NewAY38912() *AY38912 {
	return &AY38912{}
}

// WriteAddress selects the register for subsequent read/write operations.
func (ay *AY38912) WriteAddress(val uint8) {
	ay.SelectedRegister = val & 0x0F
}

// WriteData writes a value to the currently selected register.
func (ay *AY38912) WriteData(val uint8) {
	ay.Registers[ay.SelectedRegister] = val
	slog.Debug("AY-3-8912 register write", "reg", ay.SelectedRegister, "val", val)
}

// ReadData returns the value of the currently selected register.
func (ay *AY38912) ReadData() uint8 {
	return ay.Registers[ay.SelectedRegister]
}
