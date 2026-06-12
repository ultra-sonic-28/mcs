package memory

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestMemory48(t *testing.T) {
	dsl.RunScenarios(t, memory48Scenarios)
}
