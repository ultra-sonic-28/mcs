// Package z80 provides the implementation of the Zilog Z80 CPU.
package z80

import (
	"fmt"
	"log/slog"
	"strings"
)

// Handler is a function that executes an instruction on the CPU.
// It returns the number of T-cycles consumed by the instruction.
type Handler func(cpu *CPU) int

// Instruction represents a single Z80 opcode definition.
type Instruction struct {
	Mnemonic  string         // Symbolic representation (e.g., "LD A, (HL)")
	Length    int            // Instruction length in bytes
	Cycles    int            // Base T-states (cycles)
	AddrMode1 AddressingMode // First operand addressing mode
	AddrMode2 AddressingMode // Second operand addressing mode
	Execute   Handler        // Function to execute the instruction logic
}

// OpcodeTable represents a set of 256 instructions mapping to a single-byte opcode.
type OpcodeTable [256]Instruction

var (
	// MainTable contains the primary 8-bit opcodes.
	MainTable OpcodeTable

	// CBTable contains bitwise, shift, and rotate instructions (CB prefix).
	CBTable OpcodeTable

	// EDTable contains extended instructions, including Z80N extensions (ED prefix).
	EDTable OpcodeTable

	// DDTable contains IX-related instructions (DD prefix).
	DDTable OpcodeTable

	// FDTable contains IY-related instructions (FD prefix).
	FDTable OpcodeTable

	// DDCBTable contains IX bitwise instructions (DD CB prefix).
	DDCBTable OpcodeTable

	// FDCBTable contains IY bitwise instructions (FD CB prefix).
	FDCBTable OpcodeTable
)

func initTables() {
	for i := 0; i < 256; i++ {
		undefined := Instruction{
			Mnemonic:  fmt.Sprintf("UNK %02X", i),
			Length:    1,
			Cycles:    4,
			AddrMode1: AddrModeNone,
			AddrMode2: AddrModeNone,
			Execute: func(cpu *CPU) int {
				// Default behavior for undefined opcodes: act as NOP or log error.
				return 4
			},
		}
		MainTable[i] = undefined
		CBTable[i] = undefined
		EDTable[i] = undefined
		DDTable[i] = undefined
		FDTable[i] = undefined
		DDCBTable[i] = undefined
		FDCBTable[i] = undefined
	}
}

// PopulateIndexTables copies MainTable to DDTable and FDTable.
// This should be called after MainTable is fully populated but before
// DDTable/FDTable specific instructions are registered.
func PopulateIndexTables() {
	for i := 0; i < 256; i++ {
		// Only copy if the target instruction is still undefined
		if strings.HasPrefix(DDTable[i].Mnemonic, "UNK") {
			DDTable[i] = MainTable[i]
			if MainTable[i].Mnemonic != fmt.Sprintf("UNK %02X", i) {
				DDTable[i].Length++
				DDTable[i].Cycles += 4
			}
		}
		if strings.HasPrefix(FDTable[i].Mnemonic, "UNK") {
			FDTable[i] = MainTable[i]
			if MainTable[i].Mnemonic != fmt.Sprintf("UNK %02X", i) {
				FDTable[i].Length++
				FDTable[i].Cycles += 4
			}
		}
	}
}

// RegisterInstruction adds a new instruction definition to the specified table.
func RegisterInstruction(table *OpcodeTable, opcode uint8, instr Instruction) {
	table[opcode] = instr
}

// LogInstruction logs the details of a single instruction at DEBUG level.
func LogInstruction(prefix string, opcode uint8, instr Instruction) {
	slog.Debug("Instruction",
		"prefix", prefix,
		"opcode", fmt.Sprintf("0x%02X", opcode),
		"mnemonic", fmt.Sprintf("%q", instr.Mnemonic),
		"length", instr.Length,
		"cycles", instr.Cycles,
		"mode1", instr.AddrMode1.String(),
		"mode2", instr.AddrMode2.String(),
	)
}

// LogAllInstructions iterates through all opcode tables and logs every defined instruction at DEBUG level.
func LogAllInstructions() {
	slog.Debug("Logging all registered Z80 instructions")

	tables := []struct {
		name  string
		table *OpcodeTable
	}{
		{"Main", &MainTable},
		{"CB", &CBTable},
		{"DD", &DDTable},
		{"ED", &EDTable},
		{"FD", &FDTable},
		{"DDCB", &DDCBTable},
		{"FDCB", &FDCBTable},
	}

	for _, t := range tables {
		for i := 0; i < 256; i++ {
			instr := t.table[i]
			// Only log instructions that aren't "Undefined" (UNK)
			if len(instr.Mnemonic) < 3 || instr.Mnemonic[:3] != "UNK" {
				LogInstruction(t.name, uint8(i), instr)
			}
		}
	}
}

// CountInstructions returns the total number of unique non-undefined instructions registered.
func CountInstructions() int {
	count := 0
	tables := []*OpcodeTable{
		&MainTable,
		&CBTable,
		&EDTable,
		&DDTable,
		&FDTable,
		&DDCBTable,
		&FDCBTable,
	}

	for _, table := range tables {
		for i := 0; i < 256; i++ {
			if instr := table[i]; len(instr.Mnemonic) < 3 || instr.Mnemonic[:3] != "UNK" {
				count++
			}
		}
	}
	return count
}
