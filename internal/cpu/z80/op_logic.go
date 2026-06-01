package z80

// initLogic registers Logical instructions.
func initLogic() {
	regsAdc := []struct {
		name string
		get  func(r *Registers) uint8
		op   uint8
	}{
		{"B", func(r *Registers) uint8 { return r.B }, 0x88},
		{"C", func(r *Registers) uint8 { return r.C }, 0x89},
		{"D", func(r *Registers) uint8 { return r.D }, 0x8A},
		{"E", func(r *Registers) uint8 { return r.E }, 0x8B},
		{"H", func(r *Registers) uint8 { return r.H }, 0x8C},
		{"L", func(r *Registers) uint8 { return r.L }, 0x8D},
		{"A", func(r *Registers) uint8 { return r.A }, 0x8F},
	}

	// 10. AND s / OR s / XOR s
	// s can be r, n, (HL)

	// --- AND ---
	for _, r := range regsAdc {
		r := r // capture
		RegisterInstruction(&MainTable, 0xA0+r.op-0x88, Instruction{
			Mnemonic:  "AND " + r.name,
			Length:    1,
			Cycles:    4,
			AddrMode1: AddrModeAccumulator,
			AddrMode2: AddrModeRegister,
			Execute: func(cpu *CPU) int {
				val := r.get(cpu.Regs)
				res := cpu.Regs.A & val
				cpu.Regs.A = res
				cpu.Regs.UpdateFlagsLogical8(res, true)
				return 4
			},
		})
	}
	RegisterInstruction(&MainTable, 0xE6, Instruction{
		Mnemonic:  "AND n",
		Length:    2,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			res := cpu.Regs.A & n
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, true)
			return 7
		},
	})
	RegisterInstruction(&MainTable, 0xA6, Instruction{
		Mnemonic:  "AND (HL)",
		Length:    1,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndirect,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			res := cpu.Regs.A & val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, true)
			return 7
		},
	})

	// --- IX Logical Instructions (prefixed with 0xDD) ---

	// 11. AND IXH / AND IXL / AND (IX+d)
	RegisterInstruction(&DDTable, 0xA4, Instruction{
		Mnemonic:  "AND IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXH()
			res := cpu.Regs.A & val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, true)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0xA5, Instruction{
		Mnemonic:  "AND IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXL()
			res := cpu.Regs.A & val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, true)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0xA6, Instruction{
		Mnemonic:  "AND (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			val := cpu.Memory.Read(addr)
			res := cpu.Regs.A & val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, true)
			return 19
		},
	})

	// 12. XOR IXH / XOR IXL / XOR (IX+d)
	RegisterInstruction(&DDTable, 0xAC, Instruction{
		Mnemonic:  "XOR IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXH()
			res := cpu.Regs.A ^ val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0xAD, Instruction{
		Mnemonic:  "XOR IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXL()
			res := cpu.Regs.A ^ val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0xAE, Instruction{
		Mnemonic:  "XOR (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			val := cpu.Memory.Read(addr)
			res := cpu.Regs.A ^ val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 19
		},
	})

	// 13. OR IXH / OR IXL / OR (IX+d)
	RegisterInstruction(&DDTable, 0xB4, Instruction{
		Mnemonic:  "OR IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXH()
			res := cpu.Regs.A | val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0xB5, Instruction{
		Mnemonic:  "OR IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXL()
			res := cpu.Regs.A | val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0xB6, Instruction{
		Mnemonic:  "OR (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			val := cpu.Memory.Read(addr)
			res := cpu.Regs.A | val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 19
		},
	})

	// --- IY Logical Instructions (prefixed with 0xFD) ---

	// 14. AND IYH / AND IYL / AND (IY+d)
	RegisterInstruction(&FDTable, 0xA4, Instruction{
		Mnemonic:  "AND IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYH()
			res := cpu.Regs.A & val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, true)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0xA5, Instruction{
		Mnemonic:  "AND IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYL()
			res := cpu.Regs.A & val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, true)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0xA6, Instruction{
		Mnemonic:  "AND (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			val := cpu.Memory.Read(addr)
			res := cpu.Regs.A & val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, true)
			return 19
		},
	})

	// 15. XOR IYH / XOR IYL / XOR (IY+d)
	RegisterInstruction(&FDTable, 0xAC, Instruction{
		Mnemonic:  "XOR IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYH()
			res := cpu.Regs.A ^ val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0xAD, Instruction{
		Mnemonic:  "XOR IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYL()
			res := cpu.Regs.A ^ val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0xAE, Instruction{
		Mnemonic:  "XOR (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			val := cpu.Memory.Read(addr)
			res := cpu.Regs.A ^ val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 19
		},
	})

	// 16. OR IYH / OR IYL / OR (IY+d)
	RegisterInstruction(&FDTable, 0xB4, Instruction{
		Mnemonic:  "OR IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYH()
			res := cpu.Regs.A | val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0xB5, Instruction{
		Mnemonic:  "OR IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYL()
			res := cpu.Regs.A | val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0xB6, Instruction{
		Mnemonic:  "OR (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			val := cpu.Memory.Read(addr)
			res := cpu.Regs.A | val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 19
		},
	})

	// --- OR ---
	for _, r := range regsAdc {
		r := r // capture
		RegisterInstruction(&MainTable, 0xB0+r.op-0x88, Instruction{
			Mnemonic:  "OR " + r.name,
			Length:    1,
			Cycles:    4,
			AddrMode1: AddrModeAccumulator,
			AddrMode2: AddrModeRegister,
			Execute: func(cpu *CPU) int {
				val := r.get(cpu.Regs)
				res := cpu.Regs.A | val
				cpu.Regs.A = res
				cpu.Regs.UpdateFlagsLogical8(res, false)
				return 4
			},
		})
	}
	RegisterInstruction(&MainTable, 0xF6, Instruction{
		Mnemonic:  "OR n",
		Length:    2,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			res := cpu.Regs.A | n
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 7
		},
	})
	RegisterInstruction(&MainTable, 0xB6, Instruction{
		Mnemonic:  "OR (HL)",
		Length:    1,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndirect,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			res := cpu.Regs.A | val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 7
		},
	})

	// --- XOR ---
	for _, r := range regsAdc {
		r := r // capture
		RegisterInstruction(&MainTable, 0xA8+r.op-0x88, Instruction{
			Mnemonic:  "XOR " + r.name,
			Length:    1,
			Cycles:    4,
			AddrMode1: AddrModeAccumulator,
			AddrMode2: AddrModeRegister,
			Execute: func(cpu *CPU) int {
				val := r.get(cpu.Regs)
				res := cpu.Regs.A ^ val
				cpu.Regs.A = res
				cpu.Regs.UpdateFlagsLogical8(res, false)
				return 4
			},
		})
	}
	RegisterInstruction(&MainTable, 0xEE, Instruction{
		Mnemonic:  "XOR n",
		Length:    2,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			res := cpu.Regs.A ^ n
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 7
		},
	})
	RegisterInstruction(&MainTable, 0xAE, Instruction{
		Mnemonic:  "XOR (HL)",
		Length:    1,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndirect,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			res := cpu.Regs.A ^ val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsLogical8(res, false)
			return 7
		},
	})
}
