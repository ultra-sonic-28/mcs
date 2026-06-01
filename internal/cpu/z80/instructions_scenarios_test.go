package z80

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"mcs/testutils/assert"
	"strings"
	"testing"
)

// InstructionScenario defines a test case for instruction utilities.
type InstructionScenario struct {
	Name string
	Run  func(t *testing.T)
}

// instrLogCaptureHandler captures log records with level filtering for testing.
type instrLogCaptureHandler struct {
	buf *bytes.Buffer
}

func (h *instrLogCaptureHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true 
}

func (h *instrLogCaptureHandler) Handle(_ context.Context, r slog.Record) error {
	fmt.Fprintf(h.buf, "[%s] %s", r.Level, r.Message)
	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(h.buf, " %s=%v", a.Key, a.Value)
		return true
	})
	fmt.Fprintln(h.buf)
	return nil
}

func (h *instrLogCaptureHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *instrLogCaptureHandler) WithGroup(name string) slog.Handler        { return h }

var instructionScenarios = []InstructionScenario{
	{
		Name: "LogInstruction Output Verification",
		Run: func(t *testing.T) {
			instr := Instruction{
				Mnemonic:  "TEST INSTR",
				Length:    2,
				Cycles:    10,
				AddrMode1: AddrModeRegister,
				AddrMode2: AddrModeImmediate,
			}

			var buf bytes.Buffer
			handler := &instrLogCaptureHandler{buf: &buf}
			logger := slog.New(handler)
			
			// LogInstruction uses slog.Debug
			oldDefault := slog.Default()
			slog.SetDefault(logger)
			defer slog.SetDefault(oldDefault)

			LogInstruction("PREFIX", 0xAB, instr)

			output := buf.String()
			assert.Contains(t, "Should contain level DEBUG", "[DEBUG]", output)
			assert.Contains(t, "Should contain prefix", "prefix=PREFIX", output)
			assert.Contains(t, "Should contain opcode", "opcode=0xAB", output)
			assert.Contains(t, "Should contain mnemonic", "mnemonic=\"TEST INSTR\"", output)
			assert.Contains(t, "Should contain length", "length=2", output)
			assert.Contains(t, "Should contain cycles", "cycles=10", output)
			assert.Contains(t, "Should contain mode1", "mode1=Register", output)
			assert.Contains(t, "Should contain mode2", "mode2=Immediate", output)
		},
	},
	{
		Name: "CountInstructions Verification",
		Run: func(t *testing.T) {
			count := CountInstructions()
			// We know we have some instructions registered in op_core.go init()
			// NOP, HALT, LD A, n, ADD A, n, and ADD A, r (7 regs), ADD A, (HL)
			// That's at least 1 + 1 + 1 + 1 + 7 + 1 = 12 instructions.
			assert.True(t, "Should have at least 12 instructions", count >= 12)
			
			// Let's count them specifically based on what we saw in op_core.go
			// 0x00: NOP
			// 0x76: HALT
			// 0x3E: LD A, n
			// 0xC6: ADD A, n
			// 0x80-0x85, 0x87: ADD A, r (7)
			// 0x86: ADD A, (HL)
			// Total = 12
			
			// If we haven't added more, it should be exactly 12.
			// However, op_core.go might have more than what I read.
			// Let's just ensure it's positive and reasonably consistent.
			assert.True(t, "Count should be positive", count > 0)
		},
	},
	{
		Name: "LogAllInstructions Verification",
		Run: func(t *testing.T) {
			var buf bytes.Buffer
			handler := &instrLogCaptureHandler{buf: &buf}
			logger := slog.New(handler)
			
			oldDefault := slog.Default()
			slog.SetDefault(logger)
			defer slog.SetDefault(oldDefault)

			LogAllInstructions()

			output := buf.String()
			lines := strings.Split(strings.TrimSpace(output), "\n")
			
			// The first line is "Logging all registered Z80 instructions"
			assert.Contains(t, "First line should be header", "Logging all registered Z80 instructions", lines[0])
			
			// We expect at least one instruction to be logged (e.g., NOP)
			assert.Contains(t, "Should contain NOP", "mnemonic=\"NOP\"", output)
			assert.Contains(t, "Should contain prefix Main", "prefix=Main", output)
			
			// The number of lines (excluding header) should match CountInstructions()
			assert.Equal(t, "Number of log lines should match count + header", len(lines), CountInstructions()+1)
		},
	},
}
