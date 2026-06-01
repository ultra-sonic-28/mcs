package logger

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestLogger(t *testing.T) {
	// Convert LoggerScenario to dsl.Scenario
	dslScenarios := make([]dsl.Scenario, len(loggerScenarios))
	for i, s := range loggerScenarios {
		dslScenarios[i] = dsl.NewScenario(s.Name, s.Run)
	}

	// Run scenarios using the mandated DSL
	dsl.RunScenarios(t, dslScenarios)
}
