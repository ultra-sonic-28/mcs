package z80

// initMisc registers miscellaneous instructions like NOP and HALT.
func initMisc() {
	// 0x00: NOP (No Operation)
	RegisterInstruction(&MainTable, 0x00, Instruction{
		Mnemonic:  "NOP",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeImplied,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			return 4
		},
	})

	// 0x76: HALT
	RegisterInstruction(&MainTable, 0x76, Instruction{
		Mnemonic:  "HALT",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeImplied,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.SetHalt(true)
			return 4
		},
	})

	// 0x37: SCF (Set Carry Flag)
	RegisterInstruction(&MainTable, 0x37, Instruction{
		Mnemonic:  "SCF",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeImplied,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetFlag(FlagC, true)
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagN, false)
			// Undocumented flags 3 and 5 are copied from A
			cpu.Regs.SetFlag(Flag3, (cpu.Regs.A&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (cpu.Regs.A&0x20) != 0)
			return 4
		},
	})

	// 0x3F: CCF (Complement Carry Flag)
	RegisterInstruction(&MainTable, 0x3F, Instruction{
		Mnemonic:  "CCF",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeImplied,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			oldCarry := cpu.Regs.Flag(FlagC)
			cpu.Regs.SetFlag(FlagH, oldCarry)
			cpu.Regs.SetFlag(FlagC, !oldCarry)
			cpu.Regs.SetFlag(FlagN, false)
			// Undocumented flags 3 and 5 are copied from A
			cpu.Regs.SetFlag(Flag3, (cpu.Regs.A&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (cpu.Regs.A&0x20) != 0)
			return 4
		},
	})

	// 0xDD: NOP (IX Prefix treated as NOP in MainTable)
	RegisterInstruction(&MainTable, 0xDD, Instruction{
		Mnemonic:  "NOP",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeImplied,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			return 4
		},
	})

	// 0xFD: NOP (IY Prefix treated as NOP in MainTable)
	RegisterInstruction(&MainTable, 0xFD, Instruction{
		Mnemonic:  "NOP",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeImplied,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			return 4
		},
	})

	// 0xF3: DI (Disable Interrupts)
	RegisterInstruction(&MainTable, 0xF3, Instruction{
		Mnemonic:  "DI",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeImplied,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.IFF1 = false
			cpu.IFF2 = false
			return 4
		},
	})

	// 0xFB: EI (Enable Interrupts)
	RegisterInstruction(&MainTable, 0xFB, Instruction{
		Mnemonic:  "EI",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeImplied,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.IFF1 = true
			cpu.IFF2 = true
			return 4
		},
	})

	// --- ED Prefix Miscellaneous Instructions ---

	// 1. NEG (Negate Accumulator)
	// Opcode: ED 44
	// Cycles: 8 T-states
	RegisterInstruction(&EDTable, 0x44, Instruction{
		Mnemonic:  "NEG",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			oldA := cpu.Regs.A
			res := uint8(0) - oldA
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSub8(0, oldA, res)
			return 8
		},
	})

	// 2. IM 0, IM 1, IM 2
	imRegs := []struct {
		mode InterruptMode
		op   uint8
	}{
		{IM0, 0x46},
		{IM1, 0x56},
		{IM2, 0x5E},
		// Alternative opcodes for IM
		{IM0, 0x4E},
		{IM1, 0x66},
		{IM2, 0x6E},
		{IM0, 0x76},
		{IM2, 0x7E},
	}

	for _, im := range imRegs {
		im := im
		mnemonic := ""
		switch im.mode {
		case IM0: mnemonic = "IM 0"
		case IM1: mnemonic = "IM 1"
		case IM2: mnemonic = "IM 2"
		}
		RegisterInstruction(&EDTable, im.op, Instruction{
			Mnemonic:  mnemonic,
			Length:    2,
			Cycles:    8,
			AddrMode1: AddrModeImplied,
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				cpu.IM = im.mode
				return 8
			},
		})
	}

	// 3. RETI (Return from Interrupt)
	// Opcode: ED 4D
	// Cycles: 14 T-states
	RegisterInstruction(&EDTable, 0x4D, Instruction{
		Mnemonic:  "RETI",
		Length:    2,
		Cycles:    14,
		AddrMode1: AddrModeImplied,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			// Restore PC from stack
			low := uint16(cpu.Memory.Read(cpu.Regs.SP))
			cpu.Regs.SP++
			high := uint16(cpu.Memory.Read(cpu.Regs.SP))
			cpu.Regs.SP++
			cpu.Regs.PC = (high << 8) | low
			return 14
		},
	})

	// 4. RETN (Return from Non-Maskable Interrupt)
	// Opcode: ED 45
	// Cycles: 14 T-states
	RegisterInstruction(&EDTable, 0x45, Instruction{
		Mnemonic:  "RETN",
		Length:    2,
		Cycles:    14,
		AddrMode1: AddrModeImplied,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			// Restore PC from stack
			low := uint16(cpu.Memory.Read(cpu.Regs.SP))
			cpu.Regs.SP++
			high := uint16(cpu.Memory.Read(cpu.Regs.SP))
			cpu.Regs.SP++
			cpu.Regs.PC = (high << 8) | low
			// Restore IFF1 from IFF2
			cpu.IFF1 = cpu.IFF2
			return 14
		},
	})
}
