package keyboard

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestKeyboard(t *testing.T) {
	dsl.RunScenarios(t, keyboardScenarios)
}
