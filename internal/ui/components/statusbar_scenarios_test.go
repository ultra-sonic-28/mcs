// Package components defines the scenarios for statusbar tests.
package components

import (
	"image/color"
	"mcs/testutils/assert"
	"mcs/testutils/dsl"
	"testing"
)

var statusbarScenarios = []dsl.Scenario{
	dsl.NewScenario("Create Statusbar with valid width", func(t *testing.T) {
		sb := NewStatusbar(320)
		assert.Equal(t, "Width matches", sb.Width(), 320)
		assert.Equal(t, "Height equals StatusbarHeight", sb.Height(), StatusbarHeight)
	}),
	dsl.NewScenario("Create Statusbar with zero width", func(t *testing.T) {
		sb := NewStatusbar(0)
		assert.Equal(t, "Width is zero", sb.Width(), 0)
	}),
	dsl.NewScenario("Create Statusbar with negative width clamped to zero", func(t *testing.T) {
		sb := NewStatusbar(-100)
		assert.Equal(t, "Negative width clamped to 0", sb.Width(), 0)
	}),
	dsl.NewScenario("Statusbar initializes with default colors", func(t *testing.T) {
		sb := NewStatusbar(320)
		// Dark grey background: RGB(32, 32, 32)
		bgR, bgG, bgB, bgA := sb.bgColor.RGBA()
		assert.Equal(t, "BG Red channel", bgR, uint32(32*257))
		assert.Equal(t, "BG Green channel", bgG, uint32(32*257))
		assert.Equal(t, "BG Blue channel", bgB, uint32(32*257))
		assert.Equal(t, "BG Alpha channel", bgA, uint32(65535))

		// Light grey text: RGB(200, 200, 200)
		textR, textG, textB, textA := sb.textColor.RGBA()
		assert.Equal(t, "Text Red channel", textR, uint32(200*257))
		assert.Equal(t, "Text Green channel", textG, uint32(200*257))
		assert.Equal(t, "Text Blue channel", textB, uint32(200*257))
		assert.Equal(t, "Text Alpha channel", textA, uint32(65535))

		// Medium grey separator: RGB(100, 100, 100)
		sepR, sepG, sepB, sepA := sb.separatorColor.RGBA()
		assert.Equal(t, "Sep Red channel", sepR, uint32(100*257))
		assert.Equal(t, "Sep Green channel", sepG, uint32(100*257))
		assert.Equal(t, "Sep Blue channel", sepB, uint32(100*257))
		assert.Equal(t, "Sep Alpha channel", sepA, uint32(65535))
	}),
	dsl.NewScenario("Statusbar initializes with default texts", func(t *testing.T) {
		sb := NewStatusbar(320)
		assert.Equal(t, "Section1 default is 'No tape'", sb.section1Text, "No tape")
		assert.Equal(t, "Section2 default is 'Z80'", sb.section2Text, "Z80")
		assert.Equal(t, "Section3 default is empty", sb.section3Text, "")
	}),
	dsl.NewScenario("Statusbar initializes with default maxChars", func(t *testing.T) {
		sb := NewStatusbar(320)
		assert.Equal(t, "MaxChars defaults to 20", sb.maxChars, 20)
	}),
	dsl.NewScenario("SetBackgroundColor changes background color", func(t *testing.T) {
		sb := NewStatusbar(320)
		newColor := color.RGBA{255, 0, 0, 255} // Red
		sb.SetBackgroundColor(newColor)

		r, g, b, a := sb.bgColor.RGBA()
		assert.Equal(t, "BG Red set to 255", r, uint32(255*257))
		assert.Equal(t, "BG Green set to 0", g, uint32(0))
		assert.Equal(t, "BG Blue set to 0", b, uint32(0))
		assert.Equal(t, "BG Alpha set to 255", a, uint32(65535))
	}),
	dsl.NewScenario("SetTextColor changes text color", func(t *testing.T) {
		sb := NewStatusbar(320)
		newColor := color.RGBA{0, 255, 0, 255} // Green
		sb.SetTextColor(newColor)

		r, g, b, a := sb.textColor.RGBA()
		assert.Equal(t, "Text Red set to 0", r, uint32(0))
		assert.Equal(t, "Text Green set to 255", g, uint32(255*257))
		assert.Equal(t, "Text Blue set to 0", b, uint32(0))
		assert.Equal(t, "Text Alpha set to 255", a, uint32(65535))
	}),
	dsl.NewScenario("SetSeparatorColor changes separator color", func(t *testing.T) {
		sb := NewStatusbar(320)
		newColor := color.RGBA{0, 0, 255, 255} // Blue
		sb.SetSeparatorColor(newColor)

		r, g, b, a := sb.separatorColor.RGBA()
		assert.Equal(t, "Sep Red set to 0", r, uint32(0))
		assert.Equal(t, "Sep Green set to 0", g, uint32(0))
		assert.Equal(t, "Sep Blue set to 255", b, uint32(255*257))
		assert.Equal(t, "Sep Alpha set to 255", a, uint32(65535))
	}),
	dsl.NewScenario("SetMaxChars updates maxChars", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetMaxChars(30)
		assert.Equal(t, "MaxChars set to 30", sb.maxChars, 30)
	}),
	dsl.NewScenario("SetMaxChars with value less than 5 defaults to 5", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetMaxChars(3)
		assert.Equal(t, "MaxChars clamped to 5", sb.maxChars, 5)
	}),
	dsl.NewScenario("SetSection1 updates section 1 text", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetSection1("MyTape")
		assert.Equal(t, "Section1 text updated", sb.section1Text, "MyTape")
	}),
	dsl.NewScenario("SetSection2 updates section 2 text", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetSection2("6502")
		assert.Equal(t, "Section2 text updated", sb.section2Text, "6502")
	}),
	dsl.NewScenario("SetSection3 updates section 3 text", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetSection3("Commodore 64")
		assert.Equal(t, "Section3 text updated", sb.section3Text, "Commodore 64")
	}),
	dsl.NewScenario("SetTapeName with empty path sets 'No tape'", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetTapeName("")
		assert.Equal(t, "Section1 set to 'No tape' for empty path", sb.section1Text, "No tape")
	}),
	dsl.NewScenario("SetTapeName extracts base filename from path", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetTapeName("/path/to/mytape.tap")
		assert.Equal(t, "Section1 set to base filename", sb.section1Text, "mytape.tap")
	}),
	dsl.NewScenario("SetTapeName handles Windows paths", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetTapeName("C:\\Users\\test\\games\\game.tap")
		assert.Equal(t, "Section1 extracted from Windows path", sb.section1Text, "game.tap")
	}),
	dsl.NewScenario("SetCPUName updates section 2 text", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetCPUName("6502")
		assert.Equal(t, "Section2 set via SetCPUName", sb.section2Text, "6502")
	}),
	dsl.NewScenario("SetMachineName updates section 3 text", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetMachineName("Apple II")
		assert.Equal(t, "Section3 set via SetMachineName", sb.section3Text, "Apple II")
	}),
	dsl.NewScenario("truncateText returns text shorter than maxChars as-is", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetMaxChars(20)
		result := sb.truncateText("short", 100)
		assert.Equal(t, "Short text not truncated", result, "short")
	}),
	dsl.NewScenario("truncateText truncates long text with ellipsis", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetMaxChars(10)
		// Available width 100 pixels: (100-10)/6 = 15 chars max, but maxChars=10 restricts it
		// maxChars = 15, text length 46 > 15, so truncate to (15-3)=12 chars + "..."
		result := sb.truncateText("This is a very long text that should be truncated", 100)
		assert.Equal(t, "Long text truncated with '...'", result, "This is a ve...")
	}),
	dsl.NewScenario("truncateText with narrow space truncates when needed", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetMaxChars(10)
		// With narrow width 15: maxWidth=5, maxChars=0, but clamped to maxChars=10
		// Text length 12 > 10, so truncate to (10-3)=7 chars + "..."
		result := sb.truncateText("LongTextHere", 15)
		assert.Equal(t, "Text truncated in narrow space", result, "LongTex...")
	}),
	dsl.NewScenario("Width property is accessible", func(t *testing.T) {
		sb := NewStatusbar(512)
		assert.Equal(t, "Width accessor", sb.Width(), 512)
	}),
	dsl.NewScenario("Height property returns StatusbarHeight constant", func(t *testing.T) {
		sb := NewStatusbar(320)
		assert.Equal(t, "Height always StatusbarHeight", sb.Height(), 12)
	}),
	dsl.NewScenario("Multiple property updates work independently", func(t *testing.T) {
		sb := NewStatusbar(320)
		sb.SetSection1("Section1Text")
		sb.SetSection2("Section2Text")
		sb.SetSection3("Section3Text")
		sb.SetMaxChars(25)

		assert.Equal(t, "Section1 preserved", sb.section1Text, "Section1Text")
		assert.Equal(t, "Section2 preserved", sb.section2Text, "Section2Text")
		assert.Equal(t, "Section3 preserved", sb.section3Text, "Section3Text")
		assert.Equal(t, "MaxChars preserved", sb.maxChars, 25)
	}),
	dsl.NewScenario("Color properties persist after being set", func(t *testing.T) {
		sb := NewStatusbar(320)
		color1 := color.RGBA{100, 150, 200, 255}
		color2 := color.RGBA{50, 75, 100, 255}
		color3 := color.RGBA{200, 100, 50, 255}

		sb.SetBackgroundColor(color1)
		sb.SetTextColor(color2)
		sb.SetSeparatorColor(color3)

		bgR, _, _, _ := sb.bgColor.RGBA()
		textR, _, _, _ := sb.textColor.RGBA()
		sepR, _, _, _ := sb.separatorColor.RGBA()

		assert.Equal(t, "Background color persists", bgR, uint32(100*257))
		assert.Equal(t, "Text color persists", textR, uint32(50*257))
		assert.Equal(t, "Separator color persists", sepR, uint32(200*257))
	}),
}
