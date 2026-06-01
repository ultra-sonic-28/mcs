package z80

// initIncDec registers 8-bit Increment and Decrement instructions.
func initIncDec() {
	// INC r (8-bit)
	// Opcodes: 0x04 (B), 0x0C (C), 0x14 (D), 0x1C (E), 0x24 (H), 0x2C (L), 0x3C (A)
	regsInc := []struct {
		name string
		get  func(r *Registers) uint8
		set  func(r *Registers, val uint8)
		op   uint8
	}{
		{"B", func(r *Registers) uint8 { return r.B }, func(r *Registers, val uint8) { r.B = val }, 0x04},
		{"C", func(r *Registers) uint8 { return r.C }, func(r *Registers, val uint8) { r.C = val }, 0x0C},
		{"D", func(r *Registers) uint8 { return r.D }, func(r *Registers, val uint8) { r.D = val }, 0x14},
		{"E", func(r *Registers) uint8 { return r.E }, func(r *Registers, val uint8) { r.E = val }, 0x1C},
		{"H", func(r *Registers) uint8 { return r.H }, func(r *Registers, val uint8) { r.H = val }, 0x24},
		{"L", func(r *Registers) uint8 { return r.L }, func(r *Registers, val uint8) { r.L = val }, 0x2C},
		{"A", func(r *Registers) uint8 { return r.A }, func(r *Registers, val uint8) { r.A = val }, 0x3C},
	}

	for _, r := range regsInc {
		r := r // capture for closure
		RegisterInstruction(&MainTable, r.op, Instruction{
			Mnemonic:  "INC " + r.name,
			Length:    1,
			Cycles:    4,
			AddrMode1: AddrModeRegister,
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				oldVal := r.get(cpu.Regs)
				res := oldVal + 1
				r.set(cpu.Regs, res)
				cpu.Regs.UpdateFlagsInc8(oldVal, res)
				return 4
			},
		})
	}

	// DEC r (8-bit)
	// Opcodes: 0x05 (B), 0x0D (C), 0x15 (D), 0x1D (E), 0x25 (H), 0x2D (L), 0x3D (A)
	regsDec := []struct {
		name string
		get  func(r *Registers) uint8
		set  func(r *Registers, val uint8)
		op   uint8
	}{
		{"B", func(r *Registers) uint8 { return r.B }, func(r *Registers, val uint8) { r.B = val }, 0x05},
		{"C", func(r *Registers) uint8 { return r.C }, func(r *Registers, val uint8) { r.C = val }, 0x0D},
		{"D", func(r *Registers) uint8 { return r.D }, func(r *Registers, val uint8) { r.D = val }, 0x15},
		{"E", func(r *Registers) uint8 { return r.E }, func(r *Registers, val uint8) { r.E = val }, 0x1D},
		{"H", func(r *Registers) uint8 { return r.H }, func(r *Registers, val uint8) { r.H = val }, 0x25},
		{"L", func(r *Registers) uint8 { return r.L }, func(r *Registers, val uint8) { r.L = val }, 0x2D},
		{"A", func(r *Registers) uint8 { return r.A }, func(r *Registers, val uint8) { r.A = val }, 0x3D},
	}

	for _, r := range regsDec {
		r := r // capture for closure
		RegisterInstruction(&MainTable, r.op, Instruction{
			Mnemonic:  "DEC " + r.name,
			Length:    1,
			Cycles:    4,
			AddrMode1: AddrModeRegister,
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				oldVal := r.get(cpu.Regs)
				res := oldVal - 1
				r.set(cpu.Regs, res)
				cpu.Regs.UpdateFlagsDec8(oldVal, res)
				return 4
			},
		})
	}

	// 0x34: INC (HL)
	RegisterInstruction(&MainTable, 0x34, Instruction{
		Mnemonic:  "INC (HL)",
		Length:    1,
		Cycles:    11,
		AddrMode1: AddrModeIndirect,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			addr := cpu.Regs.HL()
			oldVal := cpu.Memory.Read(addr)
			res := oldVal + 1
			cpu.Memory.Write(addr, res)
			cpu.Regs.UpdateFlagsInc8(oldVal, res)
			return 11
		},
	})

	// 0x35: DEC (HL)
	RegisterInstruction(&MainTable, 0x35, Instruction{
		Mnemonic:  "DEC (HL)",
		Length:    1,
		Cycles:    11,
		AddrMode1: AddrModeIndirect,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			addr := cpu.Regs.HL()
			oldVal := cpu.Memory.Read(addr)
			res := oldVal - 1
			cpu.Memory.Write(addr, res)
			cpu.Regs.UpdateFlagsDec8(oldVal, res)
			return 11
		},
	})

	// INC rr (16-bit)
	// Opcodes: 0x03 (BC), 0x13 (DE), 0x23 (HL), 0x33 (SP)
	// Note: 16-bit INC does NOT affect flags.
	regsInc16 := []struct {
		name string
		get  func(r *Registers) uint16
		set  func(r *Registers, val uint16)
		op   uint8
	}{
		{"BC", func(r *Registers) uint16 { return r.BC() }, func(r *Registers, val uint16) { r.SetBC(val) }, 0x03},
		{"DE", func(r *Registers) uint16 { return r.DE() }, func(r *Registers, val uint16) { r.SetDE(val) }, 0x13},
		{"HL", func(r *Registers) uint16 { return r.HL() }, func(r *Registers, val uint16) { r.SetHL(val) }, 0x23},
		{"SP", func(r *Registers) uint16 { return r.SP }, func(r *Registers, val uint16) { r.SP = val }, 0x33},
	}

	for _, r := range regsInc16 {
		r := r // capture for closure
		RegisterInstruction(&MainTable, r.op, Instruction{
			Mnemonic:  "INC " + r.name,
			Length:    1,
			Cycles:    6,
			AddrMode1: AddrModeRegister,
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				r.set(cpu.Regs, r.get(cpu.Regs)+1)
				return 6
			},
		})
	}

	// DEC rr (16-bit)
	// Opcodes: 0x0B (BC), 0x1B (DE), 0x2B (HL), 0x3B (SP)
	// Note: 16-bit DEC does NOT affect flags.
	regsDec16 := []struct {
		name string
		get  func(r *Registers) uint16
		set  func(r *Registers, val uint16)
		op   uint8
	}{
		{"BC", func(r *Registers) uint16 { return r.BC() }, func(r *Registers, val uint16) { r.SetBC(val) }, 0x0B},
		{"DE", func(r *Registers) uint16 { return r.DE() }, func(r *Registers, val uint16) { r.SetDE(val) }, 0x1B},
		{"HL", func(r *Registers) uint16 { return r.HL() }, func(r *Registers, val uint16) { r.SetHL(val) }, 0x2B},
		{"SP", func(r *Registers) uint16 { return r.SP }, func(r *Registers, val uint16) { r.SP = val }, 0x3B},
	}

	for _, r := range regsDec16 {
		r := r // capture for closure
		RegisterInstruction(&MainTable, r.op, Instruction{
			Mnemonic:  "DEC " + r.name,
			Length:    1,
			Cycles:    6,
			AddrMode1: AddrModeRegister,
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				r.set(cpu.Regs, r.get(cpu.Regs)-1)
				return 6
			},
		})
	}

	// INC IX / DEC IX
	RegisterInstruction(&DDTable, 0x23, Instruction{
		Mnemonic:  "INC IX",
		Length:    2,
		Cycles:    10,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.Regs.IX++
			return 10
		},
	})
	RegisterInstruction(&DDTable, 0x2B, Instruction{
		Mnemonic:  "DEC IX",
		Length:    2,
		Cycles:    10,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.Regs.IX--
			return 10
		},
	})

	// INC IXH / DEC IXH / INC IXL / DEC IXL
	RegisterInstruction(&DDTable, 0x24, Instruction{
		Mnemonic:  "INC IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			oldVal := cpu.Regs.IXH()
			res := oldVal + 1
			cpu.Regs.SetIXH(res)
			cpu.Regs.UpdateFlagsInc8(oldVal, res)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x25, Instruction{
		Mnemonic:  "DEC IXH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			oldVal := cpu.Regs.IXH()
			res := oldVal - 1
			cpu.Regs.SetIXH(res)
			cpu.Regs.UpdateFlagsDec8(oldVal, res)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x2C, Instruction{
		Mnemonic:  "INC IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			oldVal := cpu.Regs.IXL()
			res := oldVal + 1
			cpu.Regs.SetIXL(res)
			cpu.Regs.UpdateFlagsInc8(oldVal, res)
			return 8
		},
	})
	RegisterInstruction(&DDTable, 0x2D, Instruction{
		Mnemonic:  "DEC IXL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			oldVal := cpu.Regs.IXL()
			res := oldVal - 1
			cpu.Regs.SetIXL(res)
			cpu.Regs.UpdateFlagsDec8(oldVal, res)
			return 8
		},
	})

	// 0x34: INC (IX+d)
	RegisterInstruction(&DDTable, 0x34, Instruction{
		Mnemonic:  "INC (IX+d)",
		Length:    3,
		Cycles:    23,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			oldVal := cpu.Memory.Read(addr)
			res := oldVal + 1
			cpu.Memory.Write(addr, res)
			cpu.Regs.UpdateFlagsInc8(oldVal, res)
			return 23
		},
	})

	// 0x35: DEC (IX+d)
	RegisterInstruction(&DDTable, 0x35, Instruction{
		Mnemonic:  "DEC (IX+d)",
		Length:    3,
		Cycles:    23,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IX) + int32(d))
			oldVal := cpu.Memory.Read(addr)
			res := oldVal - 1
			cpu.Memory.Write(addr, res)
			cpu.Regs.UpdateFlagsDec8(oldVal, res)
			return 23
		},
	})

	// --- IY Increment/Decrement Instructions (prefixed with 0xFD) ---

	// INC IY / DEC IY
	RegisterInstruction(&FDTable, 0x23, Instruction{
		Mnemonic:  "INC IY",
		Length:    2,
		Cycles:    10,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.Regs.IY++
			return 10
		},
	})
	RegisterInstruction(&FDTable, 0x2B, Instruction{
		Mnemonic:  "DEC IY",
		Length:    2,
		Cycles:    10,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.Regs.IY--
			return 10
		},
	})

	// INC IYH / DEC IYH / INC IYL / DEC IYL
	RegisterInstruction(&FDTable, 0x24, Instruction{
		Mnemonic:  "INC IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			oldVal := cpu.Regs.IYH()
			res := oldVal + 1
			cpu.Regs.SetIYH(res)
			cpu.Regs.UpdateFlagsInc8(oldVal, res)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x25, Instruction{
		Mnemonic:  "DEC IYH",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			oldVal := cpu.Regs.IYH()
			res := oldVal - 1
			cpu.Regs.SetIYH(res)
			cpu.Regs.UpdateFlagsDec8(oldVal, res)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x2C, Instruction{
		Mnemonic:  "INC IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			oldVal := cpu.Regs.IYL()
			res := oldVal + 1
			cpu.Regs.SetIYL(res)
			cpu.Regs.UpdateFlagsInc8(oldVal, res)
			return 8
		},
	})
	RegisterInstruction(&FDTable, 0x2D, Instruction{
		Mnemonic:  "DEC IYL",
		Length:    2,
		Cycles:    8,
		AddrMode1: AddrModeRegister,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			oldVal := cpu.Regs.IYL()
			res := oldVal - 1
			cpu.Regs.SetIYL(res)
			cpu.Regs.UpdateFlagsDec8(oldVal, res)
			return 8
		},
	})

	// 0x34: INC (IY+d)
	RegisterInstruction(&FDTable, 0x34, Instruction{
		Mnemonic:  "INC (IY+d)",
		Length:    3,
		Cycles:    23,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			oldVal := cpu.Memory.Read(addr)
			res := oldVal + 1
			cpu.Memory.Write(addr, res)
			cpu.Regs.UpdateFlagsInc8(oldVal, res)
			return 23
		},
	})

	// 0x35: DEC (IY+d)
	RegisterInstruction(&FDTable, 0x35, Instruction{
		Mnemonic:  "DEC (IY+d)",
		Length:    3,
		Cycles:    23,
		AddrMode1: AddrModeIndexed,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			d := int8(cpu.FetchByte())
			addr := uint16(int32(cpu.Regs.IY) + int32(d))
			oldVal := cpu.Memory.Read(addr)
			res := oldVal - 1
			cpu.Memory.Write(addr, res)
			cpu.Regs.UpdateFlagsDec8(oldVal, res)
			return 23
		},
	})
}
