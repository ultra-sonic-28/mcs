package logger

import (
	"testing"

	"mcs/testutils"
)

func TestMain(m *testing.M) {
	testutils.RunWithAssertTracking(m)
}
