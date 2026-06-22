// Package components implements tests for button components.
package components

import (
	"mcs/testutils/dsl"
	"testing"
)

// TestButton runs the scenario-based tests for the Button component.
func TestButton(t *testing.T) {
	dsl.RunScenarios(t, buttonScenarios)
}
