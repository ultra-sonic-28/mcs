package tape

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestTape(t *testing.T) {
	dsl.RunScenarios(t, tapeScenarios)
}
