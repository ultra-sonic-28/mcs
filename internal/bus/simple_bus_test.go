package bus

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestSimpleBus(t *testing.T) {
	// Convert BusScenario to dsl.Scenario
	dslScenarios := make([]dsl.Scenario, len(busScenarios))
	for i, s := range busScenarios {
		dslScenarios[i] = dsl.NewScenario(s.Name, s.Run)
	}

	// Run scenarios using the mandated DSL
	dsl.RunScenarios(t, dslScenarios)
}
