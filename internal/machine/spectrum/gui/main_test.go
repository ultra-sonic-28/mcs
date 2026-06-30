// Package gui implements tests for the ZX Spectrum GUI logic.
package gui

import (
	"testing"

	"mcs/testutils"
)

// TestMain coordinates the execution of the test suite and handles assertion tracking.
func TestMain(m *testing.M) {
	testutils.RunWithAssertTracking(m)
}
