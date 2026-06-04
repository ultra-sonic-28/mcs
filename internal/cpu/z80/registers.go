// Package z80 provides the implementation of the Zilog Z80 CPU.
package z80

import (
	"fmt"
	"log/slog"
)

// Registers represents the complete set of Z80 CPU registers.
type Registers struct {
	// Main Register Set
	A, F uint8
	B, C uint8
	D, E uint8
	H, L uint8

	// Alternate Register Set
	APrime, FPrime uint8
	BPrime, CPrime uint8
	DPrime, EPrime uint8
	HPrime, LPrime uint8

	// Index Registers
	IX, IY uint16

	// Special Purpose Registers
	PC, SP uint16
	I, R   uint8
}

// NewRegisters creates and initializes a new set of Z80 registers.
func NewRegisters() *Registers {
	slog.Info("initializing Z80 registers")
	return &Registers{}
}

// --- Index Register Accessors ---

// IXH returns the high byte of the IX register.
func (r *Registers) IXH() uint8 {
	return uint8(r.IX >> 8)
}

// SetIXH sets the high byte of the IX register.
func (r *Registers) SetIXH(val uint8) {
	r.IX = (uint16(val) << 8) | (r.IX & 0x00FF)
}

// IXL returns the low byte of the IX register.
func (r *Registers) IXL() uint8 {
	return uint8(r.IX & 0xFF)
}

// SetIXL sets the low byte of the IX register.
func (r *Registers) SetIXL(val uint8) {
	r.IX = (r.IX & 0xFF00) | uint16(val)
}

// IYH returns the high byte of the IY register.
func (r *Registers) IYH() uint8 {
	return uint8(r.IY >> 8)
}

// SetIYH sets the high byte of the IY register.
func (r *Registers) SetIYH(val uint8) {
	r.IY = (uint16(val) << 8) | (r.IY & 0x00FF)
}

// IYL returns the low byte of the IY register.
func (r *Registers) IYL() uint8 {
	return uint8(r.IY & 0xFF)
}

// SetIYL sets the low byte of the IY register.
func (r *Registers) SetIYL(val uint8) {
	r.IY = (r.IY & 0xFF00) | uint16(val)
}

// --- 16-bit Accessors ---

// AF returns the 16-bit value of the A and F register pair.
func (r *Registers) AF() uint16 {
	return uint16(r.A)<<8 | uint16(r.F)
}

// SetAF sets the 16-bit value of the A and F register pair.
func (r *Registers) SetAF(val uint16) {
	r.A = uint8(val >> 8)
	r.F = uint8(val & 0xFF)
}

// BC returns the 16-bit value of the B and C register pair.
func (r *Registers) BC() uint16 {
	return uint16(r.B)<<8 | uint16(r.C)
}

// SetBC sets the 16-bit value of the B and C register pair.
func (r *Registers) SetBC(val uint16) {
	r.B = uint8(val >> 8)
	r.C = uint8(val & 0xFF)
}

// DE returns the 16-bit value of the D and E register pair.
func (r *Registers) DE() uint16 {
	return uint16(r.D)<<8 | uint16(r.E)
}

// SetDE sets the 16-bit value of the D and E register pair.
func (r *Registers) SetDE(val uint16) {
	r.D = uint8(val >> 8)
	r.E = uint8(val & 0xFF)
}

// HL returns the 16-bit value of the H and L register pair.
func (r *Registers) HL() uint16 {
	return uint16(r.H)<<8 | uint16(r.L)
}

// SetHL sets the 16-bit value of the H and L register pair.
func (r *Registers) SetHL(val uint16) {
	r.H = uint8(val >> 8)
	r.L = uint8(val & 0xFF)
}

// --- Alternate Set Accessors ---

// AFPrime returns the 16-bit value of the A' and F' register pair.
func (r *Registers) AFPrime() uint16 {
	return uint16(r.APrime)<<8 | uint16(r.FPrime)
}

// SetAFPrime sets the 16-bit value of the A' and F' register pair.
func (r *Registers) SetAFPrime(val uint16) {
	r.APrime = uint8(val >> 8)
	r.FPrime = uint8(val & 0xFF)
}

// SetHLPrime sets the 16-bit value of the H' and L' register pair.
func (r *Registers) SetHLPrime(val uint16) {
	r.HPrime = uint8(val >> 8)
	r.LPrime = uint8(val & 0xFF)
}

// INC_HL increments the HL register pair.
func (r *Registers) INC_HL() {
	r.SetHL(r.HL() + 1)
}

// DEC_HL decrements the HL register pair.
func (r *Registers) DEC_HL() {
	r.SetHL(r.HL() - 1)
}

// INC_BC increments the BC register pair.
func (r *Registers) INC_BC() {
	r.SetBC(r.BC() + 1)
}

// DEC_BC decrements the BC register pair.
func (r *Registers) DEC_BC() {
	r.SetBC(r.BC() - 1)
}

// INC_DE increments the DE register pair.
func (r *Registers) INC_DE() {
	r.SetDE(r.DE() + 1)
}

// DEC_DE decrements the DE register pair.
func (r *Registers) DEC_DE() {
	r.SetDE(r.DE() - 1)
}

// --- Exchange Methods ---

// ExchangeAF swaps the AF and AF' register pairs.
func (r *Registers) ExchangeAF() {
	r.A, r.APrime = r.APrime, r.A
	r.F, r.FPrime = r.FPrime, r.F
}

// ExchangeMainSwaps swaps BC, DE, and HL with their alternate counterparts (EXX).
func (r *Registers) ExchangeMainSwaps() {
	r.B, r.BPrime = r.BPrime, r.B
	r.C, r.CPrime = r.CPrime, r.C
	r.D, r.DPrime = r.DPrime, r.D
	r.E, r.EPrime = r.EPrime, r.E
	r.H, r.HPrime = r.HPrime, r.H
	r.L, r.LPrime = r.LPrime, r.L
}

// --- Flag Helpers ---

// SetFlag sets or clears a specific flag bit.
func (r *Registers) SetFlag(flag uint8, value bool) {
	if value {
		r.F |= flag
	} else {
		r.F &= ^flag
	}
}

// Flag returns true if the specified flag bit is set.
func (r *Registers) Flag(flag uint8) bool {
	return (r.F & flag) != 0
}

// LogState logs the current state of all registers across three formatted lines.
func (r *Registers) LogState() {
	// Line 1: Main 8-bit registers
	slog.Info("Register Main",
		"A", fmt.Sprintf("0x%02X", r.A),
		"F", fmt.Sprintf("0x%02X", r.F),
		"B", fmt.Sprintf("0x%02X", r.B),
		"C", fmt.Sprintf("0x%02X", r.C),
		"D", fmt.Sprintf("0x%02X", r.D),
		"E", fmt.Sprintf("0x%02X", r.E),
		"H", fmt.Sprintf("0x%02X", r.H),
		"L", fmt.Sprintf("0x%02X", r.L),
	)

	// Line 2: Alternate 8-bit registers
	slog.Info("Register Alternate",
		"A'", fmt.Sprintf("0x%02X", r.APrime),
		"F'", fmt.Sprintf("0x%02X", r.FPrime),
		"B'", fmt.Sprintf("0x%02X", r.BPrime),
		"C'", fmt.Sprintf("0x%02X", r.CPrime),
		"D'", fmt.Sprintf("0x%02X", r.DPrime),
		"E'", fmt.Sprintf("0x%02X", r.EPrime),
		"H'", fmt.Sprintf("0x%02X", r.HPrime),
		"L'", fmt.Sprintf("0x%02X", r.LPrime),
	)

	// Line 3: 16-bit pairs and control registers
	slog.Info("Register State",
		"AF", fmt.Sprintf("0x%04X", r.AF()),
		"BC", fmt.Sprintf("0x%04X", r.BC()),
		"DE", fmt.Sprintf("0x%04X", r.DE()),
		"HL", fmt.Sprintf("0x%04X", r.HL()),
		"IX", fmt.Sprintf("0x%04X", r.IX),
		"IY", fmt.Sprintf("0x%04X", r.IY),
		"SP", fmt.Sprintf("0x%04X", r.SP),
		"PC", fmt.Sprintf("0x%04X", r.PC),
		"I", fmt.Sprintf("0x%02X", r.I),
		"R", fmt.Sprintf("0x%02X", r.R),
	)
}
