package z80

// initBlock registers Block Transfer and Search instructions.
func initBlock() {
	// --- Block Transfer Instructions ---

	// 1. LDI (Load and Increment)
	// Opcode: ED A0
	// Cycles: 16 T-states
	RegisterInstruction(&EDTable, 0xA0, Instruction{
		Mnemonic:  "LDI",
		Length:    2,
		Cycles:    16,
		AddrMode1: AddrModeNone,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			cpu.Memory.Write(cpu.Regs.DE(), val)
			cpu.Regs.INC_HL()
			cpu.Regs.INC_DE()
			cpu.Regs.DEC_BC()

			// Flags
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagPV, cpu.Regs.BC() != 0)
			cpu.Regs.SetFlag(FlagN, false)
			
			// Undocumented flags 3 and 5 are based on (A + val)
			res := val + cpu.Regs.A
			cpu.Regs.SetFlag(Flag3, (res&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (res&0x02) != 0) // Bit 1 of (A+val) for LDI/LDD

			return 16
		},
	})

	// 2. LDIR (Load, Increment and Repeat)
	// Opcode: ED B0
	// Cycles: 21 T-states if repeat, 16 T-states if done
	RegisterInstruction(&EDTable, 0xB0, Instruction{
		Mnemonic:  "LDIR",
		Length:    2,
		Cycles:    21, // Repeat case
		AddrMode1: AddrModeNone,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			cpu.Memory.Write(cpu.Regs.DE(), val)
			cpu.Regs.INC_HL()
			cpu.Regs.INC_DE()
			cpu.Regs.DEC_BC()

			// Flags
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagN, false)
			
			// Undocumented flags
			res := val + cpu.Regs.A
			cpu.Regs.SetFlag(Flag3, (res&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (res&0x02) != 0)

			if cpu.Regs.BC() != 0 {
				cpu.Regs.PC -= 2 // Repeat
				cpu.Regs.SetFlag(FlagPV, true)
				return 21
			}
			cpu.Regs.SetFlag(FlagPV, false)
			return 16
		},
	})

	// 3. LDD (Load and Decrement)
	// Opcode: ED A8
	// Cycles: 16 T-states
	RegisterInstruction(&EDTable, 0xA8, Instruction{
		Mnemonic:  "LDD",
		Length:    2,
		Cycles:    16,
		AddrMode1: AddrModeNone,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			cpu.Memory.Write(cpu.Regs.DE(), val)
			cpu.Regs.DEC_HL()
			cpu.Regs.DEC_DE()
			cpu.Regs.DEC_BC()

			// Flags
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagPV, cpu.Regs.BC() != 0)
			cpu.Regs.SetFlag(FlagN, false)
			
			res := val + cpu.Regs.A
			cpu.Regs.SetFlag(Flag3, (res&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (res&0x02) != 0)

			return 16
		},
	})

	// 4. LDDR (Load, Decrement and Repeat)
	// Opcode: ED B8
	// Cycles: 21 T-states if repeat, 16 T-states if done
	RegisterInstruction(&EDTable, 0xB8, Instruction{
		Mnemonic:  "LDDR",
		Length:    2,
		Cycles:    21,
		AddrMode1: AddrModeNone,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			cpu.Memory.Write(cpu.Regs.DE(), val)
			cpu.Regs.DEC_HL()
			cpu.Regs.DEC_DE()
			cpu.Regs.DEC_BC()

			// Flags
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagN, false)

			// Undocumented flags
			res := val + cpu.Regs.A
			cpu.Regs.SetFlag(Flag3, (res&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (res&0x02) != 0)

			if cpu.Regs.BC() != 0 {
				cpu.Regs.PC -= 2 // Repeat
				cpu.Regs.SetFlag(FlagPV, true)
				return 21
			}
			cpu.Regs.SetFlag(FlagPV, false)
			return 16
		},
	})

	// --- Block Search Instructions ---

	// 5. CPI (Compare and Increment)
	// Opcode: ED A1
	RegisterInstruction(&EDTable, 0xA1, Instruction{
		Mnemonic:  "CPI",
		Length:    2,
		Cycles:    16,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.INC_HL()
			cpu.Regs.DEC_BC()

			// Flags
			cpu.Regs.SetFlag(FlagS, (res&0x80) != 0)
			cpu.Regs.SetFlag(FlagZ, res == 0)
			cpu.Regs.SetFlag(FlagH, (oldA&0x0F) < (val&0x0F))
			cpu.Regs.SetFlag(FlagPV, cpu.Regs.BC() != 0)
			cpu.Regs.SetFlag(FlagN, true)
			
			// Undocumented flags for CPI/CPD:
			// F5 = bit 1 of (A - val - H)
			// F3 = bit 3 of (A - val - H)
			var h uint8
			if cpu.Regs.Flag(FlagH) { h = 1 }
			resUndoc := oldA - val - h
			cpu.Regs.SetFlag(Flag5, (resUndoc&0x02) != 0)
			cpu.Regs.SetFlag(Flag3, (resUndoc&0x08) != 0)

			return 16
		},
	})

	// 6. CPIR (Compare, Increment and Repeat)
	// Opcode: ED B1
	RegisterInstruction(&EDTable, 0xB1, Instruction{
		Mnemonic:  "CPIR",
		Length:    2,
		Cycles:    21,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.INC_HL()
			cpu.Regs.DEC_BC()

			// Flags
			cpu.Regs.SetFlag(FlagS, (res&0x80) != 0)
			cpu.Regs.SetFlag(FlagZ, res == 0)
			cpu.Regs.SetFlag(FlagH, (oldA&0x0F) < (val&0x0F))
			cpu.Regs.SetFlag(FlagPV, cpu.Regs.BC() != 0)
			cpu.Regs.SetFlag(FlagN, true)
			
			var h uint8
			if cpu.Regs.Flag(FlagH) { h = 1 }
			resUndoc := oldA - val - h
			cpu.Regs.SetFlag(Flag5, (resUndoc&0x02) != 0)
			cpu.Regs.SetFlag(Flag3, (resUndoc&0x08) != 0)

			if cpu.Regs.BC() != 0 && res != 0 {
				cpu.Regs.PC -= 2 // Repeat
				return 21
			}
			return 16
		},
	})

	// 7. CPD (Compare and Decrement)
	// Opcode: ED A9
	RegisterInstruction(&EDTable, 0xA9, Instruction{
		Mnemonic:  "CPD",
		Length:    2,
		Cycles:    16,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.DEC_HL()
			cpu.Regs.DEC_BC()

			// Flags
			cpu.Regs.SetFlag(FlagS, (res&0x80) != 0)
			cpu.Regs.SetFlag(FlagZ, res == 0)
			cpu.Regs.SetFlag(FlagH, (oldA&0x0F) < (val&0x0F))
			cpu.Regs.SetFlag(FlagPV, cpu.Regs.BC() != 0)
			cpu.Regs.SetFlag(FlagN, true)
			
			var h uint8
			if cpu.Regs.Flag(FlagH) { h = 1 }
			resUndoc := oldA - val - h
			cpu.Regs.SetFlag(Flag5, (resUndoc&0x02) != 0)
			cpu.Regs.SetFlag(Flag3, (resUndoc&0x08) != 0)

			return 16
		},
	})

	// 8. CPDR (Compare, Decrement and Repeat)
	// Opcode: ED B9
	RegisterInstruction(&EDTable, 0xB9, Instruction{
		Mnemonic:  "CPDR",
		Length:    2,
		Cycles:    21,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.DEC_HL()
			cpu.Regs.DEC_BC()

			// Flags
			cpu.Regs.SetFlag(FlagS, (res&0x80) != 0)
			cpu.Regs.SetFlag(FlagZ, res == 0)
			cpu.Regs.SetFlag(FlagH, (oldA&0x0F) < (val&0x0F))
			cpu.Regs.SetFlag(FlagPV, cpu.Regs.BC() != 0)
			cpu.Regs.SetFlag(FlagN, true)

			var h uint8
			if cpu.Regs.Flag(FlagH) { h = 1 }
			resUndoc := oldA - val - h
			cpu.Regs.SetFlag(Flag5, (resUndoc&0x02) != 0)
			cpu.Regs.SetFlag(Flag3, (resUndoc&0x08) != 0)

			if cpu.Regs.BC() != 0 && res != 0 {
				cpu.Regs.PC -= 2 // Repeat
				return 21
			}
			return 16
		},
	})

	// --- Block I/O Instructions ---

	// 9. INI / INIR / IND / INDR
	RegisterInstruction(&EDTable, 0xA2, Instruction{
		Mnemonic: "INI",
		Length:   2,
		Cycles:   16,
		Execute: func(cpu *CPU) int {
			val := cpu.IO.In(cpu.Regs.BC())
			cpu.Memory.Write(cpu.Regs.HL(), val)
			cpu.Regs.INC_HL()
			cpu.Regs.B--
			cpu.Regs.SetFlag(FlagZ, cpu.Regs.B == 0)
			cpu.Regs.SetFlag(FlagN, true)
			return 16
		},
	})
	RegisterInstruction(&EDTable, 0xB2, Instruction{
		Mnemonic: "INIR",
		Length:   2,
		Cycles:   21,
		Execute: func(cpu *CPU) int {
			val := cpu.IO.In(cpu.Regs.BC())
			cpu.Memory.Write(cpu.Regs.HL(), val)
			cpu.Regs.INC_HL()
			cpu.Regs.B--
			cpu.Regs.SetFlag(FlagZ, true) // Always set if it ends
			cpu.Regs.SetFlag(FlagN, true)
			if cpu.Regs.B != 0 {
				cpu.Regs.PC -= 2
				cpu.Regs.SetFlag(FlagZ, false)
				return 21
			}
			return 16
		},
	})
	RegisterInstruction(&EDTable, 0xAA, Instruction{
		Mnemonic: "IND",
		Length:   2,
		Cycles:   16,
		Execute: func(cpu *CPU) int {
			val := cpu.IO.In(cpu.Regs.BC())
			cpu.Memory.Write(cpu.Regs.HL(), val)
			cpu.Regs.DEC_HL()
			cpu.Regs.B--
			cpu.Regs.SetFlag(FlagZ, cpu.Regs.B == 0)
			cpu.Regs.SetFlag(FlagN, true)
			return 16
		},
	})
	RegisterInstruction(&EDTable, 0xBA, Instruction{
		Mnemonic: "INDR",
		Length:   2,
		Cycles:   21,
		Execute: func(cpu *CPU) int {
			val := cpu.IO.In(cpu.Regs.BC())
			cpu.Memory.Write(cpu.Regs.HL(), val)
			cpu.Regs.DEC_HL()
			cpu.Regs.B--
			cpu.Regs.SetFlag(FlagZ, true)
			cpu.Regs.SetFlag(FlagN, true)
			if cpu.Regs.B != 0 {
				cpu.Regs.PC -= 2
				cpu.Regs.SetFlag(FlagZ, false)
				return 21
			}
			return 16
		},
	})

	// 10. OUTI / OTIR / OUTD / OTDR
	RegisterInstruction(&EDTable, 0xA3, Instruction{
		Mnemonic: "OUTI",
		Length:   2,
		Cycles:   16,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			cpu.Regs.B--
			cpu.IO.Out(cpu.Regs.BC(), val)
			cpu.Regs.INC_HL()
			cpu.Regs.SetFlag(FlagZ, cpu.Regs.B == 0)
			cpu.Regs.SetFlag(FlagN, true)
			return 16
		},
	})
	RegisterInstruction(&EDTable, 0xB3, Instruction{
		Mnemonic: "OTIR",
		Length:   2,
		Cycles:   21,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			cpu.Regs.B--
			cpu.IO.Out(cpu.Regs.BC(), val)
			cpu.Regs.INC_HL()
			cpu.Regs.SetFlag(FlagZ, true)
			cpu.Regs.SetFlag(FlagN, true)
			if cpu.Regs.B != 0 {
				cpu.Regs.PC -= 2
				cpu.Regs.SetFlag(FlagZ, false)
				return 21
			}
			return 16
		},
	})
	RegisterInstruction(&EDTable, 0xAB, Instruction{
		Mnemonic: "OUTD",
		Length:   2,
		Cycles:   16,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			cpu.Regs.B--
			cpu.IO.Out(cpu.Regs.BC(), val)
			cpu.Regs.DEC_HL()
			cpu.Regs.SetFlag(FlagZ, cpu.Regs.B == 0)
			cpu.Regs.SetFlag(FlagN, true)
			return 16
		},
	})
	RegisterInstruction(&EDTable, 0xBB, Instruction{
		Mnemonic: "OTDR",
		Length:   2,
		Cycles:   21,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			cpu.Regs.B--
			cpu.IO.Out(cpu.Regs.BC(), val)
			cpu.Regs.DEC_HL()
			cpu.Regs.SetFlag(FlagZ, true)
			cpu.Regs.SetFlag(FlagN, true)
			if cpu.Regs.B != 0 {
				cpu.Regs.PC -= 2
				cpu.Regs.SetFlag(FlagZ, false)
				return 21
			}
			return 16
		},
	})
}
