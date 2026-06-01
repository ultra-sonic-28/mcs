package z80

// initPushPop registers PUSH and POP instructions.
func initPushPop() {
	// --- PUSH Instructions ---

	// PUSH rr (Standard)
	pushRegs := []struct {
		name string
		get  func(r *Registers) uint16
		op   uint8
	}{
		{"BC", func(r *Registers) uint16 { return r.BC() }, 0xC5},
		{"DE", func(r *Registers) uint16 { return r.DE() }, 0xD5},
		{"HL", func(r *Registers) uint16 { return r.HL() }, 0xE5},
		{"AF", func(r *Registers) uint16 { return r.AF() }, 0xF5},
	}

	for _, r := range pushRegs {
		r := r // capture for closure
		RegisterInstruction(&MainTable, r.op, Instruction{
			Mnemonic:  "PUSH " + r.name,
			Length:    1,
			Cycles:    11,
			AddrMode1: AddrModeRegisterPair,
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				val := r.get(cpu.Regs)
				cpu.Regs.SP--
				cpu.Memory.Write(cpu.Regs.SP, uint8(val>>8))
				cpu.Regs.SP--
				cpu.Memory.Write(cpu.Regs.SP, uint8(val&0xFF))
				return 11
			},
		})
	}

	// PUSH IX (0xDD 0xE5)
	RegisterInstruction(&DDTable, 0xE5, Instruction{
		Mnemonic:  "PUSH IX",
		Length:    2,
		Cycles:    15,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IX
			cpu.Regs.SP--
			cpu.Memory.Write(cpu.Regs.SP, uint8(val>>8))
			cpu.Regs.SP--
			cpu.Memory.Write(cpu.Regs.SP, uint8(val&0xFF))
			return 15
		},
	})

	// PUSH IY (0xFD 0xE5)
	RegisterInstruction(&FDTable, 0xE5, Instruction{
		Mnemonic:  "PUSH IY",
		Length:    2,
		Cycles:    15,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			val := cpu.Regs.IY
			cpu.Regs.SP--
			cpu.Memory.Write(cpu.Regs.SP, uint8(val>>8))
			cpu.Regs.SP--
			cpu.Memory.Write(cpu.Regs.SP, uint8(val&0xFF))
			return 15
		},
	})

	// --- POP Instructions ---

	// POP rr (Standard)
	popRegs := []struct {
		name string
		set  func(r *Registers, val uint16)
		op   uint8
	}{
		{"BC", func(r *Registers, val uint16) { r.SetBC(val) }, 0xC1},
		{"DE", func(r *Registers, val uint16) { r.SetDE(val) }, 0xD1},
		{"HL", func(r *Registers, val uint16) { r.SetHL(val) }, 0xE1},
		{"AF", func(r *Registers, val uint16) { r.SetAF(val) }, 0xF1},
	}

	for _, r := range popRegs {
		r := r // capture for closure
		RegisterInstruction(&MainTable, r.op, Instruction{
			Mnemonic:  "POP " + r.name,
			Length:    1,
			Cycles:    10,
			AddrMode1: AddrModeRegisterPair,
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				low := uint16(cpu.Memory.Read(cpu.Regs.SP))
				cpu.Regs.SP++
				high := uint16(cpu.Memory.Read(cpu.Regs.SP))
				cpu.Regs.SP++
				r.set(cpu.Regs, (high<<8)|low)
				return 10
			},
		})
	}

	// POP IX (0xDD 0xE1)
	RegisterInstruction(&DDTable, 0xE1, Instruction{
		Mnemonic:  "POP IX",
		Length:    2,
		Cycles:    14,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			low := uint16(cpu.Memory.Read(cpu.Regs.SP))
			cpu.Regs.SP++
			high := uint16(cpu.Memory.Read(cpu.Regs.SP))
			cpu.Regs.SP++
			cpu.Regs.IX = (high << 8) | low
			return 14
		},
	})

	// POP IY (0xFD 0xE1)
	RegisterInstruction(&FDTable, 0xE1, Instruction{
		Mnemonic:  "POP IY",
		Length:    2,
		Cycles:    14,
		AddrMode1: AddrModeRegisterPair,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			low := uint16(cpu.Memory.Read(cpu.Regs.SP))
			cpu.Regs.SP++
			high := uint16(cpu.Memory.Read(cpu.Regs.SP))
			cpu.Regs.SP++
			cpu.Regs.IY = (high << 8) | low
			return 14
		},
	})
}
