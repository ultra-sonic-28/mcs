// Package machine implements the ZX Spectrum 48K and 128K machine logic.
package machine

import (
	"fmt"
	"image/color"
	"log/slog"
	"mcs/internal/cpu/z80"
	"mcs/internal/machine/spectrum/bus"
	"mcs/internal/machine/spectrum/display"
	"mcs/internal/machine/spectrum/gui"
	"mcs/internal/machine/spectrum/keyboard"
	"mcs/internal/machine/spectrum/sound"
	"mcs/internal/machine/spectrum/tape"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	// StatusLineHeight is the height of the status line in pixels.
	StatusLineHeight = 12
	// ProcessorClock is 3.5 MHz (Spectrum 48K)
	ProcessorClock = 3500000
	// ProcessorClock128 is 3.5469 MHz (Spectrum 128K)
	ProcessorClock128 = 3546900
	// FramesPerSecond is 50 Hz
	FramesPerSecond = 50
	// CyclesPerFrame is the exact number of T-cycles in a 50Hz Spectrum 48K frame.
	CyclesPerFrame = 69888
	// CyclesPerFrame128 is the exact number of T-cycles in a 50Hz Spectrum 128K frame.
	CyclesPerFrame128 = 70908
)

// Bus defines the interface for the memory and I/O bus of a ZX Spectrum.
type Bus interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
	Read16(addr uint16) uint16
	In(port uint16) uint8
	Out(port uint16, val uint8)

	// GetDisplayMemory returns the memory currently being used for display.
	GetDisplayMemory() []byte

	// Common components access
	GetKeyboard() *keyboard.Keyboard
	GetTape() *tape.Tape
	GetDisplay() *display.Display
	GetAY() *sound.AY38912
	GetBorderColor() uint8
	GetTapeInState() bool
	SetTapeInState(state bool)
	GetBeeperState() bool

	// IsRom1Active returns true if the 48K BASIC ROM is currently paged in.
	IsRom1Active() bool
}

// BaseMachine represents common logic for all ZX Spectrum models.
type BaseMachine struct {
	CPU *z80.CPU
	Bus Bus

	// Metadata
	MachineName string

	// Graphics
	screenImage *ebiten.Image

	// AutoStart state
	autoStartEnabled bool
	autoStartStep    int
	autoStartTimer   int
	autoStartTyping  bool

	// Timing
	ClockRate      uint64
	CyclesPerFrame uint64

	// Audio
	audioContext *audio.Context
	audioPlayer  *audio.Player
}

// Machine48 represents the ZX Spectrum 48K emulator.
type Machine48 struct {
	BaseMachine
}

// Machine128 represents the ZX Spectrum 128K emulator.
type Machine128 struct {
	BaseMachine
}

// NewMachine creates and initializes a new Spectrum 48K machine.
func NewMachine() *Machine48 {
	slog.Info("Creating Spectrum 48K Machine")
	b := bus.NewBus48()
	cpu := z80.NewCPU(b, b)

	return &Machine48{
		BaseMachine: BaseMachine{
			CPU:            cpu,
			Bus:            b,
			MachineName:    "Spectrum 48K",
			ClockRate:      ProcessorClock,
			CyclesPerFrame: CyclesPerFrame,
		},
	}
}

// NewMachine128 creates and initializes a new Spectrum 128K machine.
func NewMachine128() *Machine128 {
	slog.Info("Creating Spectrum 128K Machine")
	b := bus.NewBus128()
	cpu := z80.NewCPU(b, b)

	return &Machine128{
		BaseMachine: BaseMachine{
			CPU:            cpu,
			Bus:            b,
			MachineName:    "Spectrum 128K",
			ClockRate:      ProcessorClock128,
			CyclesPerFrame: CyclesPerFrame128,
		},
	}
}

// EnableAutoStart prepares the machine to automatically load and run the tape.
func (m *BaseMachine) EnableAutoStart() {
	m.autoStartEnabled = true
	m.autoStartStep = 0
	m.autoStartTimer = 150 // Wait 150 frames (3 seconds) for boot
	slog.Info("Auto-start mechanism enabled")
}

// Reset performs a hardware reset of the machine.
func (m *BaseMachine) Reset() {
	//m.CPU.Reset()
	// Spectrum ROM starts with DI (0xF3), so PC=0 is correct.
}

// UpdateKeyboard reads the host keyboard state and updates the Spectrum keyboard.
func (m *BaseMachine) UpdateKeyboard() {
	if m.autoStartTyping {
		// During auto-start typing, we don't want physical keyboard to interfere.
		return
	}
	k := m.Bus.GetKeyboard()
	for eKey, sKey := range gui.KeyMap {
		k.SetKeyState(sKey, ebiten.IsKeyPressed(eKey))
	}
}

// RunFrame executes instructions for one 50Hz frame.
func (m *BaseMachine) RunFrame() {
	targetCycles := m.CyclesPerFrame
	startCycles := m.CPU.Cycles

	// Update AutoStart
	m.updateAutoStart()

	// Trigger Interrupt at the start of the frame (ULA behavior)
	m.CPU.INT = true

	ay := m.Bus.GetAY()
	cyclesToNextSample := float64(m.ClockRate) / float64(sound.SampleRate)
	accumulatedCycles := 0.0

	for (m.CPU.Cycles - startCycles) < targetCycles {
		// Instant Load Trap: Intercept ROM Tape Loading Routine (LD-BYTES at 0x0556)
		// Only trap if we are in the 48K BASIC ROM (ROM 1 on 128K).
		if m.CPU.Regs.PC == 0x0556 && m.Bus.IsRom1Active() {
			t := m.Bus.GetTape()
			if len(t.Blocks) > 0 {
				m.instantLoadBlock()
			} else {
				slog.Debug("LD-BYTES called but no tape blocks loaded")
			}
		}

		cycles := m.CPU.Step()

		// Update Sound
		if ay != nil {
			beeper := m.Bus.GetBeeperState()
			for i := 0; i < cycles; i++ {
				// AY clock is exactly half of CPU clock (1.77MHz)
				if (startCycles+(m.CPU.Cycles-startCycles-uint64(cycles)+uint64(i)))%2 == 0 {
					ay.Tick(beeper)
				}

				accumulatedCycles += 1.0
				if accumulatedCycles >= cyclesToNextSample {
					accumulatedCycles -= cyclesToNextSample
					l, r := ay.GetSample()
					ay.AddAudioSample(l, r)
				}
			}
		}

		// Update Tape Signal
		m.Bus.SetTapeInState(m.Bus.GetTape().Step(uint32(cycles)))
	}

	// Toggle Flash every 16 frames (approx 0.32s)
	if m.CPU.Cycles%(m.CyclesPerFrame*16) < m.CyclesPerFrame {
		m.Bus.GetDisplay().FlashState = !m.Bus.GetDisplay().FlashState
	}
}

func (m *BaseMachine) updateAutoStart() {
	if !m.autoStartEnabled {
		return
	}

	if m.autoStartTimer > 0 {
		m.autoStartTimer--
		return
	}

	k := m.Bus.GetKeyboard()

	// Keyboard sequence for LOAD "" : RUN <ENTER>
	// 48K mode keywords: J -> LOAD, Symbol Shift + P -> ", Symbol Shift + Z -> :, R -> RUN
	switch m.autoStartStep {
	case 0: // Press J (LOAD)
		slog.Debug("Auto-start: Pressing J")
		m.autoStartTyping = true
		k.SetKeyState(keyboard.KeyJ, true)
		m.autoStartTimer = 10
		m.autoStartStep++
	case 1: // Release J
		slog.Debug("Auto-start: Releasing J")
		k.SetKeyState(keyboard.KeyJ, false)
		m.autoStartTimer = 10
		m.autoStartStep++
	case 2: // Press Symbol Shift
		slog.Debug("Auto-start: Pressing Symbol Shift")
		k.SetKeyState(keyboard.KeySymbolShift, true)
		m.autoStartTimer = 5
		m.autoStartStep++
	case 3: // Press P (")
		slog.Debug("Auto-start: Pressing P (first quote)")
		k.SetKeyState(keyboard.KeyP, true)
		m.autoStartTimer = 10
		m.autoStartStep++
	case 4: // Release P
		slog.Debug("Auto-start: Releasing P")
		k.SetKeyState(keyboard.KeyP, false)
		m.autoStartTimer = 10
		m.autoStartStep++
	case 5: // Press P again (")
		slog.Debug("Auto-start: Pressing P (second quote)")
		k.SetKeyState(keyboard.KeyP, true)
		m.autoStartTimer = 10
		m.autoStartStep++
	case 6: // Release P
		slog.Debug("Auto-start: Releasing P")
		k.SetKeyState(keyboard.KeyP, false)
		m.autoStartTimer = 10
		m.autoStartStep++
	case 7: // Press Z (:)
		slog.Debug("Auto-start: Pressing Z (colon)")
		k.SetKeyState(keyboard.KeyZ, true)
		m.autoStartTimer = 10
		m.autoStartStep++
	case 8: // Release Z and Symbol Shift
		slog.Debug("Auto-start: Releasing Z and Symbol Shift")
		k.SetKeyState(keyboard.KeyZ, false)
		k.SetKeyState(keyboard.KeySymbolShift, false)
		m.autoStartTimer = 10
		m.autoStartStep++
	case 9: // Press R (RUN)
		slog.Debug("Auto-start: Pressing R (RUN)")
		k.SetKeyState(keyboard.KeyR, true)
		m.autoStartTimer = 10
		m.autoStartStep++
	case 10: // Release R
		slog.Debug("Auto-start: Releasing R")
		k.SetKeyState(keyboard.KeyR, false)
		m.autoStartTimer = 10
		m.autoStartStep++
	case 11: // Press Enter
		slog.Debug("Auto-start: Pressing Enter")
		k.SetKeyState(keyboard.KeyEnter, true)
		m.autoStartTimer = 10
		m.autoStartStep++
	case 12: // Release Enter
		slog.Debug("Auto-start: Releasing Enter")
		k.SetKeyState(keyboard.KeyEnter, false)
		m.autoStartTimer = 100 // Wait 2s
		m.autoStartStep++
	case 13: // Finished typing
		slog.Info("Auto-typing complete")
		m.autoStartTyping = false
		m.autoStartEnabled = false
	}
}

func (m *BaseMachine) instantLoadBlock() {
	t := m.Bus.GetTape()
	if t.CurrentBlock >= len(t.Blocks) {
		slog.Debug("LD-BYTES called but all tape blocks already loaded")
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
	if block[0] != expectedFlag {
		slog.Warn("Instant load flag mismatch, skipping block",
			"block", t.CurrentBlock+1,
			"block_flag", fmt.Sprintf("0x%02X", block[0]),
			"expected_flag", fmt.Sprintf("0x%02X", expectedFlag))
		return
	}

	slog.Debug("Instant loading tape block",
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
	m.CPU.Regs.SetBC(0x0001)           // Not strictly necessary but common
	m.CPU.Regs.L = block[len(block)-1] // Checksum
	m.CPU.Regs.H = m.CPU.Regs.L
	m.CPU.Regs.SetFlag(z80.FlagC, true) // Success
	m.CPU.Regs.SetFlag(z80.FlagZ, true)

	// Perform a RET (return to caller)
	retAddr := m.Bus.Read16(m.CPU.Regs.SP)
	slog.Debug("Instant load complete, returning to ROM",
		"addr", fmt.Sprintf("0x%04X", retAddr),
		"sp", fmt.Sprintf("0x%04X", m.CPU.Regs.SP),
		"next_block", t.CurrentBlock+1,
		"total_blocks", len(t.Blocks))

	m.CPU.Regs.SP += 2
	m.CPU.Regs.PC = retAddr

	// Advance to next block
	t.CurrentBlock++
	if t.CurrentBlock >= len(t.Blocks) {
		t.Active = false
		slog.Info("Tape loading finished (Instant)")
	}
}

func (m *BaseMachine) initAudio() {
	if m.audioContext == nil {
		m.audioContext = audio.NewContext(sound.SampleRate)
	}
	ay := m.Bus.GetAY()
	if ay != nil && m.audioPlayer == nil {
		var err error
		m.audioPlayer, err = m.audioContext.NewPlayer(ay)
		if err != nil {
			slog.Error("Failed to create audio player", "error", err)
			return
		}
		slog.Info("AudioPlayer created")
		// Play in a goroutine to avoid hanging the main thread on some platforms/drivers
		go func() {
			m.audioPlayer.Play()
			slog.Debug("AudioPlayer playback started")
		}()
	}
}

// Run executes the machine using Ebitengine.
func (m *BaseMachine) Run() error {
	m.initAudio()
	slog.Info("Setting Ebitengine UI")
	ebiten.SetWindowSize(display.ScreenWidth*2, (display.ScreenHeight+StatusLineHeight)*2)
	ebiten.SetWindowTitle("MCS - Multi CPUs System")
	ebiten.SetTPS(FramesPerSecond) // Set to 50 TPS for Spectrum
	slog.Info("Starting Spectrum Ebitengine loop")
	return ebiten.RunGame(m)
}

// Update handles logical state changes.
func (m *BaseMachine) Update() error {
	m.UpdateKeyboard()
	m.RunFrame()
	return nil
}

// Draw handles rendering.
func (m *BaseMachine) Draw(screen *ebiten.Image) {
	m.Bus.GetDisplay().RenderFrame(m.Bus.GetDisplayMemory())

	// Draw Spectrum Screen
	if m.screenImage == nil {
		m.screenImage = ebiten.NewImage(display.ScreenWidth, display.ScreenHeight)
	}
	m.screenImage.WritePixels(m.Bus.GetDisplay().FrameBuffer[:])
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(m.screenImage, op)

	// Draw Status Line Background (Dark grey)
	statusRect := ebiten.NewImage(display.ScreenWidth, StatusLineHeight)
	statusRect.Fill(color.RGBA{32, 32, 32, 255})
	opRect := &ebiten.DrawImageOptions{}
	opRect.GeoM.Translate(0, float64(display.ScreenHeight))
	screen.DrawImage(statusRect, opRect)

	// Draw Status Line Sections
	t := m.Bus.GetTape()
	tapeName := "No tape"
	if t.Filename != "" {
		tapeName = filepath.Base(t.Filename)
	}

	// Colors
	textColor := color.RGBA{200, 200, 200, 255}
	sepColor := color.RGBA{100, 100, 100, 255}

	// 1. Tape Section (50% = 128px)
	// Each char is 6px wide (5px font + 1px space).
	// 128 / 6 = 21 chars.
	const maxTapeChars = 20
	displayTapeName := tapeName
	if len(displayTapeName) > maxTapeChars {
		displayTapeName = displayTapeName[:maxTapeChars-3] + "..."
	}
	gui.DrawSmallText(screen, "|", 0, display.ScreenHeight+2, sepColor)
	gui.DrawSmallText(screen, displayTapeName, 6, display.ScreenHeight+2, textColor)

	// 2. CPU Section (15% = 38px) -> starts at 128
	// "Z80" is 18px wide. (38 - 18) / 2 = 10px offset.
	gui.DrawSmallText(screen, "|", 128, display.ScreenHeight+2, sepColor)
	gui.DrawSmallText(screen, "Z80", 128+10, display.ScreenHeight+2, textColor)

	// 3. Machine Section (35% = 90px) -> starts at 166 (128 + 38)
	gui.DrawSmallText(screen, "|", 166, display.ScreenHeight+2, sepColor)
	gui.DrawSmallText(screen, m.MachineName, 166+6, display.ScreenHeight+2, textColor)
	gui.DrawSmallText(screen, "|", 255, display.ScreenHeight+2, sepColor)
}

// Layout defines the screen dimensions.
func (m *BaseMachine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return display.ScreenWidth, display.ScreenHeight + StatusLineHeight
}
