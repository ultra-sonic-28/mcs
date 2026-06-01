package z80

// initSUB registers Subtraction instructions.
func initSUB() {
	// 8. ADC A, r / ADC A, n / ADC A, (HL) / ADC A, (IX+d) / ADC A, (IY+d)
	// (Actually SUB/SBC/AND/OR/XOR/CP were added to initADD, let's move them here)

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

	// 10. SUB s / SBC A, s / CP s
	// s can be r, n, (HL), (IX+d), (IY+d)

	// --- SUB ---
	for _, r := range regsAdc {
		r := r // capture
		RegisterInstruction(&MainTable, 0x90+r.op-0x88, Instruction{
			Mnemonic:  "SUB " + r.name,
			Length:    1,
			Cycles:    4,
			AddrMode1: AddrModeAccumulator,
			AddrMode2: AddrModeRegister,
			Execute: func(cpu *CPU) int {
				val := r.get(cpu.Regs)
				oldA := cpu.Regs.A
				res := oldA - val
				cpu.Regs.A = res
				cpu.Regs.UpdateFlagsSub8(oldA, val, res)
				return 4
			},
		})
	}
	RegisterInstruction(&MainTable, 0xD6, Instruction{
		Mnemonic:  "SUB n",
		Length:    2,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			oldA := cpu.Regs.A
			res := oldA - n
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSub8(oldA, n, res)
			return 7
		},
	})
	RegisterInstruction(&MainTable, 0x96, Instruction{
		Mnemonic:  "SUB (HL)",
		Length:    1,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndirect,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSub8(oldA, val, res)
			return 7
		},
	})

	// --- SBC A, s ---
	for _, r := range regsAdc {
		r := r // capture
		RegisterInstruction(&MainTable, 0x98+r.op-0x88, Instruction{
			Mnemonic:  "SBC A, " + r.name,
			Length:    1,
			Cycles:    4,
			AddrMode1: AddrModeAccumulator,
			AddrMode2: AddrModeRegister,
			Execute: func(cpu *CPU) int {
				val := r.get(cpu.Regs)
				oldA := cpu.Regs.A
				carry := uint8(0)
				if cpu.Regs.Flag(FlagC) {
					carry = 1
				}
				res := oldA - val - carry
				cpu.Regs.A = res
				cpu.Regs.UpdateFlagsSbc8(oldA, val, res, carry)
				return 4
			},
		})
	}
	RegisterInstruction(&MainTable, 0xDE, Instruction{
		Mnemonic:  "SBC A, n",
		Length:    2,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			oldA := cpu.Regs.A
			carry := uint8(0)
			if cpu.Regs.Flag(FlagC) {
				carry = 1
			}
			res := oldA - n - carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSbc8(oldA, n, res, carry)
			return 7
		},
	})
	RegisterInstruction(&MainTable, 0x9E, Instruction{
		Mnemonic:  "SBC A, (HL)",
		Length:    1,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndirect,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			oldA := cpu.Regs.A
			carry := uint8(0)
			if cpu.Regs.Flag(FlagC) {
				carry = 1
			}
			res := oldA - val - carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSbc8(oldA, val, res, carry)
			return 7
		},
	})

	// --- CP s ---
	for _, r := range regsAdc {
		r := r // capture
		RegisterInstruction(&MainTable, 0xB8+r.op-0x88, Instruction{
			Mnemonic:  "CP " + r.name,
			Length:    1,
			Cycles:    4,
			AddrMode1: AddrModeAccumulator,
			AddrMode2: AddrModeRegister,
			Execute: func(cpu *CPU) int {
				val := r.get(cpu.Regs)
				oldA := cpu.Regs.A
				res := oldA - val
				cpu.Regs.UpdateFlagsCp8(oldA, val, res)
				return 4
			},
		})
	}
	RegisterInstruction(&MainTable, 0xFE, Instruction{
		Mnemonic:  "CP n",
		Length:    2,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			oldA := cpu.Regs.A
			res := oldA - n
			cpu.Regs.UpdateFlagsCp8(oldA, n, res)
			return 7
		},
	})
	RegisterInstruction(&MainTable, 0xBE, Instruction{
		Mnemonic:  "CP (HL)",
		Length:    1,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndirect,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.UpdateFlagsCp8(oldA, val, res)
			return 7
		},
	})

	// SUB A, IXH / SUB A, IXL / SUB A, (IX+d)
	RegisterInstruction(&DDTable, 0x94, Instruction{
		Mnemonic:  "SUB A, IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXH()
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSub8(oldA, val, res)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x95, Instruction{
		Mnemonic:  "SUB A, IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXL()
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSub8(oldA, val, res)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x96, Instruction{
		Mnemonic:  "SUB A, (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			val := cpu.Memory.Read(addr)
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSub8(oldA, val, res)
			return 19
		},
	})

	// SUB A, IYH / SUB A, IYL / SUB A, (IY+d)
	RegisterInstruction(&FDTable, 0x94, Instruction{
		Mnemonic:  "SUB A, IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYH()
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSub8(oldA, val, res)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x95, Instruction{
		Mnemonic:  "SUB A, IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYL()
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSub8(oldA, val, res)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x96, Instruction{
		Mnemonic:  "SUB A, (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			val := cpu.Memory.Read(addr)
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSub8(oldA, val, res)
			return 19
		},
	})

	// SBC A, IXH / SBC A, IXL / SBC A, (IX+d)
	RegisterInstruction(&DDTable, 0x9C, Instruction{
		Mnemonic:  "SBC A, IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXH()
			oldA := cpu.Regs.A
			carry := uint8(0)
			if cpu.Regs.Flag(FlagC) {
				carry = 1
			}
			res := oldA - val - carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSbc8(oldA, val, res, carry)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x9D, Instruction{
		Mnemonic:  "SBC A, IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXL()
			oldA := cpu.Regs.A
			carry := uint8(0)
			if cpu.Regs.Flag(FlagC) {
				carry = 1
			}
			res := oldA - val - carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSbc8(oldA, val, res, carry)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x9E, Instruction{
		Mnemonic:  "SBC A, (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			val := cpu.Memory.Read(addr)
			oldA := cpu.Regs.A
			carry := uint8(0)
			if cpu.Regs.Flag(FlagC) {
				carry = 1
			}
			res := oldA - val - carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSbc8(oldA, val, res, carry)
			return 19
		},
	})

	// SBC A, IYH / SBC A, IYL / SBC A, (IY+d)
	RegisterInstruction(&FDTable, 0x9C, Instruction{
		Mnemonic:  "SBC A, IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYH()
			oldA := cpu.Regs.A
			carry := uint8(0)
			if cpu.Regs.Flag(FlagC) {
				carry = 1
			}
			res := oldA - val - carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSbc8(oldA, val, res, carry)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x9D, Instruction{
		Mnemonic:  "SBC A, IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYL()
			oldA := cpu.Regs.A
			carry := uint8(0)
			if cpu.Regs.Flag(FlagC) {
				carry = 1
			}
			res := oldA - val - carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSbc8(oldA, val, res, carry)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x9E, Instruction{
		Mnemonic:  "SBC A, (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			val := cpu.Memory.Read(addr)
			oldA := cpu.Regs.A
			carry := uint8(0)
			if cpu.Regs.Flag(FlagC) {
				carry = 1
			}
			res := oldA - val - carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsSbc8(oldA, val, res, carry)
			return 19
		},
	})

	// CP IXH / CP IXL / CP (IX+d)
	RegisterInstruction(&DDTable, 0xBC, Instruction{
		Mnemonic:  "CP IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXH()
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.UpdateFlagsCp8(oldA, val, res)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0xBD, Instruction{
		Mnemonic:  "CP IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXL()
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.UpdateFlagsCp8(oldA, val, res)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0xBE, Instruction{
		Mnemonic:  "CP (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			val := cpu.Memory.Read(addr)
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.UpdateFlagsCp8(oldA, val, res)
			return 19
		},
	})

	// CP IYH / CP IYL / CP (IY+d)
	RegisterInstruction(&FDTable, 0xBC, Instruction{
		Mnemonic:  "CP IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYH()
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.UpdateFlagsCp8(oldA, val, res)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0xBD, Instruction{
		Mnemonic:  "CP IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYL()
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.UpdateFlagsCp8(oldA, val, res)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0xBE, Instruction{
		Mnemonic:  "CP (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			val := cpu.Memory.Read(addr)
			oldA := cpu.Regs.A
			res := oldA - val
			cpu.Regs.UpdateFlagsCp8(oldA, val, res)
			return 19
		},
	})

	// --- ED Prefix Subtraction Instructions ---

	// SBC HL, rr
	sbcHLRegs := []struct {
		name string
		get  func(r *Registers) uint16
		op   uint8
	}{
		{"BC", func(r *Registers) uint16 { return r.BC() }, 0x42},
		{"DE", func(r *Registers) uint16 { return r.DE() }, 0x52},
		{"HL", func(r *Registers) uint16 { return r.HL() }, 0x62},
		{"SP", func(r *Registers) uint16 { return r.SP }, 0x72},
	}

	for _, rr := range sbcHLRegs {
		rr := rr // capture
		RegisterInstruction(&EDTable, rr.op, Instruction{
			Mnemonic:  "SBC HL, " + rr.name,
			Length:    2,
			Cycles:    15,
			AddrMode1: AddrModeRegisterPair,
			AddrMode2: AddrModeRegisterPair,
			Execute: func(cpu *CPU) int {
				val := rr.get(cpu.Regs)
				oldHL := cpu.Regs.HL()
				carry := uint32(0)
				if cpu.Regs.Flag(FlagC) {
					carry = 1
				}
				res := uint32(oldHL) - uint32(val) - carry
				cpu.Regs.SetHL(uint16(res))
				cpu.Regs.UpdateFlagsSbc16(oldHL, val, res, carry)
				return 15
			},
		})
	}
}
