package machine

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestMachine(t *testing.T) {
	dsl.RunScenarios(t, machineScenarios)
}
