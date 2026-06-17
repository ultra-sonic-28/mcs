// Package config implements tests for configuration loading.
package config

import (
	"mcs/testutils/dsl"
	"testing"
)

// TestConfig runs the scenario-based tests for the config package.
func TestConfig(t *testing.T) {
	dsl.RunScenarios(t, configScenarios)
}
