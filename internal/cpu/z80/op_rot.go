package z80

// initRot registers Rotation and Shift instructions.
func initRot() {
	// 0x07: RLCA (Rotate Left Circular Accumulator)
	RegisterInstruction(&MainTable, 0x07, Instruction{
		Mnemonic:  "RLCA",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			a := cpu.Regs.A
			carry := a >> 7
			res := (a << 1) | carry
			cpu.Regs.A = res
			cpu.Regs.SetFlag(FlagC, carry != 0)
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagN, false)
			cpu.Regs.SetFlag(Flag3, (res&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (res&0x20) != 0)
			return 4
		},
	})

	// 0x0F: RRCA (Rotate Right Circular Accumulator)
	RegisterInstruction(&MainTable, 0x0F, Instruction{
		Mnemonic:  "RRCA",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			a := cpu.Regs.A
			carry := a & 0x01
			res := (a >> 1) | (carry << 7)
			cpu.Regs.A = res
			cpu.Regs.SetFlag(FlagC, carry != 0)
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagN, false)
			cpu.Regs.SetFlag(Flag3, (res&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (res&0x20) != 0)
			return 4
		},
	})

	// 0x17: RLA (Rotate Left Accumulator through Carry)
	RegisterInstruction(&MainTable, 0x17, Instruction{
		Mnemonic:  "RLA",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			a := cpu.Regs.A
			oldCarry := uint8(0)
			if cpu.Regs.Flag(FlagC) {
				oldCarry = 1
			}
			newCarry := a >> 7
			res := (a << 1) | oldCarry
			cpu.Regs.A = res
			cpu.Regs.SetFlag(FlagC, newCarry != 0)
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagN, false)
			cpu.Regs.SetFlag(Flag3, (res&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (res&0x20) != 0)
			return 4
		},
	})

	// 0x1F: RRA (Rotate Right Accumulator through Carry)
	RegisterInstruction(&MainTable, 0x1F, Instruction{
		Mnemonic:  "RRA",
		Length:    1,
		Cycles:    4,
		AddrMode1: AddrModeAccumulator,
		AddrMode2: AddrModeNone,
		Execute: func(cpu *CPU) int {
			a := cpu.Regs.A
			oldCarry := uint8(0)
			if cpu.Regs.Flag(FlagC) {
				oldCarry = 1
			}
			newCarry := a & 0x01
			res := (a >> 1) | (oldCarry << 7)
			cpu.Regs.A = res
			cpu.Regs.SetFlag(FlagC, newCarry != 0)
			cpu.Regs.SetFlag(FlagH, false)
			cpu.Regs.SetFlag(FlagN, false)
			cpu.Regs.SetFlag(Flag3, (res&0x08) != 0)
			cpu.Regs.SetFlag(Flag5, (res&0x20) != 0)
			return 4
		},
	})

	// --- CB Prefix Shift/Rotate Instructions ---

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
		{"(HL)", nil, nil, 6},
		{"A", func(r *Registers) uint8 { return r.A }, func(r *Registers, val uint8) { r.A = val }, 7},
	}

	for _, r := range regs8 {
		r := r // capture

		// RLC r
		RegisterInstruction(&CBTable, 0x00|r.op, Instruction{
			Mnemonic: "RLC " + r.name,
			Length:   2,
			Cycles:   8,
			Execute: func(cpu *CPU) int {
				val, cycles := getVal8(cpu, r)
				newCarry := val >> 7
				res := (val << 1) | newCarry
				setVal8(cpu, r, res)
				cpu.Regs.UpdateFlagsIOIn(res)
				cpu.Regs.SetFlag(FlagC, newCarry != 0)
				return cycles
			},
		})

		// RRC r
		RegisterInstruction(&CBTable, 0x08|r.op, Instruction{
			Mnemonic: "RRC " + r.name,
			Length:   2,
			Cycles:   8,
			Execute: func(cpu *CPU) int {
				val, cycles := getVal8(cpu, r)
				newCarry := val & 0x01
				res := (val >> 1) | (newCarry << 7)
				setVal8(cpu, r, res)
				cpu.Regs.UpdateFlagsIOIn(res)
				cpu.Regs.SetFlag(FlagC, newCarry != 0)
				return cycles
			},
		})

		// RL r
		RegisterInstruction(&CBTable, 0x10|r.op, Instruction{
			Mnemonic: "RL " + r.name,
			Length:   2,
			Cycles:   8,
			Execute: func(cpu *CPU) int {
				val, cycles := getVal8(cpu, r)
				oldCarry := uint8(0)
				if cpu.Regs.Flag(FlagC) {
					oldCarry = 1
				}
				newCarry := val >> 7
				res := (val << 1) | oldCarry
				setVal8(cpu, r, res)
				cpu.Regs.UpdateFlagsIOIn(res)
				cpu.Regs.SetFlag(FlagC, newCarry != 0)
				return cycles
			},
		})

		// RR r
		RegisterInstruction(&CBTable, 0x18|r.op, Instruction{
			Mnemonic: "RR " + r.name,
			Length:   2,
			Cycles:   8,
			Execute: func(cpu *CPU) int {
				val, cycles := getVal8(cpu, r)
				oldCarry := uint8(0)
				if cpu.Regs.Flag(FlagC) {
					oldCarry = 1
				}
				newCarry := val & 0x01
				res := (val >> 1) | (oldCarry << 7)
				setVal8(cpu, r, res)
				cpu.Regs.UpdateFlagsIOIn(res)
				cpu.Regs.SetFlag(FlagC, newCarry != 0)
				return cycles
			},
		})

		// SLA r
		RegisterInstruction(&CBTable, 0x20|r.op, Instruction{
			Mnemonic: "SLA " + r.name,
			Length:   2,
			Cycles:   8,
			Execute: func(cpu *CPU) int {
				val, cycles := getVal8(cpu, r)
				newCarry := val >> 7
				res := val << 1
				setVal8(cpu, r, res)
				cpu.Regs.UpdateFlagsIOIn(res)
				cpu.Regs.SetFlag(FlagC, newCarry != 0)
				return cycles
			},
		})

		// SRA r
		RegisterInstruction(&CBTable, 0x28|r.op, Instruction{
			Mnemonic: "SRA " + r.name,
			Length:   2,
			Cycles:   8,
			Execute: func(cpu *CPU) int {
				val, cycles := getVal8(cpu, r)
				newCarry := val & 0x01
				res := (val >> 1) | (val & 0x80)
				setVal8(cpu, r, res)
				cpu.Regs.UpdateFlagsIOIn(res)
				cpu.Regs.SetFlag(FlagC, newCarry != 0)
				return cycles
			},
		})

		// SLL r (undocumented)
		RegisterInstruction(&CBTable, 0x30|r.op, Instruction{
			Mnemonic: "SLL " + r.name,
			Length:   2,
			Cycles:   8,
			Execute: func(cpu *CPU) int {
				val, cycles := getVal8(cpu, r)
				newCarry := val >> 7
				res := (val << 1) | 1
				setVal8(cpu, r, res)
				cpu.Regs.UpdateFlagsIOIn(res)
				cpu.Regs.SetFlag(FlagC, newCarry != 0)
				return cycles
			},
		})

		// SRL r
		RegisterInstruction(&CBTable, 0x38|r.op, Instruction{
			Mnemonic: "SRL " + r.name,
			Length:   2,
			Cycles:   8,
			Execute: func(cpu *CPU) int {
				val, cycles := getVal8(cpu, r)
				newCarry := val & 0x01
				res := val >> 1
				setVal8(cpu, r, res)
				cpu.Regs.UpdateFlagsIOIn(res)
				cpu.Regs.SetFlag(FlagC, newCarry != 0)
				return cycles
			},
		})
	}

	// --- DDCB and FDCB Prefix Shift/Rotate Instructions ---

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

	rotOps := []struct {
		name string
		op   uint8
		fn   func(cpu *CPU, val uint8) (uint8, bool)
	}{
		{"RLC", 0x00, func(cpu *CPU, val uint8) (uint8, bool) {
			c := val >> 7
			return (val << 1) | c, c != 0
		}},
		{"RRC", 0x08, func(cpu *CPU, val uint8) (uint8, bool) {
			c := val & 0x01
			return (val >> 1) | (c << 7), c != 0
		}},
		{"RL", 0x10, func(cpu *CPU, val uint8) (uint8, bool) {
			oldCarry := uint8(0)
			if cpu.Regs.Flag(FlagC) {
				oldCarry = 1
			}
			c := val >> 7
			return (val << 1) | oldCarry, c != 0
		}},
		{"RR", 0x18, func(cpu *CPU, val uint8) (uint8, bool) {
			oldCarry := uint8(0)
			if cpu.Regs.Flag(FlagC) {
				oldCarry = 1
			}
			c := val & 0x01
			return (val >> 1) | (oldCarry << 7), c != 0
		}},
		{"SLA", 0x20, func(cpu *CPU, val uint8) (uint8, bool) {
			c := val >> 7
			return val << 1, c != 0
		}},
		{"SRA", 0x28, func(cpu *CPU, val uint8) (uint8, bool) {
			c := val & 0x01
			return (val >> 1) | (val & 0x80), c != 0
		}},
		{"SLL", 0x30, func(cpu *CPU, val uint8) (uint8, bool) {
			c := val >> 7
			return (val << 1) | 1, c != 0
		}},
		{"SRL", 0x38, func(cpu *CPU, val uint8) (uint8, bool) {
			c := val & 0x01
			return val >> 1, c != 0
		}},
	}

	for _, op := range rotOps {
		for _, r := range regs8Idx {
			op := op
			r := r

			// DDCB
			mnemonicDD := op.name + " (IX+d)"
			if r.name != "" {
				mnemonicDD += ", " + r.name
			}
			RegisterInstruction(&DDCBTable, op.op|r.op, Instruction{
				Mnemonic: mnemonicDD,
				Length:   4,
				Cycles:   23,
				Execute: func(cpu *CPU) int {
					addr := uint16(int32(cpu.Regs.IX) + int32(cpu.LastDisplacement))
					val := cpu.Memory.Read(addr)
					res, newCarry := op.fn(cpu, val)
					cpu.Memory.Write(addr, res)
					if r.set != nil {
						r.set(cpu.Regs, res)
					}
					cpu.Regs.UpdateFlagsIOIn(res)
					cpu.Regs.SetFlag(FlagC, newCarry)
					return 23
				},
			})

			// FDCB
			mnemonicFD := op.name + " (IY+d)"
			if r.name != "" {
				mnemonicFD += ", " + r.name
			}
			RegisterInstruction(&FDCBTable, op.op|r.op, Instruction{
				Mnemonic: mnemonicFD,
				Length:   4,
				Cycles:   23,
				Execute: func(cpu *CPU) int {
					addr := uint16(int32(cpu.Regs.IY) + int32(cpu.LastDisplacement))
					val := cpu.Memory.Read(addr)
					res, newCarry := op.fn(cpu, val)
					cpu.Memory.Write(addr, res)
					if r.set != nil {
						r.set(cpu.Regs, res)
					}
					cpu.Regs.UpdateFlagsIOIn(res)
					cpu.Regs.SetFlag(FlagC, newCarry)
					return 23
				},
			})
		}
	}
}

// Helpers for CB instructions
func getVal8(cpu *CPU, r struct {
	name string
	get  func(r *Registers) uint8
	set  func(r *Registers, val uint8)
	op   uint8
}) (uint8, int) {
	if r.name == "(HL)" {
		return cpu.Memory.Read(cpu.Regs.HL()), 15
	}
	return r.get(cpu.Regs), 8
}

func setVal8(cpu *CPU, r struct {
	name string
	get  func(r *Registers) uint8
	set  func(r *Registers, val uint8)
	op   uint8
}, val uint8) {
	if r.name == "(HL)" {
		cpu.Memory.Write(cpu.Regs.HL(), val)
	} else {
		r.set(cpu.Regs, val)
	}
}
