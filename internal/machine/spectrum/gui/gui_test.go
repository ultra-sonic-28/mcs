// Package gui implements tests for the ZX Spectrum GUI logic.
package gui

import (
	"mcs/testutils/dsl"
	"testing"
)

// TestGui runs all the scenario-based tests defined in the guiScenarios list.
func TestGui(t *testing.T) {
	dsl.RunScenarios(t, guiScenarios)
}
