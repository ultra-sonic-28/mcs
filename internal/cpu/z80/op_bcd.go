package z80

// initBCD registers BCD-related instructions.
func initBCD() {
	// 0x27: DAA (Decimal Adjust Accumulator)
	RegisterInstruction(&MainTable, 0x27, Instruction{
		Mnemonic:  "DAA",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			a := cpu.Regs.A
			n := cpu.Regs.Flag(FlagN)
			c := cpu.Regs.Flag(FlagC)
			h := cpu.Regs.Flag(FlagH)

			correction := uint8(0)
			newC := c

			if h || (a&0x0F) > 0x09 {
				correction |= 0x06
			}

			if c || a > 0x99 {
				correction |= 0x60
				newC = true
			}

			if n {
				// Subtraction mode
				cpu.Regs.SetFlag(FlagH, h && (a&0x0F) < 6)
				a -= correction
			} else {
				// Addition mode
				cpu.Regs.SetFlag(FlagH, (a&0x0F) > 0x09)
				a += correction
			}

			cpu.Regs.A = a
			cpu.Regs.SetFlag(FlagC, newC)
			cpu.Regs.SetFlag(FlagZ, a == 0)
			cpu.Regs.SetFlag(FlagS, (a&0x80) != 0)
			
			// Parity calculation
			parity := true
			for i := 0; i < 8; i++ {
				if (a & (1 << i)) != 0 {
					parity = !parity
				}
			}
			cpu.Regs.SetFlag(FlagPV, parity)
			
			// Undocumented flags 3 and 5 are from the result
			cpu.Regs.SetFlag(Flag3, (a&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (a&0x20) != 0)

			return 4
		},
	})

	// 0xED 0x67: RRD (Rotate Right Decimal)
	RegisterInstruction(&EDTable, 0x67, Instruction{
		Mnemonic:  "RRD",
		Length:    2,
		Cycles:    18,
		AddrMode1: AddrModeIndirect,
		AddrMode2: AddrModeAccumulator,
		Execute: func(cpu *CPU) int {
			hlVal := cpu.Memory.Read(cpu.Regs.HL())
			aLow := cpu.Regs.A & 0x0F
			hlLow := hlVal & 0x0F
			hlHigh := hlVal >> 4

			// HL low -> A low
			// A low -> HL high
			// HL high -> HL low
			newHl := (aLow << 4) | hlHigh
			cpu.Memory.Write(cpu.Regs.HL(), newHl)
			cpu.Regs.A = (cpu.Regs.A & 0xF0) | hlLow

			// Flags
			res := cpu.Regs.A
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagN, false)
			cpu.Regs.SetFlag(FlagZ, res == 0)
			cpu.Regs.SetFlag(FlagS, (res&0x80) != 0)
			
			parity := true
			for i := 0; i < 8; i++ {
				if (res & (1 << i)) != 0 {
					parity = !parity
				}
			}
			cpu.Regs.SetFlag(FlagPV, parity)
			cpu.Regs.SetFlag(Flag3, (res&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (res&0x20) != 0)

			return 18
		},
	})

	// 0xED 0x6F: RLD (Rotate Left Decimal)
	RegisterInstruction(&EDTable, 0x6F, Instruction{
		Mnemonic:  "RLD",
		Length:    2,
		Cycles:    18,
		AddrMode1: AddrModeIndirect,
		AddrMode2: AddrModeAccumulator,
		Execute: func(cpu *CPU) int {
			hlVal := cpu.Memory.Read(cpu.Regs.HL())
			aLow := cpu.Regs.A & 0x0F
			hlLow := hlVal & 0x0F
			hlHigh := hlVal >> 4

			// HL low -> HL high
			// HL high -> A low
			// A low -> HL low
			newHl := (hlLow << 4) | aLow
			cpu.Memory.Write(cpu.Regs.HL(), newHl)
			cpu.Regs.A = (cpu.Regs.A & 0xF0) | hlHigh

			// Flags
			res := cpu.Regs.A
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagN, false)
			cpu.Regs.SetFlag(FlagZ, res == 0)
			cpu.Regs.SetFlag(FlagS, (res&0x80) != 0)
			
			parity := true
			for i := 0; i < 8; i++ {
				if (res & (1 << i)) != 0 {
					parity = !parity
				}
			}
			cpu.Regs.SetFlag(FlagPV, parity)
			cpu.Regs.SetFlag(Flag3, (res&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (res&0x20) != 0)

			return 18
		},
	})
}
