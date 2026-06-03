// Package spectrum implements the ZX Spectrum 48K machine logic.
package spectrum

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestBus(t *testing.T) {
	dsl.RunScenarios(t, busScenarios)
	dsl.RunScenarios(t, keyboardScenarios)
}
