// Package main is the entry point for the MCS application.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"mcs/internal/bus"
	"mcs/internal/config"
	"mcs/internal/cpu/z80"
	"mcs/internal/logger"
	"os"
	"path/filepath"
)

var (
	Version   = "dev"
	BuildDate = "unknown"
)

func main() {
	fmt.Printf("MCS v%s, built on %s\n\n", Version, BuildDate)

	// --- 1. Command Line Options ---
	programPath := flag.String("program", "", "Path to the binary program to load (e.g., test or test.out)")
	flag.Parse()

	if *programPath == "" {
		fmt.Println("Usage: mcs --program <file>[.out]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Add .out extension if not present
	filePath := *programPath
	if filepath.Ext(filePath) == "" {
		filePath += ".out"
	}

	// --- 2. Load Configuration ---
	cfg, err := config.Load("config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️ failed to load config: %v\n", err)
		os.Exit(1)
	}

	// --- 3. Setup Logging ---
	cleanup, err := logger.Setup("mcs.log", cfg.LoggingEnabled, cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️ failed to setup logging: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	slog.Info("Starting MCS (Multi CPUs System)")

	// --- 4. Initialize Hardware ---

	sharedBus := bus.NewSimpleBus()
	cpu := z80.NewCPU(sharedBus, sharedBus)

	// Report loaded instructions
	z80.LogAllInstructions()
	slog.Info("CPU initialization complete", "instructions_loaded", z80.CountInstructions())

	// Display Initial State
	cpu.Regs.LogState()

	// --- 5. Load Program into Memory ---
	if err := loadProgram(sharedBus, filePath, 0x0000); err != nil {
		slog.Error("⚠️ failed to load program", "error", err)
		os.Exit(1)
	}

	// --- 6. Execute Program ---
	slog.Info("Starting program execution")
	for !cpu.Halted {
		cpu.Step()
	}
	slog.Info("Program execution finished (HALT reached)")

	// --- 5. Display Final State ---
	cpu.Regs.LogState()

	slog.Info("MCS shutdown")
}

// loadProgram reads a binary file from disk and loads it into memory at the specified offset.
func loadProgram(mem z80.Memory, filePath string, offset uint16) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	slog.Info("Loading program into memory", "file", filePath, "size", len(data), "address", fmt.Sprintf("0x%04X", offset))
	for i, b := range data {
		mem.Write(offset+uint16(i), b)
	}
	return nil
}
