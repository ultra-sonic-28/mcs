// Package spectrum implements the ZX Spectrum 48K machine logic.
package spectrum

import (
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

// Reset performs a hardware reset of the machine.
func (m *Machine) Reset() {
	m.CPU.Reset()
	// Spectrum ROM starts with DI (0xF3), so PC=0 is correct.
}

// RunFrame executes instructions for one 50Hz frame.
func (m *Machine) RunFrame() {
	targetCycles := uint64(CyclesPerFrame)
	startCycles := m.CPU.Cycles

	// Trigger Interrupt at the start of the frame (ULA behavior)
	m.CPU.INT = true

	for (m.CPU.Cycles - startCycles) < targetCycles {
		cycles := m.CPU.Step()
		
		// Update Tape Signal
		m.Bus.TapeInState = m.Bus.Tape.Step(uint32(cycles))
	}

	// Toggle Flash every 16 frames (approx 0.32s)
	if m.CPU.Cycles % (69888 * 16) < 69888 {
		m.Bus.Display.FlashState = !m.Bus.Display.FlashState
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
