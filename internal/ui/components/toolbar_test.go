// Package components implements tests for toolbar components.
package components

import (
	"mcs/testutils/dsl"
	"testing"
)

// TestToolbar runs the scenario-based tests for the toolbar components.
func TestToolbar(t *testing.T) {
	dsl.RunScenarios(t, toolbarScenarios)
}
