// Package z80 provides the implementation of the Zilog Z80 CPU.
package z80

// Flags bit positions in the F register.
const (
	FlagC  uint8 = 1 << 0 // Carry Flag
	FlagN  uint8 = 1 << 1 // Add/Subtract Flag
	FlagPV uint8 = 1 << 2 // Parity/Overflow Flag
	Flag3  uint8 = 1 << 3 // Undocumented Bit 3
	FlagH  uint8 = 1 << 4 // Half-Carry Flag
	Flag5  uint8 = 1 << 5 // Undocumented Bit 5
	FlagZ  uint8 = 1 << 6 // Zero Flag
	FlagS  uint8 = 1 << 7 // Sign Flag
)

// UpdateFlagsAdd8 updates the flags for an 8-bit addition (A = A + n).
func (r *Registers) UpdateFlagsAdd8(oldA, val, res uint8) {
	// Carry Flag (C): Set if there's a carry from bit 7
	r.SetFlag(FlagC, uint16(oldA)+uint16(val) > 0xFF)

	// Add/Subtract Flag (N): Cleared for addition
	r.SetFlag(FlagN, false)

	// Parity/Overflow Flag (PV): Set if 2's complement overflow
	// Overflow occurs if signs of operands are same, but sign of result is different.
	overflow := ((oldA ^ res) & (val ^ res) & 0x80) != 0
	r.SetFlag(FlagPV, overflow)

	// Half-Carry Flag (H): Set if there's a carry from bit 3
	halfCarry := ((oldA & 0x0F) + (val & 0x0F)) > 0x0F
	r.SetFlag(FlagH, halfCarry)

	// Zero Flag (Z): Set if result is zero
	r.SetFlag(FlagZ, res == 0)

	// Sign Flag (S): Set if bit 7 of result is set
	r.SetFlag(FlagS, (res&0x80) != 0)

	// Undocumented flags (3 and 5): Copied from the result bits 3 and 5
	r.SetFlag(Flag3, (res&0x08) != 0)
	r.SetFlag(Flag5, (res&0x20) != 0)
}

// UpdateFlagsIOIn updates the flags for an 8-bit I/O input instruction (IN r, (C)).
func (r *Registers) UpdateFlagsIOIn(res uint8) {
	// Sign Flag (S): Set if bit 7 of result is set
	r.SetFlag(FlagS, (res&0x80) != 0)

	// Zero Flag (Z): Set if result is zero
	r.SetFlag(FlagZ, res == 0)

	// Half-Carry Flag (H): Cleared
	r.SetFlag(FlagH, false)

	// Parity/Overflow Flag (PV): Set if parity is even
	parity := true
	for i := 0; i < 8; i++ {
		if (res & (1 << i)) != 0 {
			parity = !parity
		}
	}
	r.SetFlag(FlagPV, parity)

	// Add/Subtract Flag (N): Cleared
	r.SetFlag(FlagN, false)

	// Carry Flag (C): Unaffected

	// Undocumented flags (3 and 5): Copied from the result bits 3 and 5
	r.SetFlag(Flag3, (res&0x08) != 0)
	r.SetFlag(Flag5, (res&0x20) != 0)
}
func (r *Registers) UpdateFlagsAdc8(oldA, val, res uint8, carry uint8) {
	// Carry Flag (C): Set if there's a carry from bit 7
	r.SetFlag(FlagC, uint16(oldA)+uint16(val)+uint16(carry) > 0xFF)

	// Add/Subtract Flag (N): Cleared for addition
	r.SetFlag(FlagN, false)

	// Parity/Overflow Flag (PV): Set if 2's complement overflow
	overflow := ((oldA ^ res) & (val ^ res) & 0x80) != 0
	r.SetFlag(FlagPV, overflow)

	// Half-Carry Flag (H): Set if there's a carry from bit 3
	halfCarry := ((oldA & 0x0F) + (val & 0x0F) + carry) > 0x0F
	r.SetFlag(FlagH, halfCarry)

	// Zero Flag (Z): Set if result is zero
	r.SetFlag(FlagZ, res == 0)

	// Sign Flag (S): Set if bit 7 of result is set
	r.SetFlag(FlagS, (res&0x80) != 0)

	// Undocumented flags (3 and 5): Copied from the result bits 3 and 5
	r.SetFlag(Flag3, (res&0x08) != 0)
	r.SetFlag(Flag5, (res&0x20) != 0)
}

// UpdateFlagsSub8 updates the flags for an 8-bit subtraction (A = A - n).
func (r *Registers) UpdateFlagsSub8(oldA, val, res uint8) {
	// Carry Flag (C): Set if there's a borrow from bit 7 (val > oldA)
	r.SetFlag(FlagC, uint16(val) > uint16(oldA))

	// Add/Subtract Flag (N): Set for subtraction
	r.SetFlag(FlagN, true)

	// Parity/Overflow Flag (PV): Set if 2's complement overflow
	overflow := ((oldA ^ val) & (oldA ^ res) & 0x80) != 0
	r.SetFlag(FlagPV, overflow)

	// Half-Carry Flag (H): Set if there's a borrow from bit 4
	halfCarry := (oldA & 0x0F) < (val & 0x0F)
	r.SetFlag(FlagH, halfCarry)

	// Zero Flag (Z): Set if result is zero
	r.SetFlag(FlagZ, res == 0)

	// Sign Flag (S): Set if bit 7 of result is set
	r.SetFlag(FlagS, (res&0x80) != 0)

	// Undocumented flags (3 and 5): Copied from the result bits 3 and 5
	r.SetFlag(Flag3, (res&0x08) != 0)
	r.SetFlag(Flag5, (res&0x20) != 0)
}

// UpdateFlagsSbc8 updates the flags for an 8-bit subtraction with borrow (A = A - n - CY).
func (r *Registers) UpdateFlagsSbc8(oldA, val, res uint8, carry uint8) {
	// Carry Flag (C): Set if there's a borrow from bit 7
	r.SetFlag(FlagC, uint16(val)+uint16(carry) > uint16(oldA))

	// Add/Subtract Flag (N): Set for subtraction
	r.SetFlag(FlagN, true)

	// Parity/Overflow Flag (PV): Set if 2's complement overflow
	overflow := ((oldA ^ val) & (oldA ^ res) & 0x80) != 0
	r.SetFlag(FlagPV, overflow)

	// Half-Carry Flag (H): Set if there's a borrow from bit 4
	halfCarry := (uint16(oldA & 0x0F)) < (uint16(val & 0x0F) + uint16(carry))
	r.SetFlag(FlagH, halfCarry)

	// Zero Flag (Z): Set if result is zero
	r.SetFlag(FlagZ, res == 0)

	// Sign Flag (S): Set if bit 7 of result is set
	r.SetFlag(FlagS, (res&0x80) != 0)

	// Undocumented flags (3 and 5): Copied from the result bits 3 and 5
	r.SetFlag(Flag3, (res&0x08) != 0)
	r.SetFlag(Flag5, (res&0x20) != 0)
}

// UpdateFlagsCp8 updates the flags for an 8-bit comparison (A - n).
// Undocumented flags 3 and 5 are copied from the operand, not the result.
func (r *Registers) UpdateFlagsCp8(oldA, val, res uint8) {
	// Carry Flag (C): Set if there's a borrow from bit 7
	r.SetFlag(FlagC, uint16(val) > uint16(oldA))

	// Add/Subtract Flag (N): Set for subtraction
	r.SetFlag(FlagN, true)

	// Parity/Overflow Flag (PV): Set if 2's complement overflow
	overflow := ((oldA ^ val) & (oldA ^ res) & 0x80) != 0
	r.SetFlag(FlagPV, overflow)

	// Half-Carry Flag (H): Set if there's a borrow from bit 4
	halfCarry := (oldA & 0x0F) < (val & 0x0F)
	r.SetFlag(FlagH, halfCarry)

	// Zero Flag (Z): Set if result is zero
	r.SetFlag(FlagZ, res == 0)

	// Sign Flag (S): Set if bit 7 of result is set
	r.SetFlag(FlagS, (res&0x80) != 0)

	// Undocumented flags (3 and 5): Copied from the OPERAND bits 3 and 5
	r.SetFlag(Flag3, (val&0x08) != 0)
	r.SetFlag(Flag5, (val&0x20) != 0)
}

// UpdateFlagsLogical8 updates the flags for 8-bit logical operations (AND, OR, XOR).
func (r *Registers) UpdateFlagsLogical8(res uint8, isAnd bool) {
	// Carry Flag (C): Cleared
	r.SetFlag(FlagC, false)

	// Add/Subtract Flag (N): Cleared
	r.SetFlag(FlagN, false)

	// Parity/Overflow Flag (PV): Set if parity is even
	parity := true
	for i := 0; i < 8; i++ {
		if (res & (1 << i)) != 0 {
			parity = !parity
		}
	}
	r.SetFlag(FlagPV, parity)

	// Half-Carry Flag (H): Set for AND, Cleared for OR/XOR
	r.SetFlag(FlagH, isAnd)

	// Zero Flag (Z): Set if result is zero
	r.SetFlag(FlagZ, res == 0)

	// Sign Flag (S): Set if bit 7 of result is set
	r.SetFlag(FlagS, (res&0x80) != 0)

	// Undocumented flags (3 and 5): Copied from the result bits 3 and 5
	r.SetFlag(Flag3, (res&0x08) != 0)
	r.SetFlag(Flag5, (res&0x20) != 0)
}

// UpdateFlagsAdd16 updates the flags for a 16-bit addition (HL = HL + rr).
// Note: S, Z, and PV flags are not affected.
func (r *Registers) UpdateFlagsAdd16(oldVal, val uint16, res uint32) {
	// Carry Flag (C): Set if there's a carry from bit 15
	r.SetFlag(FlagC, res > 0xFFFF)

	// Add/Subtract Flag (N): Cleared for addition
	r.SetFlag(FlagN, false)

	// Half-Carry Flag (H): Set if there's a carry from bit 11
	r.SetFlag(FlagH, ((oldVal&0x0FFF)+(val&0x0FFF)) > 0x0FFF)

	// Undocumented flags (3 and 5): Copied from bits 11 and 13 of the result
	res16 := uint16(res)
	r.SetFlag(Flag3, (res16&0x0800) != 0)
	r.SetFlag(Flag5, (res16&0x2000) != 0)
}

// UpdateFlagsAdc16 updates the flags for a 16-bit addition with carry (HL = HL + rr + CY).
func (r *Registers) UpdateFlagsAdc16(oldHL, val, res uint32, carry uint32) {
	res16 := uint16(res)

	// Carry Flag (C): Set if there's a carry from bit 15
	r.SetFlag(FlagC, res > 0xFFFF)

	// Add/Subtract Flag (N): Cleared for addition
	r.SetFlag(FlagN, false)

	// Parity/Overflow Flag (PV): Set if 2's complement overflow
	overflow := ((oldHL ^ res) & ((val + carry) ^ res) & 0x8000) != 0
	r.SetFlag(FlagPV, overflow)

	// Half-Carry Flag (H): Set if there's a carry from bit 11
	halfCarry := ((oldHL & 0x0FFF) + (val & 0x0FFF) + carry) > 0x0FFF
	r.SetFlag(FlagH, halfCarry)

	// Zero Flag (Z): Set if result is zero
	r.SetFlag(FlagZ, res16 == 0)

	// Sign Flag (S): Set if bit 15 of result is set
	r.SetFlag(FlagS, (res16&0x8000) != 0)

	// Undocumented flags (3 and 5): Copied from bits 11 and 13 of the result
	r.SetFlag(Flag3, (res16&0x0800) != 0)
	r.SetFlag(Flag5, (res16&0x2000) != 0)
}

// UpdateFlagsSbc16 updates the flags for a 16-bit subtraction with borrow (HL = HL - rr - CY).
func (r *Registers) UpdateFlagsSbc16(oldHL, val uint16, res uint32, carry uint32) {
	res16 := uint16(res)

	// Carry Flag (C): Set if there's a borrow from bit 15
	r.SetFlag(FlagC, uint32(oldHL) < (uint32(val) + carry))

	// Add/Subtract Flag (N): Set for subtraction
	r.SetFlag(FlagN, true)

	// Parity/Overflow Flag (PV): Set if 2's complement overflow
	overflow := ((oldHL ^ val) & (oldHL ^ res16) & 0x8000) != 0
	r.SetFlag(FlagPV, overflow)

	// Half-Carry Flag (H): Set if there's a borrow from bit 12
	r.SetFlag(FlagH, (oldHL&0x0FFF) < ((val&0x0FFF)+uint16(carry)))

	// Zero Flag (Z): Set if result is zero
	r.SetFlag(FlagZ, res16 == 0)

	// Sign Flag (S): Set if bit 15 of result is set
	r.SetFlag(FlagS, (res16&0x8000) != 0)

	// Undocumented flags (3 and 5): Copied from bits 11 and 13 of the result
	r.SetFlag(Flag3, (res16&0x0800) != 0)
	r.SetFlag(Flag5, (res16&0x2000) != 0)
}

// UpdateFlagsInc8 updates the flags for an 8-bit increment (r = r + 1).
// Note: Carry flag is not affected.
func (r *Registers) UpdateFlagsInc8(oldVal, res uint8) {
	// Add/Subtract Flag (N): Cleared for addition
	r.SetFlag(FlagN, false)

	// Parity/Overflow Flag (PV): Set if result overflowed (0x7F -> 0x80)
	r.SetFlag(FlagPV, oldVal == 0x7F)

	// Half-Carry Flag (H): Set if there's a carry from bit 3
	r.SetFlag(FlagH, (oldVal&0x0F) == 0x0F)

	// Zero Flag (Z): Set if result is zero
	r.SetFlag(FlagZ, res == 0)

	// Sign Flag (S): Set if bit 7 of result is set
	r.SetFlag(FlagS, (res&0x80) != 0)

	// Undocumented flags (3 and 5): Copied from the result bits 3 and 5
	r.SetFlag(Flag3, (res&0x08) != 0)
	r.SetFlag(Flag5, (res&0x20) != 0)
}

// UpdateFlagsDec8 updates the flags for an 8-bit decrement (r = r - 1).
// Note: Carry flag is not affected.
func (r *Registers) UpdateFlagsDec8(oldVal, res uint8) {
	// Add/Subtract Flag (N): Set for subtraction
	r.SetFlag(FlagN, true)

	// Parity/Overflow Flag (PV): Set if result overflowed (0x80 -> 0x7F)
	r.SetFlag(FlagPV, oldVal == 0x80)

	// Half-Carry Flag (H): Set if there's a borrow from bit 4
	r.SetFlag(FlagH, (oldVal&0x0F) == 0x00)

	// Zero Flag (Z): Set if result is zero
	r.SetFlag(FlagZ, res == 0)

	// Sign Flag (S): Set if bit 7 of result is set
	r.SetFlag(FlagS, (res&0x80) != 0)

	// Undocumented flags (3 and 5): Copied from the result bits 3 and 5
	r.SetFlag(Flag3, (res&0x08) != 0)
	r.SetFlag(Flag5, (res&0x20) != 0)
}
