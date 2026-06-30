// Package components implements reusable UI components for the MCS system.
package components

import (
	"image/color"
	"path/filepath"

	"mcs/internal/machine/spectrum/gui"

	"github.com/hajimehoshi/ebiten/v2"
)

// StatusbarHeight is the height of the statusbar in logical pixels.
const StatusbarHeight = 12

// statusbarSectionRatio defines the proportional widths of the statusbar sections.
// Values are percentages of total width: section1 (left), separator, section2, separator, section3 (right).
const (
	section1Ratio = 50 // Tape section: 50%
	section2Start = 65 // CPU and Machine sections start at 65%
)

// Statusbar represents a reusable statusbar UI component.
// It displays information in three sections: Tape, CPU, and Machine.
type Statusbar struct {
	width          int         // Width of the statusbar in logical pixels
	bgColor        color.Color // Background color
	textColor      color.Color // Text color
	separatorColor color.Color // Separator color
	section1Text   string      // Left section text (e.g., tape name)
	section2Text   string      // Middle section text (e.g., CPU name)
	section3Text   string      // Right section text (e.g., machine name)
	maxChars       int         // Maximum characters for truncation
}

// NewStatusbar creates a new Statusbar component with the given width and default styling.
// Default colors: background dark grey, text light grey, separator medium grey.
func NewStatusbar(width int) *Statusbar {
	if width < 0 {
		width = 0
	}
	return &Statusbar{
		width:          width,
		bgColor:        color.RGBA{32, 32, 32, 255},    // Dark grey
		textColor:      color.RGBA{200, 200, 200, 255}, // Light grey
		separatorColor: color.RGBA{100, 100, 100, 255}, // Medium grey
		section1Text:   "No tape",
		section2Text:   "Z80",
		section3Text:   "",
		maxChars:       20,
	}
}

// SetBackgroundColor sets the background color of the statusbar.
func (s *Statusbar) SetBackgroundColor(c color.Color) {
	s.bgColor = c
}

// SetTextColor sets the text color of the statusbar.
func (s *Statusbar) SetTextColor(c color.Color) {
	s.textColor = c
}

// SetSeparatorColor sets the separator color of the statusbar.
func (s *Statusbar) SetSeparatorColor(c color.Color) {
	s.separatorColor = c
}

// SetMaxChars sets the maximum number of characters displayed in section 1 before truncation.
// If maxChars is less than 5, it defaults to 5.
func (s *Statusbar) SetMaxChars(maxChars int) {
	if maxChars < 5 {
		maxChars = 5
	}
	s.maxChars = maxChars
}

// SetSection1 sets the text for the first (left) section. Commonly used for tape names.
func (s *Statusbar) SetSection1(text string) {
	s.section1Text = text
}

// SetSection2 sets the text for the second (middle) section. Commonly used for CPU name.
func (s *Statusbar) SetSection2(text string) {
	s.section2Text = text
}

// SetSection3 sets the text for the third (right) section. Commonly used for machine name.
func (s *Statusbar) SetSection3(text string) {
	s.section3Text = text
}

// SetTapeName is a convenience method that sets the first section to the base filename
// of the tape path. If the path is empty, it sets "No tape".
func (s *Statusbar) SetTapeName(tapePath string) {
	if tapePath == "" {
		s.section1Text = "No tape"
	} else {
		s.section1Text = filepath.Base(tapePath)
	}
}

// SetCPUName is a convenience method that sets the second section to the CPU name.
func (s *Statusbar) SetCPUName(cpuName string) {
	s.section2Text = cpuName
}

// SetMachineName is a convenience method that sets the third section to the machine name.
func (s *Statusbar) SetMachineName(machineName string) {
	s.section3Text = machineName
}

// Width returns the width of the statusbar in logical pixels.
func (s *Statusbar) Width() int {
	return s.width
}

// Height returns the height of the statusbar in logical pixels.
func (s *Statusbar) Height() int {
	return StatusbarHeight
}

// Draw renders the statusbar on the target image at the given Y position.
// The statusbar spans the full width and renders the three sections with separators.
func (s *Statusbar) Draw(screen *ebiten.Image, yPos int) {
	// Draw background
	statusRect := ebiten.NewImage(s.width, StatusbarHeight)
	statusRect.Fill(s.bgColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(yPos))
	screen.DrawImage(statusRect, op)

	// Calculate separator positions based on proportional widths
	sep1X := s.width / 2
	sep2X := s.width * section2Start / 100

	// Section 1: Tape/Device (left side)
	section1Text := s.truncateText(s.section1Text, sep1X)
	gui.DrawSmallText(screen, section1Text, 6, yPos+2, s.textColor)

	// Separator 1
	gui.DrawSmallText(screen, "|", sep1X, yPos+2, s.separatorColor)

	// Section 2: CPU (middle)
	gui.DrawSmallText(screen, s.section2Text, sep1X+10, yPos+2, s.textColor)

	// Separator 2
	gui.DrawSmallText(screen, "|", sep2X, yPos+2, s.separatorColor)

	// Section 3: Machine (right)
	gui.DrawSmallText(screen, s.section3Text, sep2X+6, yPos+2, s.textColor)
}

// truncateText truncates text to fit within the available space.
// Each character is approximately 6 pixels wide, with 10 pixels padding on each end.
func (s *Statusbar) truncateText(text string, availableWidth int) string {
	if availableWidth < 10 {
		return ""
	}
	maxWidth := availableWidth - 10
	maxChars := maxWidth / 6
	if maxChars < s.maxChars {
		maxChars = s.maxChars
	}

	if len(text) > maxChars {
		return text[:maxChars-3] + "..."
	}
	return text
}
