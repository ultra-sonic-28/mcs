package testutils

import (
	"bytes"
	"os"
	"testing"
)

// CaptureStdout capture tout ce qui est écrit sur stdout pendant l'exécution de fn
func CaptureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}

	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	return buf.String()
}
