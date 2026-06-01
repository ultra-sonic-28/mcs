package z80

// initLD registers Load instructions.
func initLD() {
	regs8 := []struct {
		name string
		get  func(r *Registers) uint8
		set  func(r *Registers, val uint8)
	}{
		{"B", func(r *Registers) uint8 { return r.B }, func(r *Registers, val uint8) { r.B = val }},
		{"C", func(r *Registers) uint8 { return r.C }, func(r *Registers, val uint8) { r.C = val }},
		{"D", func(r *Registers) uint8 { return r.D }, func(r *Registers, val uint8) { r.D = val }},
		{"E", func(r *Registers) uint8 { return r.E }, func(r *Registers, val uint8) { r.E = val }},
		{"H", func(r *Registers) uint8 { return r.H }, func(r *Registers, val uint8) { r.H = val }},
		{"L", func(r *Registers) uint8 { return r.L }, func(r *Registers, val uint8) { r.L = val }},
		{"(HL)", nil, nil}, // Placeholder for (HL)
		{"A", func(r *Registers) uint8 { return r.A }, func(r *Registers, val uint8) { r.A = val }},
	}

	// 1. LD r, r' (Register to Register)
	// 2. LD r, (HL) (Memory to Register)
	// 3. LD (HL), r (Register to Memory)
	for dIdx, dest := range regs8 {
		for sIdx, src := range regs8 {
			if dIdx == 6 && sIdx == 6 {
				// 0x76 is HALT, not LD (HL), (HL)
				continue
			}

			opcode := uint8(0x40 + dIdx*8 + sIdx)
			d, s := dest, src // capture for closure

			mnemonic := "LD " + d.name + ", " + s.name
			cycles := 4
			addr1 := AddrModeRegister
			addr2 := AddrModeRegister

			if dIdx == 6 {
				addr1 = AddrModeIndirect
				cycles = 7
			}
			if sIdx == 6 {
				addr2 = AddrModeIndirect
				cycles = 7
			}

			RegisterInstruction(&MainTable, opcode, Instruction{
				Mnemonic:  mnemonic,
				Length:    1,
				Cycles:    cycles,
				AddrMode1: addr1,
				AddrMode2: addr2,
				Execute: func(cpu *CPU) int {
					var val uint8
					if s.get != nil {
						val = s.get(cpu.Regs)
					} else {
						val = cpu.Memory.Read(cpu.Regs.HL())
					}

					if d.set != nil {
						d.set(cpu.Regs, val)
					} else {
						cpu.Memory.Write(cpu.Regs.HL(), val)
					}
					return cycles
				},
			})
		}
	}

	// 4. LD r, n (Immediate to Register)
	// 5. LD (HL), n (Immediate to Memory)
	for i, reg := range regs8 {
		opcode := uint8(0x06 + i*8)
		r := reg // capture
		mnemonic := "LD " + r.name + ", n"
		cycles := 7
		addr1 := AddrModeRegister
		addr2 := AddrModeImmediate

		if i == 6 { // (HL)
			addr1 = AddrModeIndirect
			cycles = 10
		}

		RegisterInstruction(&MainTable, opcode, Instruction{
			Mnemonic:  mnemonic,
			Length:    2,
			Cycles:    cycles,
			AddrMode1: addr1,
			AddrMode2: addr2,
			Execute: func(cpu *CPU) int {
				n := cpu.FetchByte()
				if r.set != nil {
					r.set(cpu.Regs, n)
				} else {
					cpu.Memory.Write(cpu.Regs.HL(), n)
				}
				return cycles
			},
		})
	}

	// 6. LD dd, nn (16-bit immediate to BC, DE, HL, SP)
	ddRegs := []struct {
		name string
		set  func(r *Registers, val uint16)
	}{
		{"BC", func(r *Registers, val uint16) { r.SetBC(val) }},
		{"DE", func(r *Registers, val uint16) { r.SetDE(val) }},
		{"HL", func(r *Registers, val uint16) { r.SetHL(val) }},
		{"SP", func(r *Registers, val uint16) { r.SP = val }},
	}

	for i, dd := range ddRegs {
		opcode := uint8(0x01 + i*16)
		d := dd // capture
		RegisterInstruction(&MainTable, opcode, Instruction{
			Mnemonic:  "LD " + d.name + ", nn",
			Length:    3,
			Cycles:    10,
			AddrMode1: AddrModeRegisterPair,
			AddrMode2: AddrModeImmediate16,
			Execute: func(cpu *CPU) int {
				nn := cpu.FetchWord()
				d.set(cpu.Regs, nn)
				return 10
			},
		})
	}

	// 7. LD (BC), A / LD A, (BC) / LD (DE), A / LD A, (DE)
	RegisterInstruction(&MainTable, 0x02, Instruction{
		Mnemonic:  "LD (BC), A",
		Length:    1,
		Cycles:    7,
		AddrMode1: AddrModeIndirect,
		AddrMode2: AddrModeAccumulator,
		Execute: func(cpu *CPU) int {
			cpu.Memory.Write(cpu.Regs.BC(), cpu.Regs.A)
			return 7
		},
	})
	RegisterInstruction(&MainTable, 0x0A, Instruction{
		Mnemonic:  "LD A, (BC)",
		Length:    1,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndirect,
		Execute: func(cpu *CPU) int {
			cpu.Regs.A = cpu.Memory.Read(cpu.Regs.BC())
			return 7
		},
	})
	RegisterInstruction(&MainTable, 0x12, Instruction{
		Mnemonic:  "LD (DE), A",
		Length:    1,
		Cycles:    7,
		AddrMode1: AddrModeIndirect,
		AddrMode2: AddrModeAccumulator,
		Execute: func(cpu *CPU) int {
			cpu.Memory.Write(cpu.Regs.DE(), cpu.Regs.A)
			return 7
		},
	})
	RegisterInstruction(&MainTable, 0x1A, Instruction{
		Mnemonic:  "LD A, (DE)",
		Length:    1,
		Cycles:    7,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeIndirect,
		Execute: func(cpu *CPU) int {
			cpu.Regs.A = cpu.Memory.Read(cpu.Regs.DE())
			return 7
		},
	})

	// 8. LD (nn), HL / LD HL, (nn)
	RegisterInstruction(&MainTable, 0x22, Instruction{
		Mnemonic:  "LD (nn), HL",
		Length:    3,
		Cycles:    16,
		AddrMode1: AddrModeExtended,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			cpu.Memory.Write(nn, cpu.Regs.L)
			cpu.Memory.Write(nn+1, cpu.Regs.H)
			return 16
		},
	})
	RegisterInstruction(&MainTable, 0x2A, Instruction{
		Mnemonic:  "LD HL, (nn)",
		Length:    3,
		Cycles:    16,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeExtended,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			low := cpu.Memory.Read(nn)
			high := cpu.Memory.Read(nn + 1)
			cpu.Regs.SetHL((uint16(high) << 8) | uint16(low))
			return 16
		},
	})

	// 9. LD (nn), A / LD A, (nn)
	RegisterInstruction(&MainTable, 0x32, Instruction{
		Mnemonic:  "LD (nn), A",
		Length:    3,
		Cycles:    13,
		AddrMode1: AddrModeExtended,
		AddrMode2: AddrModeAccumulator,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			cpu.Memory.Write(nn, cpu.Regs.A)
			return 13
		},
	})
	RegisterInstruction(&MainTable, 0x3A, Instruction{
		Mnemonic:  "LD A, (nn)",
		Length:    3,
		Cycles:    13,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeExtended,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			cpu.Regs.A = cpu.Memory.Read(nn)
			return 13
		},
	})

	// 10. LD SP, HL
	RegisterInstruction(&MainTable, 0xF9, Instruction{
		Mnemonic:  "LD SP, HL",
		Length:    1,
		Cycles:    6,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SP = cpu.Regs.HL()
			return 6
		},
	})

	// --- IX Load Instructions (prefixed with 0xDD) ---

	// 11. LD IX, nn (16-bit immediate to IX)
	RegisterInstruction(&DDTable, 0x21, Instruction{
		Mnemonic:  "LD IX, nn",
		Length:    4,
		Cycles:    14,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeImmediate16,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			cpu.Regs.IX = nn
			return 14
		},
	})

	// 12. LD (nn), IX (IX to memory at nn)
	RegisterInstruction(&DDTable, 0x22, Instruction{
		Mnemonic:  "LD (nn), IX",
		Length:    4,
		Cycles:    20,
		AddrMode1: AddrModeExtended,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			cpu.Memory.Write(nn, uint8(cpu.Regs.IX&0xFF))
			cpu.Memory.Write(nn+1, uint8(cpu.Regs.IX>>8))
			return 20
		},
	})

	// 13. LD IX, (nn) (Memory at nn to IX)
	RegisterInstruction(&DDTable, 0x2A, Instruction{
		Mnemonic:  "LD IX, (nn)",
		Length:    4,
		Cycles:    20,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeExtended,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			low := cpu.Memory.Read(nn)
			high := cpu.Memory.Read(nn + 1)
			cpu.Regs.IX = (uint16(high) << 8) | uint16(low)
			return 20
		},
	})

	// 14. LD IXH, n (8-bit immediate to IX high byte)
	RegisterInstruction(&DDTable, 0x26, Instruction{
		Mnemonic:  "LD IXH, n",
		Length:    3,
		Cycles:    11,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			cpu.Regs.SetIXH(n)
			return 11
		},
	})

	// 15. LD IXL, n (8-bit immediate to IX low byte)
	RegisterInstruction(&DDTable, 0x2E, Instruction{
		Mnemonic:  "LD IXL, n",
		Length:    3,
		Cycles:    11,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			cpu.Regs.SetIXL(n)
			return 11
		},
	})

	// 16. LD (IX+d), n (8-bit immediate to indexed memory)
	RegisterInstruction(&DDTable, 0x36, Instruction{
		Mnemonic:  "LD (IX+d), n",
		Length:    4,
		Cycles:    19,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			n := cpu.FetchByte()
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			cpu.Memory.Write(addr, n)
			return 19
		},
	})

	// 17. LD B, IXH / LD B, IXL / LD B, (IX+d)
	RegisterInstruction(&DDTable, 0x44, Instruction{
		Mnemonic:  "LD B, IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.B = cpu.Regs.IXH()
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x45, Instruction{
		Mnemonic:  "LD B, IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.B = cpu.Regs.IXL()
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x46, Instruction{
		Mnemonic:  "LD B, (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			cpu.Regs.B = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 18. LD C, IXH / LD C, IXL / LD C, (IX+d)
	RegisterInstruction(&DDTable, 0x4C, Instruction{
		Mnemonic:  "LD C, IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.C = cpu.Regs.IXH()
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x4D, Instruction{
		Mnemonic:  "LD C, IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.C = cpu.Regs.IXL()
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x4E, Instruction{
		Mnemonic:  "LD C, (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			cpu.Regs.C = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 19. LD D, IXH / LD D, IXL / LD D, (IX+d)
	RegisterInstruction(&DDTable, 0x54, Instruction{
		Mnemonic:  "LD D, IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.D = cpu.Regs.IXH()
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x55, Instruction{
		Mnemonic:  "LD D, IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.D = cpu.Regs.IXL()
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x56, Instruction{
		Mnemonic:  "LD D, (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			cpu.Regs.D = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 20. LD E, IXH / LD E, IXL / LD E, (IX+d)
	RegisterInstruction(&DDTable, 0x5C, Instruction{
		Mnemonic:  "LD E, IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.E = cpu.Regs.IXH()
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x5D, Instruction{
		Mnemonic:  "LD E, IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.E = cpu.Regs.IXL()
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x5E, Instruction{
		Mnemonic:  "LD E, (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			cpu.Regs.E = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 21. LD IXH, r (r = B, C, D, E, IXH, IXL, A)
	ixhRegs := []struct {
		name string
		get  func(r *Registers) uint8
		op   uint8
	}{
		{"B", func(r *Registers) uint8 { return r.B }, 0x60},
		{"C", func(r *Registers) uint8 { return r.C }, 0x61},
		{"D", func(r *Registers) uint8 { return r.D }, 0x62},
		{"E", func(r *Registers) uint8 { return r.E }, 0x63},
		{"IXH", func(r *Registers) uint8 { return r.IXH() }, 0x64},
		{"IXL", func(r *Registers) uint8 { return r.IXL() }, 0x65},
		{"A", func(r *Registers) uint8 { return r.A }, 0x67},
	}

	for _, r := range ixhRegs {
		r := r // capture
		RegisterInstruction(&DDTable, r.op, Instruction{
			Mnemonic:  "LD IXH, " + r.name,
			Length:    2,
			Cycles:    8,
			AddrMode1: AddrModeRegister,
			AddrMode2: AddrModeRegister,
			Execute: func(cpu *CPU) int {
				cpu.Regs.SetIXH(r.get(cpu.Regs))
				return 8
			},
		})
	}

	// 22. LD IXL, r (r = B, C, D, E, IXH, IXL, A)
	ixlRegs := []struct {
		name string
		get  func(r *Registers) uint8
		op   uint8
	}{
		{"B", func(r *Registers) uint8 { return r.B }, 0x68},
		{"C", func(r *Registers) uint8 { return r.C }, 0x69},
		{"D", func(r *Registers) uint8 { return r.D }, 0x6A},
		{"E", func(r *Registers) uint8 { return r.E }, 0x6B},
		{"IXH", func(r *Registers) uint8 { return r.IXH() }, 0x6C},
		{"IXL", func(r *Registers) uint8 { return r.IXL() }, 0x6D},
		{"A", func(r *Registers) uint8 { return r.A }, 0x6F},
	}

	for _, r := range ixlRegs {
		r := r // capture
		RegisterInstruction(&DDTable, r.op, Instruction{
			Mnemonic:  "LD IXL, " + r.name,
			Length:    2,
			Cycles:    8,
			AddrMode1: AddrModeRegister,
			AddrMode2: AddrModeRegister,
			Execute: func(cpu *CPU) int {
				cpu.Regs.SetIXL(r.get(cpu.Regs))
				return 8
			},
		})
	}

	// 23. LD H, (IX+d) / LD L, (IX+d)
	RegisterInstruction(&DDTable, 0x66, Instruction{
		Mnemonic:  "LD H, (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			cpu.Regs.H = cpu.Memory.Read(addr)
			return 19
		},
	})
	RegisterInstruction(&DDTable, 0x6E, Instruction{
		Mnemonic:  "LD L, (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			cpu.Regs.L = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 24. LD (IX+d), r (r = B, C, D, E, H, L, A)
	ixdRegs := []struct {
		name string
		get  func(r *Registers) uint8
		op   uint8
	}{
		{"B", func(r *Registers) uint8 { return r.B }, 0x70},
		{"C", func(r *Registers) uint8 { return r.C }, 0x71},
		{"D", func(r *Registers) uint8 { return r.D }, 0x72},
		{"E", func(r *Registers) uint8 { return r.E }, 0x73},
		{"H", func(r *Registers) uint8 { return r.H }, 0x74},
		{"L", func(r *Registers) uint8 { return r.L }, 0x75},
		{"A", func(r *Registers) uint8 { return r.A }, 0x77},
	}

	for _, r := range ixdRegs {
		r := r // capture
		RegisterInstruction(&DDTable, r.op, Instruction{
			Mnemonic:  "LD (IX+d), " + r.name,
			Length:    3,
			Cycles:    19,
			AddrMode1: AddrModeIndexed,
			AddrMode2: AddrModeRegister,
			Execute: func(cpu *CPU) int {
				d := int8(cpu.FetchByte())
				addr := uint16(int32(cpu.Regs.IX) + int32(d))
				cpu.Memory.Write(addr, r.get(cpu.Regs))
				return 19
			},
		})
	}

	// 25. LD A, IXH / LD A, IXL / LD A, (IX+d)
	RegisterInstruction(&DDTable, 0x7C, Instruction{
		Mnemonic:  "LD A, IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.A = cpu.Regs.IXH()
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x7D, Instruction{
		Mnemonic:  "LD A, IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.A = cpu.Regs.IXL()
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x7E, Instruction{
		Mnemonic:  "LD A, (IX+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			cpu.Regs.A = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 26. LD SP, IX
	RegisterInstruction(&DDTable, 0xF9, Instruction{
		Mnemonic:  "LD SP, IX",
		Length:    2,
		Cycles:    10,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SP = cpu.Regs.IX
			return 10
		},
	})

	// 27. LD SP, IY
	RegisterInstruction(&FDTable, 0xF9, Instruction{
		Mnemonic:  "LD SP, IY",
		Length:    2,
		Cycles:    10,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SP = cpu.Regs.IY
			return 10
		},
	})

	// --- IY Load Instructions (prefixed with 0xFD) ---

	// 28. LD IY, nn (16-bit immediate to IY)
	RegisterInstruction(&FDTable, 0x21, Instruction{
		Mnemonic:  "LD IY, nn",
		Length:    4,
		Cycles:    14,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeImmediate16,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			cpu.Regs.IY = nn
			return 14
		},
	})

	// 28a. LD IYH, n (8-bit immediate to IY high byte)
	RegisterInstruction(&FDTable, 0x26, Instruction{
		Mnemonic:  "LD IYH, n",
		Length:    3,
		Cycles:    11,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			cpu.Regs.SetIYH(n)
			return 11
		},
	})

	// 28b. LD IYL, n (8-bit immediate to IY low byte)
	RegisterInstruction(&FDTable, 0x2E, Instruction{
		Mnemonic:  "LD IYL, n",
		Length:    3,
		Cycles:    11,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			cpu.Regs.SetIYL(n)
			return 11
		},
	})

	// 29. LD (nn), IY (IY to memory at nn)
	RegisterInstruction(&FDTable, 0x22, Instruction{
		Mnemonic:  "LD (nn), IY",
		Length:    4,
		Cycles:    20,
		AddrMode1: AddrModeExtended,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			cpu.Memory.Write(nn, uint8(cpu.Regs.IY&0xFF))
			cpu.Memory.Write(nn+1, uint8(cpu.Regs.IY>>8))
			return 20
		},
	})

	// 30. LD IY, (nn) (Memory at nn to IY)
	RegisterInstruction(&FDTable, 0x2A, Instruction{
		Mnemonic:  "LD IY, (nn)",
		Length:    4,
		Cycles:    20,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeExtended,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			low := cpu.Memory.Read(nn)
			high := cpu.Memory.Read(nn + 1)
			cpu.Regs.IY = (uint16(high) << 8) | uint16(low)
			return 20
		},
	})

	// 31. LD (IY+d), n (8-bit immediate to indexed memory)
	RegisterInstruction(&FDTable, 0x36, Instruction{
		Mnemonic:  "LD (IY+d), n",
		Length:    4,
		Cycles:    19,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeImmediate,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			n := cpu.FetchByte()
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Memory.Write(addr, n)
			return 19
		},
	})

	// 32. LD B, IYH / LD B, IYL / LD B, (IY+d)
	RegisterInstruction(&FDTable, 0x44, Instruction{
		Mnemonic:  "LD B, IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.B = cpu.Regs.IYH()
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x45, Instruction{
		Mnemonic:  "LD B, IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.B = cpu.Regs.IYL()
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x46, Instruction{
		Mnemonic:  "LD B, (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Regs.B = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 33. LD IYH, B
	RegisterInstruction(&FDTable, 0x60, Instruction{
		Mnemonic:  "LD IYH, B",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYH(cpu.Regs.B)
			return 8
		},
	})

	// 34. LD IYL, B
	RegisterInstruction(&FDTable, 0x68, Instruction{
		Mnemonic:  "LD IYL, B",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYL(cpu.Regs.B)
			return 8
		},
	})

	// 35. LD (IY+d), B
	RegisterInstruction(&FDTable, 0x70, Instruction{
		Mnemonic:  "LD (IY+d), B",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Memory.Write(addr, cpu.Regs.B)
			return 19
		},
	})

	// 36. LD IYH, A
	RegisterInstruction(&FDTable, 0x67, Instruction{
		Mnemonic:  "LD IYH, A",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYH(cpu.Regs.A)
			return 8
		},
	})

	// 37. LD IYL, A
	RegisterInstruction(&FDTable, 0x6F, Instruction{
		Mnemonic:  "LD IYL, A",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYL(cpu.Regs.A)
			return 8
		},
	})

	// 38. LD (IY+d), A
	RegisterInstruction(&FDTable, 0x77, Instruction{
		Mnemonic:  "LD (IY+d), A",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Memory.Write(addr, cpu.Regs.A)
			return 19
		},
	})

	// 39. LD A, IYH / LD A, IYL / LD A, (IY+d)
	RegisterInstruction(&FDTable, 0x7C, Instruction{
		Mnemonic:  "LD A, IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.A = cpu.Regs.IYH()
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x7D, Instruction{
		Mnemonic:  "LD A, IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.A = cpu.Regs.IYL()
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x7E, Instruction{
		Mnemonic:  "LD A, (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Regs.A = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 40. LD C, IYH / LD C, IYL / LD C, (IY+d)
	RegisterInstruction(&FDTable, 0x4C, Instruction{
		Mnemonic:  "LD C, IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.C = cpu.Regs.IYH()
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x4D, Instruction{
		Mnemonic:  "LD C, IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.C = cpu.Regs.IYL()
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x4E, Instruction{
		Mnemonic:  "LD C, (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Regs.C = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 41. LD IYH, C
	RegisterInstruction(&FDTable, 0x61, Instruction{
		Mnemonic:  "LD IYH, C",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYH(cpu.Regs.C)
			return 8
		},
	})

	// 42. LD IYL, C
	RegisterInstruction(&FDTable, 0x69, Instruction{
		Mnemonic:  "LD IYL, C",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYL(cpu.Regs.C)
			return 8
		},
	})

	// 43. LD (IY+d), C
	RegisterInstruction(&FDTable, 0x71, Instruction{
		Mnemonic:  "LD (IY+d), C",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Memory.Write(addr, cpu.Regs.C)
			return 19
		},
	})

	// 44. LD D, IYH / LD D, IYL / LD D, (IY+d)
	RegisterInstruction(&FDTable, 0x54, Instruction{
		Mnemonic:  "LD D, IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.D = cpu.Regs.IYH()
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x55, Instruction{
		Mnemonic:  "LD D, IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.D = cpu.Regs.IYL()
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x56, Instruction{
		Mnemonic:  "LD D, (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Regs.D = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 45. LD IYH, D
	RegisterInstruction(&FDTable, 0x62, Instruction{
		Mnemonic:  "LD IYH, D",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYH(cpu.Regs.D)
			return 8
		},
	})

	// 46. LD IYL, D
	RegisterInstruction(&FDTable, 0x6A, Instruction{
		Mnemonic:  "LD IYL, D",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYL(cpu.Regs.D)
			return 8
		},
	})

	// 47. LD (IY+d), D
	RegisterInstruction(&FDTable, 0x72, Instruction{
		Mnemonic:  "LD (IY+d), D",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Memory.Write(addr, cpu.Regs.D)
			return 19
		},
	})

	// 48. LD E, IYH / LD E, IYL / LD E, (IY+d)
	RegisterInstruction(&FDTable, 0x5C, Instruction{
		Mnemonic:  "LD E, IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.E = cpu.Regs.IYH()
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x5D, Instruction{
		Mnemonic:  "LD E, IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.E = cpu.Regs.IYL()
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x5E, Instruction{
		Mnemonic:  "LD E, (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Regs.E = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 49. LD IYH, E
	RegisterInstruction(&FDTable, 0x63, Instruction{
		Mnemonic:  "LD IYH, E",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYH(cpu.Regs.E)
			return 8
		},
	})

	// 50. LD IYL, E
	RegisterInstruction(&FDTable, 0x6B, Instruction{
		Mnemonic:  "LD IYL, E",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYL(cpu.Regs.E)
			return 8
		},
	})

	// 51. LD (IY+d), E
	RegisterInstruction(&FDTable, 0x73, Instruction{
		Mnemonic:  "LD (IY+d), E",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Memory.Write(addr, cpu.Regs.E)
			return 19
		},
	})

	// 52. LD IYH, IYH / LD IYH, IYL / LD H, (IY+d)
	RegisterInstruction(&FDTable, 0x64, Instruction{
		Mnemonic:  "LD IYH, IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYH(cpu.Regs.IYH())
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x65, Instruction{
		Mnemonic:  "LD IYH, IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYH(cpu.Regs.IYL())
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x66, Instruction{
		Mnemonic:  "LD H, (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Regs.H = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 53. LD IYL, IYH / LD IYL, IYL / LD L, (IY+d)
	RegisterInstruction(&FDTable, 0x6C, Instruction{
		Mnemonic:  "LD IYL, IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYL(cpu.Regs.IYH())
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x6D, Instruction{
		Mnemonic:  "LD IYL, IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.SetIYL(cpu.Regs.IYL())
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x6E, Instruction{
		Mnemonic:  "LD L, (IY+d)",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeIndexed,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Regs.L = cpu.Memory.Read(addr)
			return 19
		},
	})

	// 54. LD (IY+d), H
	RegisterInstruction(&FDTable, 0x74, Instruction{
		Mnemonic:  "LD (IY+d), H",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Memory.Write(addr, cpu.Regs.H)
			return 19
		},
	})

	// 55. LD (IY+d), L
	RegisterInstruction(&FDTable, 0x75, Instruction{
		Mnemonic:  "LD (IY+d), L",
		Length:    3,
		Cycles:    19,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			cpu.Memory.Write(addr, cpu.Regs.L)
			return 19
		},
	})

	// --- ED Prefix Load Instructions ---

	// 56. LD (nn), BC / LD (nn), DE / LD (nn), HL / LD (nn), SP
	RegisterInstruction(&EDTable, 0x43, Instruction{
		Mnemonic:  "LD (nn), BC",
		Length:    4,
		Cycles:    20,
		AddrMode1: AddrModeExtended,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			cpu.Memory.Write(nn, cpu.Regs.C)
			cpu.Memory.Write(nn+1, cpu.Regs.B)
			return 20
		},
	})
	RegisterInstruction(&EDTable, 0x53, Instruction{
		Mnemonic:  "LD (nn), DE",
		Length:    4,
		Cycles:    20,
		AddrMode1: AddrModeExtended,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			cpu.Memory.Write(nn, cpu.Regs.E)
			cpu.Memory.Write(nn+1, cpu.Regs.D)
			return 20
		},
	})
	RegisterInstruction(&EDTable, 0x63, Instruction{
		Mnemonic:  "LD (nn), HL",
		Length:    4,
		Cycles:    20,
		AddrMode1: AddrModeExtended,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			cpu.Memory.Write(nn, cpu.Regs.L)
			cpu.Memory.Write(nn+1, cpu.Regs.H)
			return 20
		},
	})
	RegisterInstruction(&EDTable, 0x73, Instruction{
		Mnemonic:  "LD (nn), SP",
		Length:    4,
		Cycles:    20,
		AddrMode1: AddrModeExtended,
		AddrMode2: AddrModeRegisterPair,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			cpu.Memory.Write(nn, uint8(cpu.Regs.SP&0xFF))
			cpu.Memory.Write(nn+1, uint8(cpu.Regs.SP>>8))
			return 20
		},
	})

	// 57. LD BC, (nn) / LD DE, (nn) / LD HL, (nn) / LD SP, (nn)
	RegisterInstruction(&EDTable, 0x4B, Instruction{
		Mnemonic:  "LD BC, (nn)",
		Length:    4,
		Cycles:    20,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeExtended,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			low := cpu.Memory.Read(nn)
			high := cpu.Memory.Read(nn + 1)
			cpu.Regs.SetBC((uint16(high) << 8) | uint16(low))
			return 20
		},
	})
	RegisterInstruction(&EDTable, 0x5B, Instruction{
		Mnemonic:  "LD DE, (nn)",
		Length:    4,
		Cycles:    20,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeExtended,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			low := cpu.Memory.Read(nn)
			high := cpu.Memory.Read(nn + 1)
			cpu.Regs.SetDE((uint16(high) << 8) | uint16(low))
			return 20
		},
	})
	RegisterInstruction(&EDTable, 0x6B, Instruction{
		Mnemonic:  "LD HL, (nn)",
		Length:    4,
		Cycles:    20,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeExtended,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			low := cpu.Memory.Read(nn)
			high := cpu.Memory.Read(nn + 1)
			cpu.Regs.SetHL((uint16(high) << 8) | uint16(low))
			return 20
		},
	})
	RegisterInstruction(&EDTable, 0x7B, Instruction{
		Mnemonic:  "LD SP, (nn)",
		Length:    4,
		Cycles:    20,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeExtended,
		Execute: func(cpu *CPU) int {
			nn := cpu.FetchWord()
			low := cpu.Memory.Read(nn)
			high := cpu.Memory.Read(nn + 1)
			cpu.Regs.SP = (uint16(high) << 8) | uint16(low)
			return 20
		},
	})

	// 58. LD I, A / LD R, A / LD A, I / LD A, R
	RegisterInstruction(&EDTable, 0x47, Instruction{
		Mnemonic:  "LD I, A",
		Length:    2,
		Cycles:    9,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeAccumulator,
		Execute: func(cpu *CPU) int {
			cpu.Regs.I = cpu.Regs.A
			return 9
		},
	})
	RegisterInstruction(&EDTable, 0x4F, Instruction{
		Mnemonic:  "LD R, A",
		Length:    2,
		Cycles:    9,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeAccumulator,
		Execute: func(cpu *CPU) int {
			cpu.Regs.R = cpu.Regs.A
			return 9
		},
	})
	RegisterInstruction(&EDTable, 0x57, Instruction{
		Mnemonic:  "LD A, I",
		Length:    2,
		Cycles:    9,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.A = cpu.Regs.I
			cpu.Regs.SetFlag(FlagS, (cpu.Regs.A&0x80) != 0)
			cpu.Regs.SetFlag(FlagZ, cpu.Regs.A == 0)
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagPV, cpu.IFF2)
			cpu.Regs.SetFlag(FlagN, false)
			cpu.Regs.SetFlag(Flag3, (cpu.Regs.A&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (cpu.Regs.A&0x20) != 0)
			return 9
		},
	})
	RegisterInstruction(&EDTable, 0x5F, Instruction{
		Mnemonic:  "LD A, R",
		Length:    2,
		Cycles:    9,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeRegister,
		Execute: func(cpu *CPU) int {
			cpu.Regs.A = cpu.Regs.R
			cpu.Regs.SetFlag(FlagS, (cpu.Regs.A&0x80) != 0)
			cpu.Regs.SetFlag(FlagZ, cpu.Regs.A == 0)
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagPV, cpu.IFF2)
			cpu.Regs.SetFlag(FlagN, false)
			cpu.Regs.SetFlag(Flag3, (cpu.Regs.A&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (cpu.Regs.A&0x20) != 0)
			return 9
		},
	})
}
