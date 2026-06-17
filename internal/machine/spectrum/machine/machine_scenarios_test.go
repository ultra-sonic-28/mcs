// Package machine implements the ZX Spectrum machine logic.
package machine

import (
	"mcs/internal/config"
	"mcs/internal/cpu/z80"
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var machineScenarios = []dsl.Scenario{
	dsl.NewScenario("Machine RunFrame execution", func(t *testing.T) {
		m := NewMachine(nil)
		m.Reset()

		// Set PC to a HALT instruction to avoid executing random memory
		m.Bus.Write(0x0000, 0x76) // HALT
		m.CPU.IFF1 = true
		m.CPU.IM = z80.IM1

		m.RunFrame()

		// T-cycles should be at least CyclesPerFrame (69888)
		assert.True(t, "Total cycles should be >= 69888", m.CPU.Cycles >= 69888)
		assert.True(t, "Total cycles should be close to 69888", m.CPU.Cycles < 69900)
	}),
	dsl.NewScenario("Machine128 RunFrame execution", func(t *testing.T) {
		m := NewMachine128(nil)
		m.Reset()

		m.Bus.Write(0x0000, 0x76) // HALT
		m.CPU.IFF1 = true
		m.CPU.IM = z80.IM1

		m.RunFrame()

		// T-cycles should be at least CyclesPerFrame128 (70908)
		assert.True(t, "Total cycles should be >= 70908", m.CPU.Cycles >= 70908)
		assert.True(t, "Total cycles should be close to 70908", m.CPU.Cycles < 70920)
	}),
	dsl.NewScenario("Machine metadata and layout", func(t *testing.T) {
		m48 := NewMachine(nil)
		assert.Equal(t, "Machine 48K Name", m48.MachineName, "Spectrum 48K")
		w, h := m48.Layout(0, 0)
		assert.Equal(t, "Machine 48K Layout Height", h, 192+12)

		m128 := NewMachine128(nil)
		assert.Equal(t, "Machine 128K Name", m128.MachineName, "Spectrum 128K")
		w, h = m128.Layout(0, 0)
		assert.Equal(t, "Machine 128K Layout Width", w, 256)
		assert.Equal(t, "Machine 128K Layout Height", h, 192+12)
	}),
	dsl.NewScenario("Machine layout with border config", func(t *testing.T) {
		cfg := &config.Config{
			Display: config.DisplayConfig{
				Border: config.BorderConfig{
					Color: "#FF0000",
					Width: 15,
				},
			},
		}
		m := NewMachine(cfg)
		assert.Equal(t, "Border width", m.borderWidth, 15)
		w, h := m.Layout(0, 0)
		// ScreenWidth (256) + 2*15 = 286
		assert.Equal(t, "Layout Width", w, 286)
		// ScreenHeight (192) + StatusLineHeight (12) + 2*15 = 234
		assert.Equal(t, "Layout Height", h, 234)

		r, g, b, a := m.borderColor.RGBA()
		// RGBA returns values in [0, 65535]
		assert.Equal(t, "Border Color R", r, uint32(65535))
		assert.Equal(t, "Border Color G", g, uint32(0))
		assert.Equal(t, "Border Color B", b, uint32(0))
		assert.Equal(t, "Border Color A", a, uint32(65535))

		// Invalid color should fallback to default color #D6CDC9 (214, 205, 201)
		cfgInvalidColor := &config.Config{
			Display: config.DisplayConfig{
				Border: config.BorderConfig{
					Color: "invalid",
					Width: 10,
				},
			},
		}
		m2 := NewMachine(cfgInvalidColor)
		r2, g2, b2, a2 := m2.borderColor.RGBA()
		// 214 * 257 = 54998
		// 205 * 257 = 52685
		// 201 * 257 = 51657
		assert.Equal(t, "Fallback R", r2, uint32(54998))
		assert.Equal(t, "Fallback G", g2, uint32(52685))
		assert.Equal(t, "Fallback B", b2, uint32(51657))
		assert.Equal(t, "Fallback A", a2, uint32(65535))

		// Negative width should be clamped to 0
		cfgNeg := &config.Config{
			Display: config.DisplayConfig{
				Border: config.BorderConfig{
					Color: "#00FF00",
					Width: -10,
				},
			},
		}
		m3 := NewMachine(cfgNeg)
		assert.Equal(t, "Negative border width clamped to 0", m3.borderWidth, 0)
	}),
}

