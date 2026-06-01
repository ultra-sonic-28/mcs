package z80

// initExchange registers Exchange instructions.
func initExchange() {
	// 0x08: EX AF, AF'
	RegisterInstruction(&MainTable, 0x08, Instruction{
		Mnemonic:  "EX AF, AF'",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			cpu.Regs.ExchangeAF()
			return 4
		},
	})

	// 0xEB: EX DE, HL
	RegisterInstruction(&MainTable, 0xEB, Instruction{
		Mnemonic:  "EX DE, HL",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			oldDE := cpu.Regs.DE()
			cpu.Regs.SetDE(cpu.Regs.HL())
			cpu.Regs.SetHL(oldDE)
			return 4
		},
	})

	// 0xD9: EXX
	RegisterInstruction(&MainTable, 0xD9, Instruction{
		Mnemonic:  "EXX",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeNone,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.Regs.ExchangeMainSwaps()
			return 4
		},
	})

	// 0xE3: EX (SP), HL
	RegisterInstruction(&MainTable, 0xE3, Instruction{
		Mnemonic:  "EX (SP), HL",
		Length:    1,
		Cycles:    19,
		AddrMode1: AddrModeIndirect,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			low := cpu.Memory.Read(cpu.Regs.SP)
			high := cpu.Memory.Read(cpu.Regs.SP + 1)
			
			cpu.Memory.Write(cpu.Regs.SP, cpu.Regs.L)
			cpu.Memory.Write(cpu.Regs.SP+1, cpu.Regs.H)
			
			cpu.Regs.L = low
			cpu.Regs.H = high
			return 19
		},
	})

	// 0xDD 0xE3: EX (SP), IX
	RegisterInstruction(&DDTable, 0xE3, Instruction{
		Mnemonic:  "EX (SP), IX",
		Length:    2,
		Cycles:    23,
		AddrMode1: AddrModeIndirect,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			low := cpu.Memory.Read(cpu.Regs.SP)
			high := cpu.Memory.Read(cpu.Regs.SP + 1)
			
			cpu.Memory.Write(cpu.Regs.SP, uint8(cpu.Regs.IX&0xFF))
			cpu.Memory.Write(cpu.Regs.SP+1, uint8(cpu.Regs.IX>>8))
			
			cpu.Regs.IX = (uint16(high) << 8) | uint16(low)
			return 23
		},
	})

	// 0xFD 0xE3: EX (SP), IY
	RegisterInstruction(&FDTable, 0xE3, Instruction{
		Mnemonic:  "EX (SP), IY",
		Length:    2,
		Cycles:    23,
		AddrMode1: AddrModeIndirect,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			low := cpu.Memory.Read(cpu.Regs.SP)
			high := cpu.Memory.Read(cpu.Regs.SP + 1)
			
			cpu.Memory.Write(cpu.Regs.SP, uint8(cpu.Regs.IY&0xFF))
			cpu.Memory.Write(cpu.Regs.SP+1, uint8(cpu.Regs.IY>>8))
			
			cpu.Regs.IY = (uint16(high) << 8) | uint16(low)
			return 23
		},
	})
}
