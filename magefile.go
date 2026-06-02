//go:build mage
// +build mage

package main

import (
	"archive/zip"
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

const (
	mainCmd    = ".\\cmd\\mcs\\main.go"
	projectDir = "MCS"
)

// ///////////////////////////////////////////////////////////////////////////
// Target definition
// ///////////////////////////////////////////////////////////////////////////

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Run

// Configure and build MCS binary to /bin directory
func Build() error {
	if err := compileExamples(); err != nil {
		return err
	}

	// Incrémenter le numéro de build AVANT toute lecture
	fmt.Println("Updating build number")
	versionBytes, err := incrementBuildNumber("VERSION")
	if err != nil {
		return err
	}

	version := strings.TrimSpace(string(versionBytes))
	parts := strings.Split(version, ".")

	if len(parts) != 4 {
		// handle error: unexpected version format
		panic("invalid version string")
	}

	fmt.Println("Building MCS resource files for version", version)

	// 16x16: Menus, title bar, notification area, lists in "details" view.
	// 24x24: Interfaces with 125% scaling or some modern controls in Windows 10/11.
	// 32x32: Standard "icons" view in File Explorer, some shortcuts.
	// 48x48: Large icons in File Explorer or on the Desktop.
	// 256x256: Display with very large icons, high-resolution screens; Windows then scales down if necessary.
	tpl := `{
		"RT_GROUP_ICON": {
		  "APP": {
			"0000": [
			  "icon.png",
			  "icon16.png",
			  "icon24.png",
			  "icon32.png",
			  "icon48.png",
			  "icon64.png",
			  "icon128.png"
			]
		  }
		},
		"RT_MANIFEST": {
		  "#1": {
			"0409": {
			  "identity": {
				"name": "MCS",
				"version": "%[1]s"
			  },
			  "description": "",
			  "minimum-os": "win7",
			  "execution-level": "as invoker",
			  "ui-access": false,
			  "auto-elevate": false,
			  "dpi-awareness": "system",
			  "disable-theming": false,
			  "disable-window-filtering": false,
			  "high-resolution-scrolling-aware": false,
			  "ultra-high-resolution-scrolling-aware": false,
			  "long-path-aware": false,
			  "printer-driver-isolation": false,
			  "gdi-scaling": false,
			  "segment-heap": false,
			  "use-common-controls-v6": false
			}
		  }
		},
		"RT_VERSION": {
		  "#1": {
			"0000": {
			  "fixed": {
				"file_version": "%[1]s",
				"product_version": "%[1]s"
			  },
			  "info": {
				"0409": {
				  "Comments": "",
				  "CompanyName": "ultra-sonic-28",
				  "FileDescription": "MCS - Multi CPUs System",
				  "FileVersion": "%[1]s",
				  "InternalName": "MCS",
				  "LegalCopyright": "© 2026 - ultra-sonic-28 - MIT License",
				  "LegalTrademarks": "",
				  "OriginalFilename": "mcs.exe",
				  "PrivateBuild": "",
				  "ProductName": "MCS",
				  "ProductVersion": "%[1]s",
				  "SpecialBuild": ""
				}
			  }
			}
		  }
		}
	  }
`
	// Générer winres.json
	os.WriteFile(
		"./winres/winres.json",
		[]byte(fmt.Sprintf(tpl, version)),
		0644,
	)

	// Generate windows resource files for embbeding
	cmd := exec.Command("go-winres", "make")
	cmd.Run()

	if err := moveFile("./", "./cmd/mcs", "rsrc_windows_386.syso"); err != nil {
		log.Fatal(err)
	}
	if err := moveFile("./", "./cmd/mcs", "rsrc_windows_amd64.syso"); err != nil {
		log.Fatal(err)
	}

	// Build binaire
	fmt.Println("Building MCS binary...")
	now := time.Now()
	builddate := now.Format("2006-01-02")
	flags := fmt.Sprintf("-X main.Version=%s -X main.BuildDate=%s", version, builddate)
	cmd = exec.Command(
		"go",
		"build",
		"-ldflags", flags,
		"-o", "./bin/mcs.exe",
		"./cmd/mcs",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Build release archive for github
func Release() error {
	fmt.Println("Building release...")
	mg.Deps(Build)

	version, err := os.ReadFile("VERSION")
	if err != nil {
		return fmt.Errorf("cannot read VERSION file: %w", err)
	}
	ver := strings.TrimSpace(string(version))

	osName := mapOS(runtime.GOOS)
	archName := mapArch(runtime.GOARCH)

	zipName := fmt.Sprintf(
		"mcs-%s-%s-v%s.zip",
		osName,
		archName,
		ver,
	)

	// -------------------------------------------------
	// Copy assets
	// -------------------------------------------------
	fmt.Println("Copying assets...")
	srcAssets := "./assets"
	dstAssets := "./bin/assets"

	if _, err := os.Stat(srcAssets); os.IsNotExist(err) {
		fmt.Printf("Warning: source directory %s does not exist, skipping assets copy\n", srcAssets)
	} else {
		if err := copyDirFiltered(srcAssets, dstAssets); err != nil {
			return err
		}
	}

	// -------------------------------------------------
	// Create zip
	// -------------------------------------------------
	fmt.Println("Creating zip archive...")
	if err := os.MkdirAll("release", 0755); err != nil {
		return err
	}

	tmpZip := filepath.Join(os.TempDir(), zipName)
	if err := zipBinDir("./bin", tmpZip); err != nil {
		return err
	}

	finalZip := filepath.Join("release", zipName)
	if err := os.Rename(tmpZip, finalZip); err != nil {
		return err
	}

	fmt.Println("Computing SHA256...")
	if err := writeSHA256(finalZip); err != nil {
		return err
	}

	fmt.Println("Release created:", finalZip)

	return nil
}

// Launch MCS
func Run() error {
	cmd := exec.Command("go", "run", "./cmd/mcs/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running MCS...")
	return cmd.Run()
}

// Run unit tests with coverage support
func Test() error {
	cmd := exec.Command("go", "run", ".\\test_summary.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running tests...")
	fmt.Println(time.Now().Format("Monday, January 2, 2006 at 15:04:05"))
	return cmd.Run()
}

// Run unit tests in verbose mode (display tests per package and assertions per package) with coverage support
func TestVerbose() error {
	cmd := exec.Command("go", "run", ".\\test_summary.go", "-verbose")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running tests in verbose mode...")
	fmt.Println(time.Now().Format("Monday, January 2, 2006 at 15:04:05"))
	return cmd.Run()
}

// Run Z80 Instruction Set Exercisers (ZEXALL/ZEXDOC)
func Zex() error {
	fmt.Println("Running Z80 Instruction Set Exercisers...")
	cmd := exec.Command("go", "test", "-v", "-timeout", "20m", "-tags", "zex", "-run", "TestZex", "./internal/cpu/z80/...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Delete all .exe files, .log files, and clean bin/ and release/ directories
func Clean() error {
	root, _ := os.Getwd()

	fmt.Println("Cleaning project build artifacts and logs...")

	// 1. Remove directories
	dirs := []string{"bin", "release"}
	for _, d := range dirs {
		if _, err := os.Stat(d); err == nil {
			fmt.Println("Removing directory:", d)
			os.RemoveAll(d)
		}
	}

	// 2. Remove files via Walk
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore dossiers cachés
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		name := strings.ToLower(info.Name())
		// Supprime les .exe, les .log et les fichiers de debug
		if strings.HasSuffix(name, ".exe") || strings.HasSuffix(name, ".log") || strings.HasPrefix(name, "__debug_bin") {
			fmt.Println("Removing file:", path)
			return os.Remove(path)
		}

		return nil
	})
}

// Create a backup archive of the project directory (depends on Clean)
func Backup() error {
	mg.Deps(Clean)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	projectName := filepath.Base(cwd)

	now := time.Now()
	// Format: projectname-YYYYMMDD-HHhmm.zip
	timestamp := now.Format("20060102-15h04")
	zipName := fmt.Sprintf("%s-%s.zip", projectName, timestamp)

	// Le fichier de backup doit être créé dans le répertoire parent
	zipPath := filepath.Join("..", zipName)

	fmt.Printf("Creating project backup in parent directory: %s...\n", zipPath)

	return zipProject(".", zipPath, projectName)
}

// Install tools : go-winres, doc2go
func Tools() error {
	fmt.Println("Installing tools...")

	tools := map[string]string{
		"go-winres": "github.com/tc-hib/go-winres@v0.3.1",
		"doc2go":    "go.abhg.dev/doc2go@latest",
	}

	for name, path := range tools {
		if isToolInstalled(name) {
			fmt.Printf("✅ %s already installed\n", name)
			continue
		}

		fmt.Printf("Installing %s...\n", name)
		cmd := exec.Command("go", "install", path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install %s: %w", name, err)
		}
		fmt.Printf("✅ %s installed\n", name)
	}

	return nil
}

// Running sources analysis and statistics generation
func Stats() error {
	cloc_executable := "./.bintools/cloc-2.06.exe"
	report_file := "./.tmp/cloc_output_raw.md"
	clean_report_file := "./.tmp/cloc_output_clean.md"

	cmd := exec.Command(
		cloc_executable,
		"--skip-uniqueness",
		"--quiet",
		"--skip-archive=(zip|tar(.(gz|Z|bz2|xz|7z))?)",
		"--skip-win-hidden",
		//"--thousands-delimiter=_",
		//"--fmt=2",
		"--md",
		"--report-file="+report_file,
		"--found=./.tmp/found.txt",
		"--ignored=./.tmp/ignored.txt",
		"./cmd/",
		"./docs/",
		"./internal/",
		"./testutils/",
		"./winres/",
		"./GEMINi.md",
		"./CHANGELOG.md",
		"./README.md",
		"./*.go",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running stats...")
	err := cmd.Run()

	fmt.Println("Formating stats...")
	formatStatsFile()

	fmt.Println("")
	convertMarkdownTable2Text(clean_report_file)

	return err
}

// Generate HTML docs and serve it localy in browser
func DocsHTML() error {
	outDir := ".output/docs"
	fmt.Println("Generating MCS docs...")

	cmd := exec.Command("doc2go",
		"-out="+outDir,
		"-internal",
		"./...")
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	fmt.Printf("✅ Documentation MCS - Multi CPUs System : %s/index.html\n", outDir)
	return err
}

// Serve packages documentation in browser
func ServeDocs() error {
	dir := ".output/docs" // Dossier à servir
	port := ":8000"

	if len(os.Args) > 1 && os.Args[1] == "serve" {
		dir = os.Args[2] // mage serve ./mon-dossier
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("dossier %s introuvable", dir)
	}

	fmt.Printf("🚀 Serveur statique démarré : http://localhost:%s/mcs/\n", port)
	fmt.Printf("📁 Servir : %s\n", dir)
	fmt.Printf("⌨️  Ctrl+C pour arrêter\n\n")

	// FileServer pour tout servir
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	// Graceful shutdown
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("\n🛑 Signal reçu : %v\n", sig)
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Erreur shutdown : %v", err)
		}
		cancel()
	}()

	// Auto-ouvrir navigateur
	go openBrowser(fmt.Sprintf("http://localhost%s/mcs/", port))

	// Démarrer serveur
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("serveur échoué : %w", err)
	}

	fmt.Println("✅ Serveur arrêté proprement")
	return nil
}

// Compile Z80 examples if sources are newer than binaries
func compileExamples() error {
	fmt.Println("Checking for Z80 examples to compile...")

	type ExampleSet struct {
		CPU      string
		Compiler string
		Sources  []string
	}

	examples := ExampleSet{
		CPU:      "Z80",
		Compiler: filepath.Join(".bintools", "vasmz80_std.exe"),
		Sources:  []string{"add", "sub", "div", "multiply", "fact", "fibonacci", "prime_number"},
	}

	examplesDir := filepath.Join("assets", "z80", "examples")
	compilerPath, err := filepath.Abs(examples.Compiler)
	if err != nil {
		return err
	}

	if !fileExists(compilerPath) {
		return fmt.Errorf("compiler not found at %s", compilerPath)
	}

	for _, srcBase := range examples.Sources {
		srcName := srcBase + ".z80"
		outName := srcBase + ".out"
		srcPath := filepath.Join(examplesDir, srcName)
		outPath := filepath.Join(examplesDir, outName)

		srcStat, err := os.Stat(srcPath)
		if err != nil {
			return fmt.Errorf("source file %s not found: %w", srcPath, err)
		}

		outStat, err := os.Stat(outPath)
		needsCompile := false
		if err != nil {
			if os.IsNotExist(err) {
				needsCompile = true
			} else {
				return err
			}
		} else {
			if srcStat.ModTime().After(outStat.ModTime()) {
				needsCompile = true
			}
		}

		if needsCompile {
			fmt.Printf("Compiling %s -> %s\n", srcName, outName)
			cmd := exec.Command(compilerPath, "-Fbin", "-o", outName, srcName)
			cmd.Dir = examplesDir
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("compilation failed for %s: %w\n%s", srcName, err, string(output))
			}
		} else {
			fmt.Printf("Example %s is up to date.\n", outName)
		}
	}

	return nil
}

// ///////////////////////////////////////////////////////////////////////////
// Utility functions
// ///////////////////////////////////////////////////////////////////////////

func moveFile(srcDir, dstDir, name string) error {
	srcPath := filepath.Join(srcDir, name)
	dstPath := filepath.Join(dstDir, name)

	// Ensure destination directory exists.
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return err
	}

	return os.Rename(srcPath, dstPath)
}

// formatThousands insère _ tous les trois chiffres à partir de la droite.
func formatThousands(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	out := []byte{}
	count := 0
	for i := n - 1; i >= 0; i-- {
		out = append(out, s[i])
		count++
		if count%3 == 0 && i != 0 {
			out = append(out, '_')
		}
	}
	// reverse
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return string(out)
}

// Capitalise le premier caractère de chaque "mot" entre les pipes.
func capitalizeHeader(line string) string {
	parts := strings.Split(line, "|")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) > 0 {
			runes := []rune(p)
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}
	return strings.Join(parts, "|")
}

func formatStatsFile() {
	input := "./.tmp/cloc_output_raw.md"
	output := "./.tmp/cloc_output_clean.md"

	in, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)

	lineNum := 0
	reNumber := regexp.MustCompile(`\b\d{1,3}(\d{3})*\b`)

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Ignore les trois premières lignes du fichier
		if lineNum <= 3 {
			continue
		}

		// Supprime la ligne de tirets "--------|--------|..."
		if strings.HasPrefix(line, "--------|") {
			continue
		}

		// Capitaliser l'en-tête Language|files|blank|comment|code
		if strings.HasPrefix(line, "Language|") {
			line = capitalizeHeader(line)
		}

		// Substitutions de texte
		line = strings.ReplaceAll(line, "SUM:", "TOTAL:")

		// Formattage des nombres avec _
		line = reNumber.ReplaceAllStringFunc(line, func(num string) string {
			return formatThousands(num)
		})

		// Si la ligne contient "TOTAL:", mettre chaque cellule de TOTAL en gras
		if strings.HasPrefix(line, "TOTAL:|") || strings.HasPrefix(line, "**TOTAL:**|") {
			parts := strings.Split(line, "|")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
				if parts[i] != "" {
					parts[i] = fmt.Sprintf("**%s**", parts[i])
				}
			}
			line = strings.Join(parts, "|")
		}

		fmt.Fprintln(writer, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	writer.Flush()
}

func incrementBuildNumber(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(data))
	parts := strings.Split(version, ".")
	if len(parts) != 4 {
		return "", fmt.Errorf("invalid version string: %q", version)
	}

	build, err := strconv.Atoi(parts[3])
	if err != nil {
		return "", fmt.Errorf("invalid build number %q: %w", parts[3], err)
	}
	build++

	parts[3] = strconv.Itoa(build)
	newVersion := strings.Join(parts, ".")

	if err := os.WriteFile(path, []byte(newVersion+"\n"), 0644); err != nil {
		return "", err
	}

	return newVersion, nil
}

func mapOS(goos string) string {
	switch goos {
	case "windows":
		return "win"
	case "linux":
		return "linux"
	case "darwin":
		return "mac"
	default:
		return goos
	}
}

func mapArch(goarch string) string {
	switch goarch {
	case "amd64":
		return "x64"
	case "arm64":
		return "arm64"
	default:
		return goarch
	}
}

func copyDirFiltered(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if info.Name() == "_do_not_delete_.txt" {
			return nil
		}

		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func zipProject(baseDir, zipFile, prefix string) error {
	zf, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer zf.Close()

	w := zip.NewWriter(zf)
	defer w.Close()

	absZipFile, _ := filepath.Abs(zipFile)

	return filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ne pas sauter le répertoire de base lui-même
		if path == baseDir {
			return nil
		}

		// Ignore dossiers cachés (comme .git, .tmp, .vscode), SAUF .bintools
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") && info.Name() != ".bintools" {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		// Ne pas s'inclure soi-même
		absPath, _ := filepath.Abs(path)
		if absPath == absZipFile {
			return nil
		}

		rel, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		fh, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Utiliser le prefix comme répertoire racine dans le ZIP
		fh.Name = filepath.ToSlash(filepath.Join(prefix, rel))
		fh.Method = zip.Deflate

		out, err := w.CreateHeader(fh)
		if err != nil {
			return err
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		return err
	})
}

func zipBinDir(binDir, zipFile string) error {
	zf, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer zf.Close()

	w := zip.NewWriter(zf)
	defer w.Close()

	return filepath.Walk(binDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Skip config.json
		if info.Name() == "config.json" {
			return nil
		}

		// Chemin relatif à ./bin
		rel, err := filepath.Rel(binDir, path)
		if err != nil {
			return err
		}

		fh, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		fh.Name = filepath.ToSlash(rel)
		fh.Method = zip.Deflate

		out, err := w.CreateHeader(fh)
		if err != nil {
			return err
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		return err
	})
}

func writeSHA256(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	sum := hex.EncodeToString(h.Sum(nil))
	out := fmt.Sprintf("%s  %s\n", sum, filepath.Base(filePath))

	return os.WriteFile(filePath+".sha256", []byte(out), 0644)
}

type Align int

const (
	AlignLeft Align = iota
	AlignRight
)

func convertMarkdownTable2Text(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan: %v", err)
	}

	// Extraction naïve de la première table (lignes contenant des '|')
	var tableLines []string
	for _, l := range lines {
		if strings.Contains(l, "|") {
			tableLines = append(tableLines, l)
		}
	}
	if len(tableLines) == 0 {
		return
	}

	// 1. Parser les lignes en cellules
	var rowsRaw [][]string
	for _, line := range tableLines {
		parts := strings.Split(line, "|")
		for i := range parts {
			parts[i] = stripMarkdownBold(strings.TrimSpace(parts[i]))
		}
		rowsRaw = append(rowsRaw, parts)
	}

	if len(rowsRaw) < 2 {
		return
	}

	header := rowsRaw[0]
	sep := rowsRaw[1]

	// 2. Déduire l’alignement par colonne depuis la ligne sep
	alignments := make([]Align, len(sep))
	for i, col := range sep {
		col = strings.TrimSpace(col)
		left := strings.HasPrefix(col, ":")
		right := strings.HasSuffix(col, ":")
		switch {
		case right && !left:
			alignments[i] = AlignRight
		default:
			// ":", ":...:", ou rien → gauche
			alignments[i] = AlignLeft
		}
	}

	// 3. Construire la liste des lignes de données SANS la ligne de séparation
	var rows [][]string
	rows = append(rows, header)
	for i := 2; i < len(rowsRaw); i++ {
		rows = append(rows, rowsRaw[i])
	}

	if len(rows) < 2 {
		return
	}

	colCount := len(header)
	widths := make([]int, colCount)

	// 4. Calcul des largeurs max
	updateWidths := func(cells []string) {
		for i := 0; i < colCount && i < len(cells); i++ {
			if l := len(cells[i]); l > widths[i] {
				widths[i] = l
			}
		}
	}
	for _, r := range rows {
		updateWidths(r)
	}

	// 5. Fonctions pour dessiner les lignes avec bordures

	buildBorder := func(left, mid, right, fill rune) string {
		var b strings.Builder
		b.WriteRune(left)
		for c := 0; c < colCount; c++ {
			for i := 0; i < widths[c]+2; i++ { // +2 pour un espace de chaque côté du texte
				b.WriteRune(fill)
			}
			if c < colCount-1 {
				b.WriteRune(mid)
			}
		}
		b.WriteRune(right)
		return b.String()
	}

	rowToString := func(cells []string) string {
		var b strings.Builder
		b.WriteRune('│')
		for i := 0; i < colCount && i < len(cells); i++ {
			cell := cells[i]
			width := widths[i]
			b.WriteRune(' ')
			if alignments[i] == AlignRight {
				fmt.Fprintf(&b, "%*s", width, cell)
			} else {
				fmt.Fprintf(&b, "%-*s", width, cell)
			}
			b.WriteRune(' ')
			b.WriteRune('│')
		}
		return b.String()
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	// 6. Impression

	// ligne du haut
	fmt.Fprintln(w, buildBorder('┌', '┬', '┐', '─'))

	// entête
	fmt.Fprintln(w, rowToString(rows[0]))

	// séparation entête / reste
	fmt.Fprintln(w, buildBorder('├', '┼', '┤', '─'))

	// lignes de données (sauf dernière = total)
	dataEnd := len(rows) - 1
	for i := 1; i < dataEnd; i++ {
		fmt.Fprintln(w, rowToString(rows[i]))
	}

	// séparation avant le total
	if dataEnd > 1 {
		fmt.Fprintln(w, buildBorder('├', '┼', '┤', '─'))
	}

	// ligne de total (dernière ligne)
	fmt.Fprintln(w, rowToString(rows[dataEnd]))

	// ligne du bas
	fmt.Fprintln(w, buildBorder('└', '┴', '┘', '─'))
}

func stripMarkdownBold(s string) string {
	return strings.ReplaceAll(s, "**", "")
}

func isToolInstalled(name string) bool {
	// Récupérer GOPATH/bin et GOBIN
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("USERPROFILE"), "go")
	}
	binPaths := []string{
		filepath.Join(gopath, "bin"),
		os.Getenv("GOBIN"),
	}

	// Ajouter les PATH standards Windows
	binPaths = append(binPaths, os.Getenv("PATH"))

	// Tester l'existence du fichier
	exts := []string{"", ".exe"}
	for _, binPath := range strings.Split(binPaths[2], ";") { // PATH
		binPath = strings.TrimSpace(binPath)
		if binPath == "" {
			continue
		}
		for _, ext := range exts {
			exePath := filepath.Join(binPath, name+ext)
			if fileExists(exePath) {
				return true
			}
		}
	}
	return false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Start()
}
