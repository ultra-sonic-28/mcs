// Package components implements tests for statusbar components.
package components

import (
	"mcs/testutils/dsl"
	"testing"
)

// TestStatusbar runs the scenario-based tests for the statusbar component.
func TestStatusbar(t *testing.T) {
	dsl.RunScenarios(t, statusbarScenarios)
}
