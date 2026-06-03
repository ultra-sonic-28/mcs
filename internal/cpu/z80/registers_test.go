package z80

import (
	"mcs/testutils/dsl"
	"testing"
)

func TestZ80(t *testing.T) {
	// Combine all register scenarios
	allRegScenarios := append([]RegisterScenario{}, register16Scenarios...)
	allRegScenarios = append(allRegScenarios, exchangeScenarios...)
	allRegScenarios = append(allRegScenarios, flagScenarios...)
	allRegScenarios = append(allRegScenarios, registersLogScenarios...)

	// Convert RegisterScenario to dsl.Scenario
	dslScenarios := make([]dsl.Scenario, 0, len(allRegScenarios)+len(cpuScenarios)+len(interruptModeScenarios))
	for _, s := range allRegScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add CPU scenarios
	for _, s := range cpuScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add Step scenarios
	for _, s := range stepScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add Fetch scenarios
	for _, s := range fetchScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add InterruptMode scenarios
	for _, s := range interruptModeScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add AddressingMode scenarios
	for _, s := range addressingModeScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add Interrupt scenarios
	dslScenarios = append(dslScenarios, interruptScenarios...)

	// Add Instruction scenarios
	for _, s := range instructionScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add ADD scenarios
	for _, s := range addScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add SUB scenarios
	for _, s := range subScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add LOGIC scenarios
	for _, s := range logicScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add INC/DEC scenarios
	for _, s := range incDecScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add Additional LD scenarios
	for _, s := range ldScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add ADC scenarios
	for _, s := range adcScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add SBC scenarios
	for _, s := range subcScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add PushPop scenarios
	for _, s := range pushPopScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add Jump scenarios
	for _, s := range jumpScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add Exchange scenarios
	for _, s := range exchangeInstructionScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add Rotation scenarios
	for _, s := range rotScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add BCD scenarios
	for _, s := range bcdScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add BIT scenarios
	for _, s := range bitScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add I/O scenarios
	for _, s := range ioScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Add Misc scenarios
	for _, s := range miscScenarios {
		dslScenarios = append(dslScenarios, dsl.NewScenario(s.Name, s.Run))
	}

	// Run scenarios using the mandated DSL
	dsl.RunScenarios(t, dslScenarios)
}
