package testutils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type AssertStats struct {
	PerPackage map[string]int `json:"per_package"`
	Total      int            `json:"total"`
}

func ExportIfRequested() {
	dir := os.Getenv("ASSERT_STATS_DIR")
	if dir == "" {
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		return
	}

	// Trouver la racine du module (par ex. le dossier contenant go.mod)
	modRoot := findModuleRoot(wd)
	if modRoot == "" {
		// fallback : ancien comportement
		pkg := filepath.Base(wd)
		writeStats(dir, pkg)
		return
	}

	// chemin relatif au module, ex: "internal/foo/bar"
	rel, err := filepath.Rel(modRoot, wd)
	if err != nil {
		return
	}

	// normaliser pour le nom de fichier : remplacer les séparateurs
	relSafe := strings.ReplaceAll(rel, string(filepath.Separator), "_")
	// ex: "internal_foo_bar"
	writeStats(dir, relSafe)
}

func writeStats(dir, pkgKey string) {
	assertCountsMu.Lock()
	defer assertCountsMu.Unlock()

	total := 0
	for _, n := range assertCounts {
		total += n
	}

	stats := AssertStats{
		PerPackage: assertCounts,
		Total:      total,
	}

	filename := filepath.Join(dir, "asserts."+pkgKey+".json")
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	_ = enc.Encode(stats)
}

func findModuleRoot(start string) string {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
