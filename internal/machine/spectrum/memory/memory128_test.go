package memory

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestMemory128(t *testing.T) {
	dsl.RunScenarios(t, memory128Scenarios)
}
