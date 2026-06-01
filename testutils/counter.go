package testutils

import (
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var (
	assertCountsMu sync.Mutex
	assertCounts   = make(map[string]int)
)

// RecordAssertion incrémente le compteur pour le package du test
func RecordAssertion(t *testing.T) {
	assertCountsMu.Lock()
	defer assertCountsMu.Unlock()

	// t.Name() = "TestFuncName" ou "TestFuncName/subtest"
	// t.Package() n’existe pas → on utilise le nom du package depuis t
	// On récupère le path complet via runtime.Caller depuis t
	pkgPath := getPackagePathFromCaller()
	assertCounts[pkgPath]++
}

// ---------------- Helper pour déterminer le package ----------------
func getPackagePathFromCaller() string {
	// On remonte dans la stack pour trouver le fichier du test réel
	for i := 2; i < 10; i++ {
		_, file, _, ok := runtime.Caller(i)
		if !ok {
			continue
		}
		if strings.HasSuffix(file, "_test.go") {
			dir := filepath.Dir(file)
			dir = filepath.ToSlash(dir) // Unix-style
			// On coupe la partie avant "starspace/"
			if idx := strings.Index(dir, "starspace/"); idx >= 0 {
				return dir[idx:] // "starspace/internal/errors" etc.
			}
			return dir
		}
	}
	return "unknown"
}
