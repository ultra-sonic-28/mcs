package z80

// initADD registers Addition instructions.
func initADD() {
	// 0xC6: ADD A, n (Add Immediate to A)
	RegisterInstruction(&MainTable, 0xC6, Instruction{
		Mnemonic:  "ADD A, n",
		Length:    2,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			oldA := cpu.Regs.A
			res := oldA + n
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdd8(oldA, n, res)
			return 7
		},
	})

	// ADD A, r (Add Register to A)
	// Opcodes: 0x80 (B), 0x81 (C), 0x82 (D), 0x83 (E), 0x84 (H), 0x85 (L), 0x87 (A)
	regsAdd := []struct {
		name string
		get  func(r *Registers) uint8
		op   uint8
	}{
		{"B", func(r *Registers) uint8 { return r.B }, 0x80},
		{"C", func(r *Registers) uint8 { return r.C }, 0x81},
		{"D", func(r *Registers) uint8 { return r.D }, 0x82},
		{"E", func(r *Registers) uint8 { return r.E }, 0x83},
		{"H", func(r *Registers) uint8 { return r.H }, 0x84},
		{"L", func(r *Registers) uint8 { return r.L }, 0x85},
		{"A", func(r *Registers) uint8 { return r.A }, 0x87},
	}

	for _, r := range regsAdd {
		r := r // capture for closure
		RegisterInstruction(&MainTable, r.op, Instruction{
			Mnemonic:  "ADD A, " + r.name,
			Length:    1,
			Cycles:    4,
			AddrMode1: AddrModeAccumulator,
			AddrMode2: AddrModeRegister,
			Execute: func(cpu *CPU) int {
				val := r.get(cpu.Regs)
				oldA := cpu.Regs.A
				res := oldA + val
				cpu.Regs.A = res
				cpu.Regs.UpdateFlagsAdd8(oldA, val, res)
				return 4
			},
		})
	}

	// 0x86: ADD A, (HL)
	RegisterInstruction(&MainTable, 0x86, Instruction{
		Mnemonic:  "ADD A, (HL)",
		Length:    1,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndirect,
		Execute: func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.Regs.HL())
			oldA := cpu.Regs.A
			res := oldA + val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdd8(oldA, val, res)
			return 7
		},
	})

	// ADD HL, rr (Add Register Pair to HL)
	// Opcodes: 0x09 (BC), 0x19 (DE), 0x29 (HL), 0x39 (SP)
	rrAdd := []struct {
		name string
		get  func(r *Registers) uint16
		op   uint8
	}{
		{"BC", func(r *Registers) uint16 { return r.BC() }, 0x09},
		{"DE", func(r *Registers) uint16 { return r.DE() }, 0x19},
		{"HL", func(r *Registers) uint16 { return r.HL() }, 0x29},
		{"SP", func(r *Registers) uint16 { return r.SP }, 0x39},
	}

	for _, rr := range rrAdd {
		rr := rr // capture for closure
		RegisterInstruction(&MainTable, rr.op, Instruction{
			Mnemonic:  "ADD HL, " + rr.name,
			Length:    1,
			Cycles:    11,
			AddrMode1: AddrModeRegisterPair,
			AddrMode2: AddrModeRegisterPair,
			Execute: func(cpu *CPU) int {
				val := rr.get(cpu.Regs)
				oldHL := cpu.Regs.HL()
				res32 := uint32(oldHL) + uint32(val)
				res16 := uint16(res32)
				cpu.Regs.SetHL(res16)
				cpu.Regs.UpdateFlagsAdd16(oldHL, val, res32)
				return 11
			},
		})
	}

	// ADD IX, rr (Add Register Pair to IX)
	// Opcodes: 0xDD 0x09 (BC), 0xDD 0x19 (DE), 0xDD 0x29 (IX), 0xDD 0x39 (SP)
	rrAddIX := []struct {
		name string
		get  func(r *Registers) uint16
		op   uint8
	}{
		{"BC", func(r *Registers) uint16 { return r.BC() }, 0x09},
		{"DE", func(r *Registers) uint16 { return r.DE() }, 0x19},
		{"IX", func(r *Registers) uint16 { return r.IX }, 0x29},
		{"SP", func(r *Registers) uint16 { return r.SP }, 0x39},
	}

	for _, rr := range rrAddIX {
		rr := rr // capture for closure
		RegisterInstruction(&DDTable, rr.op, Instruction{
			Mnemonic:  "ADD IX, " + rr.name,
			Length:    2,
			Cycles:    15,
			AddrMode1: AddrModeRegisterPair,
			AddrMode2: AddrModeRegisterPair,
			Execute: func(cpu *CPU) int {
				val := rr.get(cpu.Regs)
				oldIX := cpu.Regs.IX
				res32 := uint32(oldIX) + uint32(val)
				res16 := uint16(res32)
				cpu.Regs.IX = res16
				cpu.Regs.UpdateFlagsAdd16(oldIX, val, res32)
				return 15
			},
		})
	}

	// ADD IY, rr (Add Register Pair to IY)
	// Opcodes: 0xFD 0x09 (BC), 0xFD 0x19 (DE), 0xFD 0x29 (IY), 0xFD 0x39 (SP)
	rrAddIY := []struct {
		name string
		get  func(r *Registers) uint16
		op   uint8
	}{
		{"BC", func(r *Registers) uint16 { return r.BC() }, 0x09},
		{"DE", func(r *Registers) uint16 { return r.DE() }, 0x19},
		{"IY", func(r *Registers) uint16 { return r.IY }, 0x29},
		{"SP", func(r *Registers) uint16 { return r.SP }, 0x39},
	}

	for _, rr := range rrAddIY {
		rr := rr // capture for closure
		RegisterInstruction(&FDTable, rr.op, Instruction{
			Mnemonic:  "ADD IY, " + rr.name,
			Length:    2,
			Cycles:    15,
			AddrMode1: AddrModeRegisterPair,
			AddrMode2: AddrModeRegisterPair,
			Execute: func(cpu *CPU) int {
				val := rr.get(cpu.Regs)
				oldIY := cpu.Regs.IY
				res32 := uint32(oldIY) + uint32(val)
				res16 := uint16(res32)
				cpu.Regs.IY = res16
				cpu.Regs.UpdateFlagsAdd16(oldIY, val, res32)
				return 15
			},
		})
	}

	// ADD A, IXH / ADD A, IXL / ADD A, (IX+d)
	RegisterInstruction(&DDTable, 0x84, Instruction{
		Mnemonic:  "ADD A, IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXH()
			oldA := cpu.Regs.A
			res := oldA + val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdd8(oldA, val, res)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x85, Instruction{
		Mnemonic:  "ADD A, IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IXL()
			oldA := cpu.Regs.A
			res := oldA + val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdd8(oldA, val, res)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x86, Instruction{
		Mnemonic:  "ADD A, (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			val := cpu.Memory.Read(addr)
			oldA := cpu.Regs.A
			res := oldA + val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdd8(oldA, val, res)
			return 19
		},
	})

	// ADD A, IYH / ADD A, IYL / ADD A, (IY+d)
	RegisterInstruction(&FDTable, 0x84, Instruction{
		Mnemonic:  "ADD A, IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYH()
			oldA := cpu.Regs.A
			res := oldA + val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdd8(oldA, val, res)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x85, Instruction{
		Mnemonic:  "ADD A, IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IYL()
			oldA := cpu.Regs.A
			res := oldA + val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdd8(oldA, val, res)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x86, Instruction{
		Mnemonic:  "ADD A, (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			val := cpu.Memory.Read(addr)
			oldA := cpu.Regs.A
			res := oldA + val
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdd8(oldA, val, res)
			return 19
		},
	})

	// ADC A, IXH / ADC A, IXL
	RegisterInstruction(&DDTable, 0x8C, Instruction{
		Mnemonic:  "ADC A, IXH",
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
			res := oldA + val + carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdc8(oldA, val, res, carry)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x8D, Instruction{
		Mnemonic:  "ADC A, IXL",
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
			res := oldA + val + carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdc8(oldA, val, res, carry)
			return 8
		},
	})

	// ADC A, IYH / ADC A, IYL
	RegisterInstruction(&FDTable, 0x8C, Instruction{
		Mnemonic:  "ADC A, IYH",
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
			res := oldA + val + carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdc8(oldA, val, res, carry)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x8D, Instruction{
		Mnemonic:  "ADC A, IYL",
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
			res := oldA + val + carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdc8(oldA, val, res, carry)
			return 8
		},
	})

	// 8. ADC A, r / ADC A, n / ADC A, (HL) / ADC A, (IX+d) / ADC A, (IY+d)
	// ADC A, r
	// Opcodes: 0x88 (B), 0x89 (C), 0x8A (D), 0x8B (E), 0x8C (H), 0x8D (L), 0x8F (A)
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

	for _, r := range regsAdc {
		r := r // capture for closure
		RegisterInstruction(&MainTable, r.op, Instruction{
			Mnemonic:  "ADC A, " + r.name,
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
				res := oldA + val + carry
				cpu.Regs.A = res
				cpu.Regs.UpdateFlagsAdc8(oldA, val, res, carry)
				return 4
			},
		})
	}

	// 0xCE: ADC A, n
	RegisterInstruction(&MainTable, 0xCE, Instruction{
		Mnemonic:  "ADC A, n",
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
			res := oldA + n + carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdc8(oldA, n, res, carry)
			return 7
		},
	})

	// 0x8E: ADC A, (HL)
	RegisterInstruction(&MainTable, 0x8E, Instruction{
		Mnemonic:  "ADC A, (HL)",
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
			res := oldA + val + carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdc8(oldA, val, res, carry)
			return 7
		},
	})

	// 0xDD 0x8E: ADC A, (IX+d)
	RegisterInstruction(&DDTable, 0x8E, Instruction{
		Mnemonic:  "ADC A, (IX+d)",
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
			res := oldA + val + carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdc8(oldA, val, res, carry)
			return 19
		},
	})

	// 0xFD 0x8E: ADC A, (IY+d)
	RegisterInstruction(&FDTable, 0x8E, Instruction{
		Mnemonic:  "ADC A, (IY+d)",
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
			res := oldA + val + carry
			cpu.Regs.A = res
			cpu.Regs.UpdateFlagsAdc8(oldA, val, res, carry)
			return 19
		},
	})

	// 9. ADC HL, rr
	// ADC HL, rr
	// Opcodes: 0xED 0x4A (BC), 0xED 0x5A (DE), 0xED 0x6A (HL), 0xED 0x7A (SP)
	rrAdc := []struct {
		name string
		get  func(r *Registers) uint16
		op   uint8
	}{
		{"BC", func(r *Registers) uint16 { return r.BC() }, 0x4A},
		{"DE", func(r *Registers) uint16 { return r.DE() }, 0x5A},
		{"HL", func(r *Registers) uint16 { return r.HL() }, 0x6A},
		{"SP", func(r *Registers) uint16 { return r.SP }, 0x7A},
	}

	for _, rr := range rrAdc {
		rr := rr // capture for closure
		RegisterInstruction(&EDTable, rr.op, Instruction{
			Mnemonic:  "ADC HL, " + rr.name,
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
				res32 := uint32(oldHL) + uint32(val) + carry
				res16 := uint16(res32)
				cpu.Regs.SetHL(res16)
				cpu.Regs.UpdateFlagsAdc16(uint32(oldHL), uint32(val), res32, carry)
				return 15
			},
		})
	}
}
