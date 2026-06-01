package z80

import (
	"mcs/testutils/assert"
	"testing"
)

type InterruptModeScenario struct {
	Name string
	Run  func(t *testing.T)
}

var interruptModeScenarios = []InterruptModeScenario{
	{
		Name: "InterruptMode Constants",
		Run: func(t *testing.T) {
			assert.Equal(t, "IM0 should be 0", int(IM0), 0)
			assert.Equal(t, "IM1 should be 1", int(IM1), 1)
			assert.Equal(t, "IM2 should be 2", int(IM2), 2)
		},
	},
	{
		Name: "InterruptMode String Representation",
		Run: func(t *testing.T) {
			assert.Equal(t, "IM0 string", IM0.String(), "IM0")
			assert.Equal(t, "IM1 string", IM1.String(), "IM1")
			assert.Equal(t, "IM2 string", IM2.String(), "IM2")
			assert.Equal(t, "Unknown IM string", InterruptMode(99).String(), "Unknown")
		},
	},
}
