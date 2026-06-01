package z80

// initBIT registers bit manipulation instructions.
func initBIT() {
	// 0x2F: CPL (Complement Accumulator)
	RegisterInstruction(&MainTable, 0x2F, Instruction{
		Mnemonic:  "CPL",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			cpu.Regs.A = ^cpu.Regs.A
			cpu.Regs.SetFlag(FlagH, true)
			cpu.Regs.SetFlag(FlagN, true)
			// Undocumented flags 3 and 5 are copied from A
			cpu.Regs.SetFlag(Flag3, (cpu.Regs.A&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (cpu.Regs.A&0x20) != 0)
			return 4
		},
	})

	// --- CB Prefix Instructions ---

	regs8 := []struct {
		name string
		get  func(r *Registers) uint8
		set  func(r *Registers, val uint8)
		op   uint8
	}{
		{"B", func(r *Registers) uint8 { return r.B }, func(r *Registers, val uint8) { r.B = val }, 0},
		{"C", func(r *Registers) uint8 { return r.C }, func(r *Registers, val uint8) { r.C = val }, 1},
		{"D", func(r *Registers) uint8 { return r.D }, func(r *Registers, val uint8) { r.D = val }, 2},
		{"E", func(r *Registers) uint8 { return r.E }, func(r *Registers, val uint8) { r.E = val }, 3},
		{"H", func(r *Registers) uint8 { return r.H }, func(r *Registers, val uint8) { r.H = val }, 4},
		{"L", func(r *Registers) uint8 { return r.L }, func(r *Registers, val uint8) { r.L = val }, 5},
		{"(HL)", nil, nil, 6}, // Special handling
		{"A", func(r *Registers) uint8 { return r.A }, func(r *Registers, val uint8) { r.A = val }, 7},
	}

	for bit := uint8(0); bit < 8; bit++ {
		for _, r := range regs8 {
			r := r // capture
			bit := bit

			// BIT b, r
			opBIT := 0x40 | (bit << 3) | r.op
			RegisterInstruction(&CBTable, opBIT, Instruction{
				Mnemonic: "BIT " + string('0'+bit) + ", " + r.name,
				Length:   2,
				Cycles:   8,
				Execute: func(cpu *CPU) int {
					var val uint8
					cycles := 8
					if r.name == "(HL)" {
						val = cpu.Memory.Read(cpu.Regs.HL())
						cycles = 12
					} else {
						val = r.get(cpu.Regs)
					}
					res := val & (1 << bit)
					cpu.Regs.SetFlag(FlagZ, res == 0)
					cpu.Regs.SetFlag(FlagH, true)
					cpu.Regs.SetFlag(FlagN, false)
					cpu.Regs.SetFlag(FlagS, bit == 7 && res != 0)
					cpu.Regs.SetFlag(FlagPV, res == 0)

					// Undocumented flags 3 and 5
					if r.name == "(HL)" {
						// For (HL), it's often based on high byte of HL or MEMPTR
						cpu.Regs.SetFlag(Flag3, (cpu.Regs.H&0x08) != 0)
						cpu.Regs.SetFlag(Flag5, (cpu.Regs.H&0x20) != 0)
					} else {
						cpu.Regs.SetFlag(Flag3, (val&0x08) != 0)
						cpu.Regs.SetFlag(Flag5, (val&0x20) != 0)
					}
					return cycles
				},
			})

			// RES b, r
			opRES := 0x80 | (bit << 3) | r.op
			RegisterInstruction(&CBTable, opRES, Instruction{
				Mnemonic: "RES " + string('0'+bit) + ", " + r.name,
				Length:   2,
				Cycles:   8,
				Execute: func(cpu *CPU) int {
					var val uint8
					cycles := 8
					if r.name == "(HL)" {
						val = cpu.Memory.Read(cpu.Regs.HL())
						val &= ^(1 << bit)
						cpu.Memory.Write(cpu.Regs.HL(), val)
						cycles = 15
					} else {
						val = r.get(cpu.Regs)
						val &= ^(1 << bit)
						r.set(cpu.Regs, val)
					}
					return cycles
				},
			})

			// SET b, r
			opSET := 0xC0 | (bit << 3) | r.op
			RegisterInstruction(&CBTable, opSET, Instruction{
				Mnemonic: "SET " + string('0'+bit) + ", " + r.name,
				Length:   2,
				Cycles:   8,
				Execute: func(cpu *CPU) int {
					var val uint8
					cycles := 8
					if r.name == "(HL)" {
						val = cpu.Memory.Read(cpu.Regs.HL())
						val |= (1 << bit)
						cpu.Memory.Write(cpu.Regs.HL(), val)
						cycles = 15
					} else {
						val = r.get(cpu.Regs)
						val |= (1 << bit)
						r.set(cpu.Regs, val)
					}
					return cycles
				},
			})
		}
	}

	// --- DDCB and FDCB Prefix Instructions ---

	regs8Idx := []struct {
		name string
		set  func(r *Registers, val uint8)
		op   uint8
	}{
		{"B", func(r *Registers, val uint8) { r.B = val }, 0},
		{"C", func(r *Registers, val uint8) { r.C = val }, 1},
		{"D", func(r *Registers, val uint8) { r.D = val }, 2},
		{"E", func(r *Registers, val uint8) { r.E = val }, 3},
		{"H", func(r *Registers, val uint8) { r.H = val }, 4},
		{"L", func(r *Registers, val uint8) { r.L = val }, 5},
		{"", nil, 6},
		{"A", func(r *Registers, val uint8) { r.A = val }, 7},
	}

	for bit := uint8(0); bit < 8; bit++ {
		for _, r := range regs8Idx {
			bit := bit
			r := r

			// BIT b, (IX+d)
			opBIT := 0x40 | (bit << 3) | r.op
			RegisterInstruction(&DDCBTable, opBIT, Instruction{
				Mnemonic: "BIT " + string('0'+bit) + ", (IX+d)",
				Length:   4,
				Cycles:   20,
				Execute: func(cpu *CPU) int {
					addr := uint16(int32(cpu.Regs.IX) + int32(cpu.LastDisplacement))
					val := cpu.Memory.Read(addr)
					res := val & (1 << bit)
					cpu.Regs.SetFlag(FlagZ, res == 0)
					cpu.Regs.SetFlag(FlagH, true)
					cpu.Regs.SetFlag(FlagN, false)
					cpu.Regs.SetFlag(FlagS, bit == 7 && res != 0)
					cpu.Regs.SetFlag(FlagPV, res == 0)
					// Undocumented flags 3 and 5 are typically from high byte of address
					cpu.Regs.SetFlag(Flag3, (uint8(addr>>8)&0x08) != 0)
					cpu.Regs.SetFlag(Flag5, (uint8(addr>>8)&0x20) != 0)
					return 20
				},
			})

			// BIT b, (IY+d)
			RegisterInstruction(&FDCBTable, opBIT, Instruction{
				Mnemonic: "BIT " + string('0'+bit) + ", (IY+d)",
				Length:   4,
				Cycles:   20,
				Execute: func(cpu *CPU) int {
					addr := uint16(int32(cpu.Regs.IY) + int32(cpu.LastDisplacement))
					val := cpu.Memory.Read(addr)
					res := val & (1 << bit)
					cpu.Regs.SetFlag(FlagZ, res == 0)
					cpu.Regs.SetFlag(FlagH, true)
					cpu.Regs.SetFlag(FlagN, false)
					cpu.Regs.SetFlag(FlagS, bit == 7 && res != 0)
					cpu.Regs.SetFlag(FlagPV, res == 0)
					// Undocumented flags 3 and 5 are typically from high byte of address
					cpu.Regs.SetFlag(Flag3, (uint8(addr>>8)&0x08) != 0)
					cpu.Regs.SetFlag(Flag5, (uint8(addr>>8)&0x20) != 0)
					return 20
				},
			})

			// RES b, (IX+d)
			opRES := 0x80 | (bit << 3) | r.op
			mnemonicRES_DD := "RES " + string('0'+bit) + ", (IX+d)"
			if r.name != "" {
				mnemonicRES_DD += ", " + r.name
			}
			RegisterInstruction(&DDCBTable, opRES, Instruction{
				Mnemonic: mnemonicRES_DD,
				Length:   4,
				Cycles:   23,
				Execute: func(cpu *CPU) int {
					addr := uint16(int32(cpu.Regs.IX) + int32(cpu.LastDisplacement))
					val := cpu.Memory.Read(addr)
					val &= ^(1 << bit)
					cpu.Memory.Write(addr, val)
					if r.set != nil {
						r.set(cpu.Regs, val)
					}
					return 23
				},
			})

			// RES b, (IY+d)
			mnemonicRES_FD := "RES " + string('0'+bit) + ", (IY+d)"
			if r.name != "" {
				mnemonicRES_FD += ", " + r.name
			}
			RegisterInstruction(&FDCBTable, opRES, Instruction{
				Mnemonic: mnemonicRES_FD,
				Length:   4,
				Cycles:   23,
				Execute: func(cpu *CPU) int {
					addr := uint16(int32(cpu.Regs.IY) + int32(cpu.LastDisplacement))
					val := cpu.Memory.Read(addr)
					val &= ^(1 << bit)
					cpu.Memory.Write(addr, val)
					if r.set != nil {
						r.set(cpu.Regs, val)
					}
					return 23
				},
			})

			// SET b, (IX+d)
			opSET := 0xC0 | (bit << 3) | r.op
			mnemonicSET_DD := "SET " + string('0'+bit) + ", (IX+d)"
			if r.name != "" {
				mnemonicSET_DD += ", " + r.name
			}
			RegisterInstruction(&DDCBTable, opSET, Instruction{
				Mnemonic: mnemonicSET_DD,
				Length:   4,
				Cycles:   23,
				Execute: func(cpu *CPU) int {
					addr := uint16(int32(cpu.Regs.IX) + int32(cpu.LastDisplacement))
					val := cpu.Memory.Read(addr)
					val |= (1 << bit)
					cpu.Memory.Write(addr, val)
					if r.set != nil {
						r.set(cpu.Regs, val)
					}
					return 23
				},
			})

			// SET b, (IY+d)
			mnemonicSET_FD := "SET " + string('0'+bit) + ", (IY+d)"
			if r.name != "" {
				mnemonicSET_FD += ", " + r.name
			}
			RegisterInstruction(&FDCBTable, opSET, Instruction{
				Mnemonic: mnemonicSET_FD,
				Length:   4,
				Cycles:   23,
				Execute: func(cpu *CPU) int {
					addr := uint16(int32(cpu.Regs.IY) + int32(cpu.LastDisplacement))
					val := cpu.Memory.Read(addr)
					val |= (1 << bit)
					cpu.Memory.Write(addr, val)
					if r.set != nil {
						r.set(cpu.Regs, val)
					}
					return 23
				},
			})
		}
	}
}
