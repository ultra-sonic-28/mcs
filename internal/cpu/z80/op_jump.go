package z80

// initJump registers Jump, Call, and Return instructions.
func initJump() {
	// 1. JP nn (Unconditional Jump)
	RegisterInstruction(&MainTable, 0xC3, Instruction{
		Mnemonic:  "JP nn",
		Length:    3,
		Cycles:    10,
		AddrMode1: AddrModeImmediate16,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			cpu.Regs.PC = nn
			return 10
		},
	})

	// 2. JP cc, nn (Conditional Jump)
	conditions := []struct {
		name string
		flag uint8
		val  bool
		op   uint8
	}{
		{"NZ", FlagZ, false, 0xC2},
		{"Z", FlagZ, true, 0xCA},
		{"NC", FlagC, false, 0xD2},
		{"C", FlagC, true, 0xDA},
		{"PO", FlagPV, false, 0xE2},
		{"PE", FlagPV, true, 0xEA},
		{"P", FlagS, false, 0xF2},
		{"M", FlagS, true, 0xFA},
	}

	for _, c := range conditions {
		c := c // capture
		RegisterInstruction(&MainTable, c.op, Instruction{
			Mnemonic:  "JP " + c.name + ", nn",
			Length:    3,
			Cycles:    10,
			AddrMode1: AddrModeImmediate16,
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				nn := cpu.FetchWord()
				if cpu.Regs.Flag(c.flag) == c.val {
					cpu.Regs.PC = nn
				}
				return 10
			},
		})
	}

	// 3. JR e (Unconditional Relative Jump)
	RegisterInstruction(&MainTable, 0x18, Instruction{
		Mnemonic:  "JR e",
		Length:    2,
		Cycles:    12,
		AddrMode1: AddrModeRelative,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			e := int8(cpu.FetchByte())
			cpu.Regs.PC = uint16(int32(cpu.Regs.PC) + int32(e))
			return 12
		},
	})

	// 4. JR cc, e (Conditional Relative Jump)
	jrConds := []struct {
		name string
		flag uint8
		val  bool
		op   uint8
	}{
		{"NZ", FlagZ, false, 0x20},
		{"Z", FlagZ, true, 0x28},
		{"NC", FlagC, false, 0x30},
		{"C", FlagC, true, 0x38},
	}

	for _, c := range jrConds {
		c := c // capture
		RegisterInstruction(&MainTable, c.op, Instruction{
			Mnemonic:  "JR " + c.name + ", e",
			Length:    2,
			Cycles:    7, // 7 if no jump, 12 if jump
			AddrMode1: AddrModeRelative,
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				e := int8(cpu.FetchByte())
				if cpu.Regs.Flag(c.flag) == c.val {
					cpu.Regs.PC = uint16(int32(cpu.Regs.PC) + int32(e))
					return 12
				}
				return 7
			},
		})
	}

	// 5. DJNZ e (Decrement B and Jump if Not Zero)
	RegisterInstruction(&MainTable, 0x10, Instruction{
		Mnemonic:  "DJNZ e",
		Length:    2,
		Cycles:    8, // 8 if no jump (B becomes 0), 13 if jump (B != 0)
		AddrMode1: AddrModeRelative,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			e := int8(cpu.FetchByte())
			cpu.Regs.B--
			if cpu.Regs.B != 0 {
				cpu.Regs.PC = uint16(int32(cpu.Regs.PC) + int32(e))
				return 13
			}
			return 8
		},
	})

	// 6. CALL nn (Unconditional Call)
	RegisterInstruction(&MainTable, 0xCD, Instruction{
		Mnemonic:  "CALL nn",
		Length:    3,
		Cycles:    17,
		AddrMode1: AddrModeImmediate16,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			// Push PC (which is already pointing to next instruction)
			cpu.Regs.SP--
			cpu.Memory.Write(cpu.Regs.SP, uint8(cpu.Regs.PC>>8))
			cpu.Regs.SP--
			cpu.Memory.Write(cpu.Regs.SP, uint8(cpu.Regs.PC&0xFF))
			cpu.Regs.PC = nn
			return 17
		},
	})

	// 6. CALL cc, nn (Conditional Call)
	for _, c := range conditions {
		c := c // capture
		RegisterInstruction(&MainTable, c.op+2, Instruction{ // Offset is +2 for non-ED conditional calls? Wait.
			// NZ: 0xC4, Z: 0xCC, NC: 0xD4, C: 0xDC, PO: 0xE4, PE: 0xEC, P: 0xF4, M: 0xFC
		})
	}
	// Let's do them explicitly to be safe with opcodes
	callConds := []struct {
		name string
		flag uint8
		val  bool
		op   uint8
	}{
		{"NZ", FlagZ, false, 0xC4},
		{"Z", FlagZ, true, 0xCC},
		{"NC", FlagC, false, 0xD4},
		{"C", FlagC, true, 0xDC},
		{"PO", FlagPV, false, 0xE4},
		{"PE", FlagPV, true, 0xEC},
		{"P", FlagS, false, 0xF4},
		{"M", FlagS, true, 0xFC},
	}
	for _, c := range callConds {
		c := c // capture
		RegisterInstruction(&MainTable, c.op, Instruction{
			Mnemonic:  "CALL " + c.name + ", nn",
			Length:    3,
			Cycles:    10, // 10 if no call, 17 if call
			AddrMode1: AddrModeImmediate16,
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				nn := cpu.FetchWord()
				if cpu.Regs.Flag(c.flag) == c.val {
					cpu.Regs.SP--
					cpu.Memory.Write(cpu.Regs.SP, uint8(cpu.Regs.PC>>8))
					cpu.Regs.SP--
					cpu.Memory.Write(cpu.Regs.SP, uint8(cpu.Regs.PC&0xFF))
					cpu.Regs.PC = nn
					return 17
				}
				return 10
			},
		})
	}

	// 7. RET (Unconditional Return)
	RegisterInstruction(&MainTable, 0xC9, Instruction{
		Mnemonic:  "RET",
		Length:    1,
		Cycles:    10,
		AddrMode1: AddrModeNone,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			low := uint16(cpu.Memory.Read(cpu.Regs.SP))
			cpu.Regs.SP++
			high := uint16(cpu.Memory.Read(cpu.Regs.SP))
			cpu.Regs.SP++
			cpu.Regs.PC = (high << 8) | low
			return 10
		},
	})

	// 8. RET cc (Conditional Return)
	retConds := []struct {
		name string
		flag uint8
		val  bool
		op   uint8
	}{
		{"NZ", FlagZ, false, 0xC0},
		{"Z", FlagZ, true, 0xC8},
		{"NC", FlagC, false, 0xD0},
		{"C", FlagC, true, 0xD8},
		{"PO", FlagPV, false, 0xE0},
		{"PE", FlagPV, true, 0xE8},
		{"P", FlagS, false, 0xF0},
		{"M", FlagS, true, 0xF8},
	}
	for _, c := range retConds {
		c := c // capture
		RegisterInstruction(&MainTable, c.op, Instruction{
			Mnemonic:  "RET " + c.name,
			Length:    1,
			Cycles:    5, // 5 if no ret, 11 if ret
			AddrMode1: AddrModeNone,
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				if cpu.Regs.Flag(c.flag) == c.val {
					low := uint16(cpu.Memory.Read(cpu.Regs.SP))
					cpu.Regs.SP++
					high := uint16(cpu.Memory.Read(cpu.Regs.SP))
					cpu.Regs.SP++
					cpu.Regs.PC = (high << 8) | low
					return 11
				}
				return 5
			},
		})
	}

	// 9. RST p (Restart)
	// Opcodes: 0xC7 (00), 0xCF (08), 0xD7 (10), 0xDF (18), 0xE7 (20), 0xEF (28), 0xF7 (30), 0xFF (38)
	rstTargets := []struct {
		name string
		val  uint16
		op   uint8
	}{
		{"00H", 0x0000, 0xC7},
		{"08H", 0x0008, 0xCF},
		{"10H", 0x0010, 0xD7},
		{"18H", 0x0018, 0xDF},
		{"20H", 0x0020, 0xE7},
		{"28H", 0x0028, 0xEF},
		{"30H", 0x0030, 0xF7},
		{"38H", 0x0038, 0xFF},
	}
	for _, rst := range rstTargets {
		rst := rst // capture
		RegisterInstruction(&MainTable, rst.op, Instruction{
			Mnemonic:  "RST " + rst.name,
			Length:    1,
			Cycles:    11,
			AddrMode1: AddrModeImmediate, // Fixed target
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				cpu.Regs.SP--
				cpu.Memory.Write(cpu.Regs.SP, uint8(cpu.Regs.PC>>8))
				cpu.Regs.SP--
				cpu.Memory.Write(cpu.Regs.SP, uint8(cpu.Regs.PC&0xFF))
				cpu.Regs.PC = rst.val
				return 11
			},
		})
	}

	// 10. JP (HL) (Unconditional Jump to HL)
	RegisterInstruction(&MainTable, 0xE9, Instruction{
		Mnemonic:  "JP (HL)",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeIndirect,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.Regs.PC = cpu.Regs.HL()
			return 4
		},
	})

	// 11. JP (IX) (Unconditional Jump to IX)
	RegisterInstruction(&DDTable, 0xE9, Instruction{
		Mnemonic:  "JP (IX)",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeIndirect,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.Regs.PC = cpu.Regs.IX
			return 8
		},
	})

	// 12. JP (IY) (Unconditional Jump to IY)
	RegisterInstruction(&FDTable, 0xE9, Instruction{
		Mnemonic:  "JP (IY)",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeIndirect,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.Regs.PC = cpu.Regs.IY
			return 8
		},
	})
}
