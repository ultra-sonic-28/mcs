package display

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestDisplay(t *testing.T) {
	dsl.RunScenarios(t, displayScenarios)
}
