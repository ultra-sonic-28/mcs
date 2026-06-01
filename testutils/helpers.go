package testutils

import (
	"bytes"
	"os/exec"
	"strings"
)

// ListPackagesExcluding liste tous les packages Go correspondant au pattern donné
// et exclut ceux dont le chemin se termine par l'un des éléments de excludes.
// Exemple:
//
//	ListPackagesExcluding("./internal/...", []string{"testutils"})
func ListPackagesExcluding(pattern string, excludes []string) ([]string, error) {
	cmd := exec.Command("go", "list", pattern)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []string
	lines := bytes.Split(out, []byte{'\n'})
	for _, line := range lines {
		pkg := strings.TrimSpace(string(line))
		if pkg == "" {
			continue
		}

		excluded := false
		for _, ex := range excludes {
			if strings.HasSuffix(pkg, "/"+ex) {
				excluded = true
				break
			}
		}
		if !excluded {
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}
