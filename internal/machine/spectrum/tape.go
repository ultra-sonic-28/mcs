// Package spectrum implements the ZX Spectrum 48K machine logic.
package spectrum

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// TapeState represents the current part of the tape block being played.
type TapeState int

const (
	TapeIdle TapeState = iota
	TapePilot
	TapeSync1
	TapeSync2
	TapeData
	TapePause
)

// Tape handles the loading and playback of .tap cassette files.
type Tape struct {
	Blocks [][]byte
	
	// Playback state
	Active        bool
	CurrentBlock  int
	State         TapeState
	PulseCount    int
	ByteIndex     int
	BitIndex      int
	PulseLength   uint32
	CyclesInPulse uint32
	Signal        bool
}

// NewTape creates a new empty Tape instance.
func NewTape() *Tape {
	return &Tape{}
}

// LoadTAP reads a .tap file and extracts its blocks.
func (t *Tape) LoadTAP(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	t.Blocks = nil
	for i := 0; i < len(data); {
		if i+2 > len(data) { break }
		length := uint16(data[i]) | uint16(data[i+1])<<8
		i += 2
		if i+int(length) > len(data) { break }
		
		block := make([]byte, length)
		copy(block, data[i:i+int(length)])
		t.Blocks = append(t.Blocks, block)
		i += int(length)
	}

	slog.Debug("Loaded .tap file", "file", filename, "blocks_count", len(t.Blocks))
	for i, block := range t.Blocks {
		detail := t.getBlockDetail(block)
		slog.Debug("Tape block info", "block", i+1, "detail", detail, "length", len(block))
	}
	return nil
}

func (t *Tape) getBlockDetail(block []byte) string {
	if len(block) < 1 {
		return "empty"
	}

	flag := block[0]
	if flag == 0x00 && len(block) >= 18 {
		// Standard Spectrum ROM Header (19 bytes in .tap: 1 flag + 17 header + 1 checksum)
		typ := block[1]
		name := strings.TrimRight(string(block[2:12]), " ")
		size := uint16(block[12]) | uint16(block[13])<<8
		param1 := uint16(block[14]) | uint16(block[15])<<8
		
		switch typ {
		case 0:
			lineInfo := ""
			if param1 <= 32767 {
				lineInfo = fmt.Sprintf(" (line %d)", param1)
			}
			return fmt.Sprintf("Program \"%s\"%s (size %d)", name, lineInfo, size)
		case 1:
			return fmt.Sprintf("Number array \"%s\" (size %d)", name, size)
		case 2:
			return fmt.Sprintf("Character array \"%s\" (size %d)", name, size)
		case 3:
			return fmt.Sprintf("Code \"%s\" (addr %d) (size %d)", name, param1, size)
		default:
			return "Unknown header"
		}
	} else if flag == 0xFF {
		return "data"
	}

	return "unknown"
}

// Play starts the tape playback.
func (t *Tape) Play() {
	if len(t.Blocks) == 0 { return }
	t.Active = true
	t.CurrentBlock = 0
	t.startBlock()
}

// Stop stops the tape playback.
func (t *Tape) Stop() {
	t.Active = false
}

func (t *Tape) startBlock() {
	t.State = TapePilot
	// Headers (Flag < 128) have 8063 pulses, Data blocks have 3223 pulses.
	if t.Blocks[t.CurrentBlock][0] < 128 {
		t.PulseCount = 8063
	} else {
		t.PulseCount = 3223
	}
	t.PulseLength = 2168
	t.CyclesInPulse = 0
	t.Signal = true
}

// Step advances the tape state by the given number of T-cycles.
// It returns the current signal state (true = EAR high).
func (t *Tape) Step(cycles uint32) bool {
	if !t.Active { return true }

	t.CyclesInPulse += cycles
	if t.CyclesInPulse >= t.PulseLength {
		t.CyclesInPulse -= t.PulseLength
		t.Signal = !t.Signal
		t.PulseCount--

		if t.PulseCount <= 0 {
			t.transition()
		}
	}
	return t.Signal
}

func (t *Tape) transition() {
	switch t.State {
	case TapePilot:
		t.State = TapeSync1
		t.PulseCount = 1
		t.PulseLength = 667
	case TapeSync1:
		t.State = TapeSync2
		t.PulseCount = 1
		t.PulseLength = 735
	case TapeSync2:
		t.State = TapeData
		t.ByteIndex = 0
		t.BitIndex = 7
		t.nextBit()
	case TapeData:
		t.BitIndex--
		if t.BitIndex < 0 {
			t.BitIndex = 7
			t.ByteIndex++
		}
		if t.ByteIndex >= len(t.Blocks[t.CurrentBlock]) {
			t.State = TapePause
			t.PulseCount = 1
			t.PulseLength = 3500000 // 1 second pause
		} else {
			t.nextBit()
		}
	case TapePause:
		t.CurrentBlock++
		if t.CurrentBlock < len(t.Blocks) {
			t.startBlock()
		} else {
			t.Active = false
		}
	}
}

func (t *Tape) nextBit() {
	val := t.Blocks[t.CurrentBlock][t.ByteIndex]
	bit := (val >> uint(t.BitIndex)) & 0x01
	t.PulseCount = 2
	if bit == 0 {
		t.PulseLength = 855
	} else {
		t.PulseLength = 1710
	}
}
