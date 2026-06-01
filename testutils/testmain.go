package testutils

import (
	"os"
	"testing"
)

func RunWithAssertTracking(m *testing.M) {
	code := m.Run()
	ExportIfRequested()
	os.Exit(code)
}
