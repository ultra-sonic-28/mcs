// Package spectrum implements the ZX Spectrum machine logic.
package spectrum

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestSpectrum(t *testing.T) {
	dsl.RunScenarios(t, baseBusScenarios)
	dsl.RunScenarios(t, bus48Scenarios)
	dsl.RunScenarios(t, bus128Scenarios)
	dsl.RunScenarios(t, keyboardScenarios)
	dsl.RunScenarios(t, displayScenarios)
	dsl.RunScenarios(t, machineScenarios)
	dsl.RunScenarios(t, tapeScenarios)
}
