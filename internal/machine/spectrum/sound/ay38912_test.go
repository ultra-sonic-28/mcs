package sound

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestAY38912(t *testing.T) {
	dsl.RunScenarios(t, ay38912Scenarios)
}
