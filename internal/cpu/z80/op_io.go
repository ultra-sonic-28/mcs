package z80

// initIO registers I/O instructions.
func initIO() {
	// 0xD3: OUT (n), A
	RegisterInstruction(&MainTable, 0xD3, Instruction{
		Mnemonic:  "OUT (n), A",
		Length:    2,
		Cycles:    11,
		AddrMode1: AddrModePort,
		AddrMode2: AddrModeAccumulator,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			port := (uint16(cpu.Regs.A) << 8) | uint16(n)
			cpu.IO.Out(port, cpu.Regs.A)
			return 11
		},
	})

	// 0xDB: IN A, (n)
	RegisterInstruction(&MainTable, 0xDB, Instruction{
		Mnemonic:  "IN A, (n)",
		Length:    2,
		Cycles:    11,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModePort,
		Execute: func(cpu *CPU) int {
			n := cpu.FetchByte()
			port := (uint16(cpu.Regs.A) << 8) | uint16(n)
			cpu.Regs.A = cpu.IO.In(port)
			return 11
		},
	})

	// --- ED Prefix I/O Instructions ---

	// IN r, (C)
	inRegs := []struct {
		name string
		set  func(r *Registers, val uint8)
		op   uint8
	}{
		{"B", func(r *Registers, val uint8) { r.B = val }, 0x40},
		{"C", func(r *Registers, val uint8) { r.C = val }, 0x48},
		{"D", func(r *Registers, val uint8) { r.D = val }, 0x50},
		{"E", func(r *Registers, val uint8) { r.E = val }, 0x58},
		{"H", func(r *Registers, val uint8) { r.H = val }, 0x60},
		{"L", func(r *Registers, val uint8) { r.L = val }, 0x68},
		{"(C)", nil, 0x70}, // Special case: affects flags only
		{"A", func(r *Registers, val uint8) { r.A = val }, 0x78},
	}

	for _, r := range inRegs {
		r := r // capture
		mnemonic := "IN " + r.name
		if r.set == nil {
			mnemonic = "IN (C)"
		} else {
			mnemonic = "IN " + r.name + ", (C)"
		}

		RegisterInstruction(&EDTable, r.op, Instruction{
			Mnemonic:  mnemonic,
			Length:    2,
			Cycles:    12,
			AddrMode1: AddrModeRegister,
			AddrMode2: AddrModePort,
			Execute: func(cpu *CPU) int {
				port := cpu.Regs.BC()
				val := cpu.IO.In(port)
				if r.set != nil {
					r.set(cpu.Regs, val)
				}
				cpu.Regs.UpdateFlagsIOIn(val)
				return 12
			},
		})
	}

	// OUT (C), r
	outRegs := []struct {
		name string
		get  func(r *Registers) uint8
		op   uint8
	}{
		{"B", func(r *Registers) uint8 { return r.B }, 0x41},
		{"C", func(r *Registers) uint8 { return r.C }, 0x49},
		{"D", func(r *Registers) uint8 { return r.D }, 0x51},
		{"E", func(r *Registers) uint8 { return r.E }, 0x59},
		{"H", func(r *Registers) uint8 { return r.H }, 0x61},
		{"L", func(r *Registers) uint8 { return r.L }, 0x69},
		{"0", func(r *Registers) uint8 { return 0 }, 0x71}, // Special case: outputs 0
		{"A", func(r *Registers) uint8 { return r.A }, 0x79},
	}

	for _, r := range outRegs {
		r := r // capture
		mnemonic := "OUT (C), " + r.name

		RegisterInstruction(&EDTable, r.op, Instruction{
			Mnemonic:  mnemonic,
			Length:    2,
			Cycles:    12,
			AddrMode1: AddrModePort,
			AddrMode2: AddrModeRegister,
			Execute: func(cpu *CPU) int {
				port := cpu.Regs.BC()
				val := r.get(cpu.Regs)
				cpu.IO.Out(port, val)
				return 12
			},
		})
	}
}
