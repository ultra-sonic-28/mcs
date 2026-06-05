// Package spectrum implements the ZX Spectrum 48K machine logic.
package spectrum

import (
	"fmt"
	"log/slog"
	"mcs/internal/cpu/z80"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// ProcessorClock is 3.5 MHz
	ProcessorClock = 3500000
	// FramesPerSecond is 50 Hz
	FramesPerSecond = 50
	// CyclesPerFrame is the exact number of T-cycles in a 50Hz Spectrum 48K frame.
	CyclesPerFrame = 69888
)

// Machine represents the complete ZX Spectrum 48K emulator.
type Machine struct {
	CPU *z80.CPU
	Bus *Bus

	// AutoStart state
	autoStartEnabled bool
	autoStartStep    int
	autoStartTimer   int
}

// NewMachine creates and initializes a new Spectrum 48K machine.
func NewMachine() *Machine {
	slog.Info("Creating Spectrum 48K Machine")
	bus := NewBus()
	cpu := z80.NewCPU(bus, bus)

	return &Machine{
		CPU: cpu,
		Bus: bus,
	}
}

// EnableAutoStart prepares the machine to automatically load and run the tape.
func (m *Machine) EnableAutoStart() {
	m.autoStartEnabled = true
	m.autoStartStep = 0
	m.autoStartTimer = 100 // Wait 100 frames (2 seconds) for boot
	slog.Info("Auto-start mechanism enabled")
}

// Reset performs a hardware reset of the machine.
func (m *Machine) Reset() {
	//m.CPU.Reset()
	// Spectrum ROM starts with DI (0xF3), so PC=0 is correct.
}

// RunFrame executes instructions for one 50Hz frame.
func (m *Machine) RunFrame() {
	targetCycles := uint64(CyclesPerFrame)
	startCycles := m.CPU.Cycles

	// Update AutoStart
	m.updateAutoStart()

	// Trigger Interrupt at the start of the frame (ULA behavior)
	m.CPU.INT = true

	for (m.CPU.Cycles - startCycles) < targetCycles {
		// Instant Load Trap: Intercept ROM Tape Loading Routine (LD-BYTES at 0x0556)
		if m.CPU.Regs.PC == 0x0556 && len(m.Bus.Tape.Blocks) > 0 {
			m.instantLoadBlock()
		}

		cycles := m.CPU.Step()

		// Update Tape Signal
		m.Bus.TapeInState = m.Bus.Tape.Step(uint32(cycles))
	}

	// Toggle Flash every 16 frames (approx 0.32s)
	if m.CPU.Cycles%(69888*16) < 69888 {
		m.Bus.Display.FlashState = !m.Bus.Display.FlashState
	}
}

func (m *Machine) updateAutoStart() {
	if !m.autoStartEnabled {
		return
	}

	if m.autoStartTimer > 0 {
		m.autoStartTimer--
		return
	}

	// Keyboard sequence for LOAD "" <ENTER>
	// 48K mode keywords: J -> LOAD, Symbol Shift + P -> "
	switch m.autoStartStep {
	case 0: // Press J
		m.Bus.Keyboard.SetKeyState(KeyJ, true)
		m.autoStartTimer = 5
		m.autoStartStep++
	case 1: // Release J
		m.Bus.Keyboard.SetKeyState(KeyJ, false)
		m.autoStartTimer = 5
		m.autoStartStep++
	case 2: // Press Symbol Shift + P (")
		m.Bus.Keyboard.SetKeyState(KeySymbolShift, true)
		m.Bus.Keyboard.SetKeyState(KeyP, true)
		m.autoStartTimer = 5
		m.autoStartStep++
	case 3: // Release Symbol Shift + P
		m.Bus.Keyboard.SetKeyState(KeySymbolShift, false)
		m.Bus.Keyboard.SetKeyState(KeyP, false)
		m.autoStartTimer = 5
		m.autoStartStep++
	case 4: // Press Symbol Shift + P (")
		m.Bus.Keyboard.SetKeyState(KeySymbolShift, true)
		m.Bus.Keyboard.SetKeyState(KeyP, true)
		m.autoStartTimer = 5
		m.autoStartStep++
	case 5: // Release Symbol Shift + P
		m.Bus.Keyboard.SetKeyState(KeySymbolShift, false)
		m.Bus.Keyboard.SetKeyState(KeyP, false)
		m.autoStartTimer = 5
		m.autoStartStep++
	case 6: // Press Enter
		m.Bus.Keyboard.SetKeyState(KeyEnter, true)
		m.autoStartTimer = 5
		m.autoStartStep++
	case 7: // Release Enter
		m.Bus.Keyboard.SetKeyState(KeyEnter, false)
		m.autoStartTimer = 20
		m.autoStartStep++
	case 8: // Ready for trap
		slog.Info("Auto-typing complete, waiting for ROM to start loading")
		// We don't call m.Bus.Tape.Play() here to avoid advancing the tape via audio cycles.
		// The trap will handle it.
		m.autoStartEnabled = false
	}
}

func (m *Machine) instantLoadBlock() {
	t := m.Bus.Tape
	if t.CurrentBlock >= len(t.Blocks) {
		return
	}

	block := t.Blocks[t.CurrentBlock]
	
	// A = expected flag (0x00 for header, 0xFF for data)
	// IX = destination address
	// DE = expected length
	expectedFlag := m.CPU.Regs.A
	destAddr := m.CPU.Regs.IX
	expectedLen := m.CPU.Regs.DE()

	// If the ROM is looking for a header (A=0), but we are at a data block, we skip.
	// Actually, standard LOAD routine expects blocks in order.
	// If the flags don't match, it's usually because the ROM is searching for a header
	// while we are at the data block of a PREVIOUS program, or vice-versa.
	if block[0] != expectedFlag {
		slog.Debug("Instant load flag mismatch, skipping block", 
			"block", t.CurrentBlock+1, 
			"block_flag", fmt.Sprintf("0x%02X", block[0]), 
			"expected_flag", fmt.Sprintf("0x%02X", expectedFlag))
		return
	}

	slog.Info("Instant loading tape block", 
		"block", t.CurrentBlock+1, 
		"type", block[0],
		"dest", fmt.Sprintf("0x%04X", destAddr),
		"len", expectedLen)

	// Copy data (skipping the flag byte)
	// Block contains: [Flag] [Data...] [Checksum]
	dataLen := uint16(len(block)) - 2
	if expectedLen < dataLen {
		dataLen = expectedLen
	}

	for i := uint16(0); i < dataLen; i++ {
		m.Bus.Write(destAddr+i, block[i+1])
	}

	// Update Registers as if LD-BYTES finished successfully
	m.CPU.Regs.IX += dataLen
	m.CPU.Regs.SetDE(0)
	m.CPU.Regs.SetBC(0x0001) // Not strictly necessary but common
	m.CPU.Regs.L = block[len(block)-1] // Checksum
	m.CPU.Regs.H = m.CPU.Regs.L
	m.CPU.Regs.SetFlag(z80.FlagC, true) // Success
	m.CPU.Regs.SetFlag(z80.FlagZ, true)

	// Perform a RET (return to caller)
	retAddr := m.Bus.Read16(m.CPU.Regs.SP)
	m.CPU.Regs.SP += 2
	m.CPU.Regs.PC = retAddr

	slog.Debug("Instant load complete, returning to", "addr", fmt.Sprintf("0x%04X", retAddr))

	// Advance to next block
	t.CurrentBlock++
	if t.CurrentBlock >= len(t.Blocks) {
		t.Active = false
		slog.Info("Tape loading finished (Instant)")
	}
}

// Run executes the machine using Ebitengine.
// This is a blocking call.
func (m *Machine) Run() error {
	slog.Info("Starting Spectrum 48K Ebitengine loop")
	ebiten.SetWindowSize(ScreenWidth*2, ScreenHeight*2)
	ebiten.SetWindowTitle("MCS - ZX Spectrum 48K")
	ebiten.SetTPS(FramesPerSecond) // Set to 50 TPS for Spectrum
	return ebiten.RunGame(m)
}
