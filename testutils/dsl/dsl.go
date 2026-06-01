package dsl

import "testing"

type Scenario struct {
	Name string
	Run  func(t *testing.T)
}

func NewScenario(name string, fn func(*testing.T)) Scenario {
	return Scenario{
		Name: name,
		Run:  fn,
	}
}

func RunScenarios(t *testing.T, scenarios []Scenario) {
	t.Helper()

	for _, s := range scenarios {
		s := s

		t.Run(s.Name, func(t *testing.T) {
			//t.Parallel()
			s.Run(t)
		})
	}
}

// Table-driven Scenario generator
func GenerateScenarios[T any](cases []T, nameFn func(T) string, runFn func(*testing.T, T)) []Scenario {
	scenarios := make([]Scenario, len(cases))
	for i, c := range cases {
		c := c
		scenarios[i] = Scenario{
			Name: nameFn(c),
			Run: func(t *testing.T) {
				t.Parallel()
				runFn(t, c)
			},
		}
	}
	return scenarios
}
