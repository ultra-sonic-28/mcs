//go:build !test_summary
// +build !test_summary

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type TestEvent struct {
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Output  string  `json:"Output"`
	Elapsed float64 `json:"Elapsed"`
}

type PackageSummary struct {
	Passed  int
	Failed  int
	Skipped int
	Total   int
}

type FailedTest struct {
	Package string
	Test    string
	Output  string
	Elapsed float64
}

func main() {
	// -------------------------
	// Options CLI
	// -------------------------
	var verboseMode bool

	flag.BoolVar(&verboseMode, "verbose", false, "Display tests per package and assertions per package")
	flag.Parse()

	coverageFile := "coverage.out"
	htmlFile := "coverage.html"

	rootDir, _ := os.Getwd()
	assertDir := filepath.Join(rootDir, ".asserts")

	_ = os.MkdirAll(assertDir, 0755)

	cmd := exec.Command(
		"go", "test", "./internal/...",
		"-json",
		//"-coverpkg=./internal/...",
		"-coverprofile="+coverageFile,
	)

	cmd.Env = append(os.Environ(),
		"ASSERT_STATS_DIR="+assertDir,
	)

	cmd.Env = append(cmd.Env,
		"GOFLAGS=-count=1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stdout)

	packageSummaries := make(map[string]*PackageSummary)
	globalPassed, globalFailed, globalSkipped := 0, 0, 0

	// Map pour stocker la sortie d'un test en cours
	testOutputs := make(map[string][]string)
	var failedTests []FailedTest

	for scanner.Scan() {
		line := scanner.Text()
		var ev TestEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}

		if ev.Test != "" && ev.Output != "" {
			// On concatène toutes les sorties pour ce test
			key := ev.Package + "/" + ev.Test
			testOutputs[key] = append(testOutputs[key], strings.TrimSpace(ev.Output))
		}

		if ev.Test == "" {
			continue // ignore package-level events
		}

		pkg := ev.Package
		if _, ok := packageSummaries[pkg]; !ok {
			packageSummaries[pkg] = &PackageSummary{}
		}

		s := packageSummaries[pkg]
		s.Total++

		switch ev.Action {
		case "pass":
			s.Passed++
			globalPassed++
			// on peut supprimer la sortie si le test passe
			key := ev.Package + "/" + ev.Test
			delete(testOutputs, key)
		case "fail":
			s.Failed++
			globalFailed++
			key := ev.Package + "/" + ev.Test
			// On stocke la sortie pour récupérer le message et la ligne
			failedTests = append(failedTests, FailedTest{
				Package: pkg,
				Test:    ev.Test,
				Output:  strings.Join(testOutputs[key], "\n"),
				Elapsed: ev.Elapsed,
			})
		case "skip":
			s.Skipped++
			globalSkipped++
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if err := cmd.Wait(); err != nil {
		// go test return non-zero si un test échoue, on ignore ici
	}

	stats, errStats := loadAssertStats()

	// Counts padding
	padding := 6

	// --- Résumé global ---
	totalTests := globalPassed + globalFailed + globalSkipped
	fmt.Println("\n\n================ Résumé Global =====================")
	fmt.Printf("Total tests      : %*d | %s%d%s | %s%d%s | %s%d%s | 🎯 %.1f%% passés\n",
		padding, totalTests,
		green("✅ "), globalPassed, reset(),
		red("❌ "), globalFailed, reset(),
		yellow("⚠️ "), globalSkipped, reset(),
		calcPassRate(globalPassed, globalFailed, globalSkipped),
	)
	if errStats == nil {
		fmt.Printf("Total assertions : %*d\n", padding, stats.Total)
	}

	if verboseMode {
		// Calculer la largeur max du nom du package
		maxNameLen := 0
		for pkg := range packageSummaries {
			if strings.HasSuffix(pkg, "/internal_test") {
				if len(pkg)-len("/internal_test") > maxNameLen {
					maxNameLen = len(pkg) - len("/internal_test")
				}
			} else {
				if len(pkg) > maxNameLen {
					maxNameLen = len(pkg)
				}
			}
		}

		// extraire + trier les packages
		packages := make([]string, 0, len(packageSummaries))
		for pkg := range packageSummaries {
			packages = append(packages, pkg)
		}
		sort.Strings(packages)

		// --- Résumé par package ---
		fmt.Println("\n\n================ Tests par Package ================")
		for _, pkg := range packages {
			s := packageSummaries[pkg]
			bar := buildInteractiveBar(s.Passed, s.Failed, s.Skipped)
			rate := calcPassRate(s.Passed, s.Failed, s.Skipped)
			if strings.HasSuffix(pkg, "/internal_test") {
				pkg = strings.Replace(pkg, "/internal_test", "", 1)
			}
			fmt.Printf("📦 %-*s: %s | %s%-4d%s | %s%d%s | %s%d%s | 🎯 %.1f%%\n",
				maxNameLen, pkg, bar,
				green("✅ "), s.Passed, reset(),
				red("❌ "), s.Failed, reset(),
				yellow("⚠️ "), s.Skipped, reset(),
				rate,
			)
		}

		fmt.Println("\n\n================ Assertions par Package =============")
		if errStats != nil {
			fmt.Println("⚠️ Assertions non disponibles")
		} else {
			packages := make([]string, 0, len(stats.PerPackage))
			for pkg := range stats.PerPackage {
				packages = append(packages, pkg)
			}
			sort.Strings(packages)

			for _, pkg := range packages {
				count := stats.PerPackage[pkg]
				if strings.HasSuffix(pkg, "/internal_test") {
					pkg = strings.Replace(pkg, "/internal_test", "", 1)
				}
				fmt.Printf("📦 %-*s: 🔢 %*d assertions\n", maxNameLen, pkg, padding, count)
			}
		}
	}

	// --- Tests échoués avec détails ---
	if len(failedTests) > 0 {
		fmt.Println("\n\n================ Tests Échoués =====================")
		for _, ft := range failedTests {
			fmt.Printf("%s%s%s - %s (%.3fs)\n", red("❌ "), ft.Package, reset(), ft.Test, ft.Elapsed)
			lines := strings.Split(ft.Output, "\n")
			for _, l := range lines {
				l = strings.TrimSpace(l)
				// garder uniquement les lignes qui ressemblent à des messages de test
				if l == "" {
					continue
				}
				if strings.Contains(l, "tests[") && strings.Contains(l, ":") {
					fmt.Printf("   %s\n", l)
				}
			}
		}

	}

	// --- Taux de couverture global ---
	coverageRate, totalStatements, totalCovered := readCoverageRate(coverageFile)
	fmt.Printf("\n\n================ Couverture du code =================\n")
	fmt.Printf("Taux de couverture global : %s (%d / %d)\n", colorCoverage(coverageRate), totalCovered, totalStatements)

	// --- Génération du HTML ---
	err = generateCoverageHTML(coverageFile, htmlFile)
	if err != nil {
		fmt.Println("Impossible de générer le HTML de couverture :", err)
	} else {
		fmt.Printf("Fichier HTML généré : %s\n", htmlFile)
	}
}

// -------------------- Fonctions utilitaires --------------------
type AssertStats struct {
	PerPackage map[string]int `json:"per_package"`
	Total      int            `json:"total"`
}

func loadAssertStats() (*AssertStats, error) {
	stats := &AssertStats{
		PerPackage: make(map[string]int),
	}

	files, err := filepath.Glob(".asserts/asserts.*.json")
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}

		var s AssertStats
		if json.Unmarshal(data, &s) != nil {
			continue
		}

		for pkg, count := range s.PerPackage {
			stats.PerPackage[pkg] += count
			stats.Total += count
		}
	}

	return stats, nil
}

func readCoverageRate(coverageFile string) (float64, uint32, uint32) {
	data, err := os.ReadFile(coverageFile)
	if err != nil {
		return 0, 0, 0
	}

	var totalStatements uint32 = 0
	var totalCovered uint32 = 0

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "mode:") {
			continue
		}
		var fileRange string
		var statements, covered uint32 = 0, 0

		_, err := fmt.Sscanf(line, "%s %d %d", &fileRange, &statements, &covered)
		if err != nil {
			continue
		}
		totalStatements += statements
		if covered == 1 {
			totalCovered += statements
		}
	}
	if totalStatements == 0 {
		return 0, 0, 0
	}
	return (float64(totalCovered) / float64(totalStatements) * 100), totalStatements, totalCovered
}

func colorCoverage(rate float64) string {
	switch {
	case rate >= 80:
		return green_cov(fmt.Sprintf(" %.1f%% ", rate))
	case rate >= 50:
		return yellow_cov(fmt.Sprintf(" %.1f%% ", rate))
	default:
		return red_cov(fmt.Sprintf(" %.1f%% ", rate))
	}
}

func generateCoverageHTML(profile, htmlFile string) error {
	cmd := exec.Command("go", "tool", "cover", "-html="+profile, "-o", htmlFile)
	return cmd.Run()
}

// Barre interactive avec ✅/❌/░
func buildInteractiveBar(passed, failed, skipped int) string {
	total := passed + failed + skipped
	if total == 0 {
		return "[" + strings.Repeat("░", 20) + "]"
	}
	const barWidth = 20
	passWidth := int(float64(passed) / float64(total) * float64(barWidth))
	failWidth := int(float64(failed) / float64(total) * float64(barWidth))
	emptyWidth := barWidth - passWidth - failWidth

	return fmt.Sprintf("[%s%s]",
		green(strings.Repeat("█", passWidth)),
		strings.Repeat("░", failWidth+emptyWidth),
	)
}

func calcPassRate(passed, failed, skipped int) float64 {
	total := passed + failed + skipped
	if total == 0 {
		return 0
	}
	return float64(passed) / float64(total) * 100
}

// Couleurs ANSI
func green(s string) string      { return "\033[32m" + s + "\033[0m" }
func red(s string) string        { return "\033[31m" + s + "\033[0m" }
func yellow(s string) string     { return "\033[33m" + s + "\033[0m" }
func green_cov(s string) string  { return "\033[97;42m" + s + "\033[0m" }
func red_cov(s string) string    { return "\033[97;41m" + s + "\033[0m" }
func yellow_cov(s string) string { return "\033[97;43m" + s + "\033[0m" }
func reset() string              { return "\033[0m" }
