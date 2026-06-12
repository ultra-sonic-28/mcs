package bus

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestBus(t *testing.T) {
	dsl.RunScenarios(t, baseBusScenarios)
	dsl.RunScenarios(t, bus48Scenarios)
	dsl.RunScenarios(t, bus128Scenarios)
}
