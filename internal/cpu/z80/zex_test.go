//go:build zex
// +build zex

package z80

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ZexTraceEntry stores the state of the CPU at a specific execution step.
type ZexTraceEntry struct {
	PC       uint16
	Mnemonic string
	A        uint8
	BC       uint16
	DE       uint16
	HL       uint16
	IX       uint16
	IY       uint16
	SP       uint16
	F        uint8
}

// TestZex runs the Z80 Instruction Set Exerciser binaries (ZEXALL/ZEXDOC).
// These binaries must be present in the project root or a specified directory.
func TestZex(t *testing.T) {
	tests := []string{"zexdoc.com", "zexall.com"}
	const historySize = 50

	for _, filename := range tests {
		t.Run(filename, func(t *testing.T) {
			path := filepath.Join("..", "..", "..", filename)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Skipf("%s not found, skipping. Please place it in the project root.", filename)
			}

			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read %s: %v", filename, err)
			}

			bus := &ZexBus{}
			cpu := NewCPU(bus, bus)

			// Load program at 0x100 (standard CP/M)
			for i, b := range data {
				bus.Write(0x100+uint16(i), b)
			}

			// BDOS trap at 0x0005
			// RET instruction (0xC9) at 0x0005 to return to the program after call
			bus.Write(0x0005, 0xC9)

			cpu.Regs.PC = 0x100
			cpu.Regs.SP = 0xF000 // Safe stack area

			fmt.Printf("\n--- Running %s ---\n", filename)

			history := make([]ZexTraceEntry, historySize)
			historyIdx := 0
			outputBuffer := ""
			lastUnkPC := uint16(0)

			// Run until JP 0 (standard CP/M exit)
			for {
				if cpu.Regs.PC == 0x0005 {
					// BDOS Trap
					funcID := cpu.Regs.C
					if funcID == 2 {
						// Output character in E
						char := byte(cpu.Regs.E)
						fmt.Printf("%c", char)
						outputBuffer += string(char)
					} else if funcID == 9 {
						// Output string at DE until '$'
						addr := cpu.Regs.DE()
						str := ""
						for {
							char := bus.Read(addr)
							if char == '$' {
								break
							}
							str += string(char)
							addr++
						}
						fmt.Print(str)
						outputBuffer += str
					}

					// Check for ERROR in the last added output
					if strings.Contains(outputBuffer, "ERROR") {
						// Continue running for a bit to capture more output
						for j := 0; j < 100; j++ {
							if cpu.Regs.PC == 0x0005 {
								funcID := cpu.Regs.C
								if funcID == 2 {
									char := byte(cpu.Regs.E)
									fmt.Printf("%c", char)
									outputBuffer += string(char)
								} else if funcID == 9 {
									addr := cpu.Regs.DE()
									for {
										char := bus.Read(addr)
										if char == '$' {
											break
										}
										fmt.Printf("%c", char)
										outputBuffer += string(char)
										addr++
									}
								}
								cpu.Regs.PC = uint16(bus.Read(cpu.Regs.SP)) | (uint16(bus.Read(cpu.Regs.SP+1)) << 8)
								cpu.Regs.SP += 2
							}
							cpu.Step()
						}
						fmt.Printf("\n\nFAILURE DETECTED in %s. Dumping last %d instructions:\n", filename, historySize)
						dumpHistory(history, historyIdx)
						t.Errorf("%s reported an ERROR", filename)
						return
					}
					// Keep buffer size reasonable
					if len(outputBuffer) > 1000 {
						outputBuffer = outputBuffer[len(outputBuffer)-500:]
					}
				}

				if cpu.Regs.PC == 0x0000 {
					fmt.Println("\nProgram exited via JP 0")
					break
				}

				// Record state before execution
				pc := cpu.Regs.PC
				opcode := bus.Read(pc)
				var instr Instruction
				isUnk := false
				if opcode == 0xDD {
					next := bus.Read(pc + 1)
					if next == 0xCB {
						instr = DDCBTable[bus.Read(pc+3)]
					} else {
						instr = DDTable[next]
					}
				} else if opcode == 0xFD {
					next := bus.Read(pc + 1)
					if next == 0xCB {
						instr = FDCBTable[bus.Read(pc+3)]
					} else {
						instr = FDTable[next]
					}
				} else if opcode == 0xED {
					instr = EDTable[bus.Read(pc+1)]
				} else if opcode == 0xCB {
					instr = CBTable[bus.Read(pc+1)]
				} else {
					instr = MainTable[opcode]
				}

				if strings.HasPrefix(instr.Mnemonic, "UNK") {
					isUnk = true
				}

				history[historyIdx] = ZexTraceEntry{
					PC:       pc,
					Mnemonic: instr.Mnemonic,
					A:        cpu.Regs.A,
					BC:       cpu.Regs.BC(),
					DE:       cpu.Regs.DE(),
					HL:       cpu.Regs.HL(),
					IX:       cpu.Regs.IX,
					IY:       cpu.Regs.IY,
					SP:       cpu.Regs.SP,
					F:        cpu.Regs.F,
				}
				historyIdx = (historyIdx + 1) % historySize

				if isUnk && pc != lastUnkPC {
					fmt.Printf("\nWARNING: Unimplemented instruction at %04X: %02X (Mnemonic: %s)\n", pc, opcode, instr.Mnemonic)
					lastUnkPC = pc
				}

				cpu.Step()

				//if cpu.Cycles%100000000 == 0 { // Every 100M cycles
				//	fmt.Printf("(%dM) ", cpu.Cycles/1000000)
				//}
			}
		})
	}
}

func dumpHistory(history []ZexTraceEntry, currentIdx int) {
	size := len(history)
	fmt.Println("PC     Mnemonic           A  BC   DE   HL   IX   IY   SP   F")
	fmt.Println("----------------------------------------------------------------")
	for i := 0; i < size; i++ {
		entry := history[(currentIdx+i)%size]
		if entry.PC == 0 {
			continue
		}
		fmt.Printf("%04X   %-18s %02X %04X %04X %04X %04X %04X %04X %02X\n",
			entry.PC, entry.Mnemonic, entry.A, entry.BC, entry.DE, entry.HL, entry.IX, entry.IY, entry.SP, entry.F)
	}
}

type ZexBus struct {
	Memory [65536]uint8
	IO     [65536]uint8
	Last   uint16
}

func (b *ZexBus) Read(addr uint16) uint8       { return b.Memory[addr] }
func (b *ZexBus) Write(addr uint16, val uint8) { b.Memory[addr] = val }
func (b *ZexBus) In(port uint16) uint8 {
	b.Last = port
	return b.IO[port]
}
func (b *ZexBus) Out(port uint16, val uint8) {
	b.Last = port
	b.IO[port] = val
}
