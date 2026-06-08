# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Auto-Start mechanism for Spectrum 48K: automatically types `LOAD "" : RUN` and executes programs when a tape is provided via the `--tape` flag. The addition of `: RUN` ensures compatibility with tapes that lack an auto-start line in their header.
- Refined keyboard macro timing and sequence to ensure reliable command entry, including a fix for host keyboard interference during auto-typing.
- Instant Load (Fast Load) support: traps the ROM's `LD-BYTES` routine (0x0556) to inject tape blocks directly into memory, bypassing the slow audio loading process. Supports multi-block tapes by allowing the trap to trigger multiple times as the loader script progresses.
- Automatic CPU state restoration after Instant Load: explicitly enables interrupts (IFF1/IFF2) and releases virtual keys after the final block is loaded to ensure a smooth transition to the game's execution.
- New `Read16` helper in Spectrum `Bus` for easier 16-bit memory access.
- Detailed debug logging for `.tap` file loading, including block types, names, lengths, auto-start lines for programs, and loading addresses for code blocks.
- Integrated Ebitengine (v2) for Spectrum 48K GUI window management and keyboard mapping.
- New `--machine spectrum` flag in `cmd/mcs/main.go` to boot the ZX Spectrum emulator.
- New `--tape <file>.tap` flag to automatically load and play cassettes on boot.
- Spectrum 48K Cassette support: `.tap` file parser and real-time pulse generation for EAR port.
- Support for Tape Pilot, Sync, and Data pulse timings (T-cycles).
- Beeper state tracking in Port 0xFE Out.
- DSL-based unit tests for `.tap` block parsing and pulse timing.
- Z80 CPU: Support for maskable (INT) and non-maskable (NMI) interrupts.
- Z80 CPU: Implementation of interrupt modes IM0, IM1, and IM2.
- Spectrum 48K Machine implementation with 50Hz frame loop and 3.5MHz timing logic.
- Automated 50Hz interrupt triggering at the start of each frame.
- DSL-based unit tests for Z80 interrupts and Spectrum frame timing.
- Spectrum 48K Video Engine: non-linear display memory mapping (256x192).
- Spectrum attribute system support: Ink, Paper, Bright, and Flash attributes.
- Spectrum 16-color palette implementation in `internal/machine/spectrum/display.go`.
- DSL-based unit tests for non-linear memory rendering and attribute logic.
- Spectrum 48K ULA I/O logic: Port 0xFE support for Border color, Beeper, and MIC state.
- Keyboard matrix scanning (40 keys) in `internal/machine/spectrum/keyboard.go`.
- Support for Tape/EAR input bit in Port 0xFE.
- DSL-based unit tests for keyboard scanning and ULA I/O.
- Spectrum 48K foundation in `internal/machine/spectrum`: `Bus` implementation with 16KB ROM and 48KB RAM support.
- Embedded Spectrum 48K ROM data in `assets/machines/spectrum/rom.go`.
- DSL-based unit tests for Spectrum `Bus` memory and I/O logic.
- New Z80 assembler examples in `assets\z80\examples\`: `fact.z80` (factorial calculation), `fibonacci.z80` (Fibonacci sequence), and `prime_number.z80` (prime number detection).
- Support for Z80 'Preliminary tests' which are executed before `zexdoc` and `zexall` to verify core instruction logic.
- Instruction counting for Z80 exercisers, providing total execution counts after each test completion.

### Changed
- Compilation of Z80 example programs is now performed during `mage build` instead of `mage release` to ensure assets are ready for development and testing.

## [0.0.0.14] - 2026-06-02

### Fixed
- Corrected `logging_enabled=false` behavior: it now correctly suppresses all terminal output and removes the `mcs.log` file if it exists.
- Refactored `DAA` (Decimal Adjust Accumulator) instruction with more accurate Half-Carry (H) and undocumented flag logic to pass `zexall` tests.
- Updated `RLCA`, `RRCA`, `RLA`, and `RRA` instructions to correctly update undocumented flags 3 and 5 from the result.
- Fixed `LDIR` and `LDDR` block instructions to correctly update undocumented flags 3 and 5 in every iteration.
- Updated `BIT` instructions in `DDCB` and `FDCB` tables to correctly set undocumented flags 3 and 5 based on the high byte of the effective address (MEMPTR behavior).
- Fixed `RES` and `SET` instructions in `DDCB` and `FDCB` tables to support undocumented "copy to register" behavior.

### Added
- New command-line option `--program <file>[.out]` to load external binary programs.
- Automatic `.out` extension appending for program files if not specified.
- Added some example Z80 binary files to the `assets/z80` directory: `add`, `sub`, `multiply`, and `div`. These files are available in Z80 assembler and binary formats (compiled with `vasm 1.8g` by Volker Barthelmann, whose documentation is available here: http://sun.hasenbraten.de/vasm/release/vasm.html).
- Integration of `config.json` to control logging enabling/disabling and logging level.
- Support for automatic creation of `config.json` with default values if the file is missing.
- New `internal/config` package to manage application configuration.
- Support for configurable logging levels (DEBUG, INFO, WARN, ERROR) in `internal/logger`.
- Unit tests for new logging configuration features.
- Implementation of standard and undocumented Z80 rotate and shift instructions for `CB`, `DDCB`, and `FDCB` tables: `RLC`, `RRC`, `RL`, `RR`, `SLA`, `SRA`, `SLL` (undocumented), and `SRL`.
- Support for "copy to register" behavior in `DDCB` and `FDCB` rotate, shift, `RES`, and `SET` instructions.
- Implementation of all 8 register mappings for `DDCB` and `FDCB` `BIT`, `RES`, and `SET` instructions.
- Support for the undocumented `SLL` instruction (Shift Left Logical) in `CB`, `DDCB`, and `FDCB` tables.
- SBC HL, rr instructions (BC, DE, HL, SP) support.
- LD (nn), dd and LD dd, (nn) instructions (BC, DE, HL, SP) support.
- LD I, A, LD R, A, LD A, I, and LD A, R instructions support.
- IN r, (C) and OUT (C), r instructions support for all 8-bit registers.
- CP IYH, CP IYL, and CP (IY+d) instructions support.
- AND, XOR, and OR instructions with IYH, IYL, and (IY+d) operands support.
- SUB A, IYH, SUB A, IYL, and SUB A, (IY+d) instructions support.
- SBC A, IYH, SBC A, IYL, and SBC A, (IY+d) instructions support.
- ADC A, IYH and ADC A, IYL instructions support.
- ADD A, IYH/IYL and ADD A, (IY+d) instructions support.
- LD (IY+d), H and LD (IY+d), L instructions support.
- LD IYH, IYH/IYL, LD H, (IY+d), LD IYL, IYH/IYL, and LD L, (IY+d) instructions support.
- LD E, IYH/IYL, LD E, (IY+d), LD IYH/IYL, E, and LD (IY+d), E instructions support.
- LD D, IYH/IYL, LD D, (IY+d), LD IYH/IYL, D, and LD (IY+d), D instructions support.
- LD C, IYH/IYL, LD C, (IY+d), LD IYH/IYL, C, and LD (IY+d), C instructions support.
- LD A, IYH/IYL, LD A, (IY+d), LD IYH/IYL, A, and LD (IY+d), A instructions support.
- LD (IY+d), n instruction support.
- INC (IX+d), DEC (IX+d), INC (IY+d), and DEC (IY+d) instructions support.
- LD IYH, n and LD IYL, n instructions support.
- INC IY, DEC IY, INC IYH, DEC IYH, INC IYL, and DEC IYL instructions support.
- LD IY, nn, LD (nn), IY, and LD IY, (nn) instructions support.
- LD SP, IX and LD SP, IY instructions support.
- JP (IX) and JP (IY) instructions support.
- Z80 registers structure and 16-bit accessors in `internal/cpu/z80`.
- Flag management for the Z80 F register.
- Z80 `CPU` structure with interrupt management (IFF1, IFF2, IM, NMI, INT).
- State management for CPU halt and T-cycle tracking.
- DSL-based unit tests for CPU state and reset logic.
- Memory and IO interfaces for Z80.
- Integration of Memory and IO interfaces into the `CPU` struct.
- Project `README.md` with goals, stack, and architectural overview.
- Updated `README.md` with MIT License information and Mage recommendation.
- Instruction management system with opcode tables and Z80N support.
- Initial set of Z80 instructions (`NOP`, `LD A, n`, `LD A, r`, `LD A, (HL)`, `ADD A, r`, `ADD A, n`, `ADD A, (HL)`).
- Type-safe `InterruptMode` enumeration (IM0, IM1, IM2).
- Specialized flag update utility for 8-bit addition (`UpdateFlagsAdd8`).
- Updated `CPU` struct to use `InterruptMode`.
- Refactored `InterruptMode` into a dedicated file `interrupt_mode.go` with `String()` support.
- DSL-based unit tests for `InterruptMode` constants and string representation.
- Main entry point in `cmd/mcs/main.go` with custom logging and CPU initialization.
- Reusable `internal/bus` package with `SimpleBus` implementation.
- DSL-based unit tests for `SimpleBus` memory and I/O operations.
- Centralized logging in `internal/logger` with a modular `Setup()` function.
- Unit tests for the custom `LogHandler` formatting.
- Utility functions `LogInstruction`, `LogAllInstructions` (DEBUG level), and `CountInstructions` for system observability.
- `HALT` instruction implementation.
- CPU `Step()` method for single-instruction execution with per-opcode debug tracing.
- Enhanced 3-line detailed register state logging (Main, Alternate, Control).
- Support for loading binary programs from disk and executing them until `HALT` in the main entry point.
- Type-safe `AddressingMode` enumeration for standard Z80 and Z80N modes.
- Integration of addressing mode metadata into the `Instruction` structure and debug logging.
- Precise architectural addressing mode assignments for all implemented instructions.
- DSL-based unit tests for `AddressingMode` string representation.
- DSL-based unit tests for the CPU `Step()` method.
- DSL-based unit tests for the CPU `FetchByte()` and `FetchWord()` methods.
- DSL-based unit tests for the `Registers.LogState()` method.
- DSL-based unit tests for `LogInstruction`, `LogAllInstructions`, and `CountInstructions`.
- Comprehensive unit tests for the `logger` package covering `Enabled`, `WithAttrs`, `WithGroup`, and `Setup`.
- DSL-based unit tests for all Z80 `ADD A, r` and `LD A, r` register variants.
- Prefix handling for `0xCB`, `0xED`, `0xDD`, `0xFD` in CPU `Step()` method.
- Implementation of standard Z80 `LD r, r'`, `LD r, n`, `LD r, (HL)`, `LD (HL), r`, `LD (HL), n`, and `LD dd, nn` instructions.
- Implementation of standard Z80 `LD (BC), A`, `LD A, (BC)`, `LD (DE), A`, and `LD A, (DE)` instructions.
- Implementation of standard Z80 8-bit Load instructions: `LD (nn), A`, `LD A, (nn)`, and index components `LD IXH, n`, `LD IXL, n`, `LD B, IXH`, `LD B, IXL`, `LD C, IXH`, `LD C, IXL`, `LD D, IXH`, `LD D, IXL`, `LD E, IXH`, `LD E, IXL`, `LD IXH, r`, `LD IXL, r`, `LD A, IXH`, `LD A, IXL`.
- Implementation of standard Z80 16-bit Load instructions: `LD (nn), HL`, `LD HL, (nn)`, `LD SP, HL`, `LD IX, nn`, `LD (nn), IX`, and `LD IX, (nn)`.
- Implementation of standard Z80 indexed Load instruction: `LD (IX+d), n`, `LD B, (IX+d)`, `LD C, (IX+d)`, `LD D, (IX+d)`, `LD E, (IX+d)`, `LD L, (IX+d)`, `LD A, (IX+d)`, and `LD (IX+d), r`.
- Implementation of standard Z80 16-bit addition (`ADD HL, rr`) for register pairs BC, DE, HL, and SP.
- Implementation of standard Z80 16-bit addition for index registers (`ADD IX, rr` and `ADD IY, rr`) for register pairs BC, DE, IX/IY, and SP.
- Implementation of standard Z80 8-bit addition with carry (`ADC A, s`) for registers, immediate, memory (HL), and indexed (IX+d, IY+d).
- Implementation of standard Z80 16-bit addition with carry (`ADC HL, rr`) for register pairs BC, DE, HL, and SP.
- Implementation of standard Z80 control flow instructions: `JP`, `JR`, `CALL`, `RET`, `RST p` (all variants), and `JP (HL)`.
- Implementation of standard Z80 exchange instructions: `EX AF, AF'`, `EX DE, HL`, `EXX`, `EX (SP), HL`, and `EX (SP), IX/IY`.
- Implementation of standard Z80 8-bit arithmetic/logical instructions: `SUB`, `SBC A, s`, `AND`, `OR`, `XOR`, and `CP`.
- Implementation of standard Z80 8-bit compare (`CP s`) for index components `IXH`, `IXL`, and indexed memory `(IX+d)`.
- Implementation of standard Z80 8-bit logical AND (`AND s`) for index components `IXH`, `IXL`, and indexed memory `(IX+d)`.
- Implementation of standard Z80 8-bit logical OR (`OR s`) for index components `IXH`, `IXL`, and indexed memory `(IX+d)`.
- Implementation of standard Z80 8-bit logical XOR (`XOR s`) for index components `IXH`, `IXL`, and indexed memory `(IX+d)`.
- Implementation of standard Z80 8-bit rotation instructions for the accumulator: `RLCA`, `RRCA`, `RLA`, and `RRA`.
- Implementation of standard Z80 8-bit increment (`INC (HL)`) and decrement (`DEC (HL)`) instructions.
- Implementation of standard Z80 BCD instructions: `DAA`, `RLD`, and `RRD`.
- Implementation of standard Z80 bitwise instruction: `CPL`.
- Implementation of standard Z80 flag manipulation instructions: `SCF` and `CCF`.
- Implementation of standard Z80 interrupt control instructions: `DI` and `EI`.
- Implementation of standard Z80 8-bit addition (`ADD A, s`) for index components `IXH`, `IXL`, and indexed memory `(IX+d)`.
- Implementation of standard Z80 8-bit addition with carry (`ADC A, s`) for index components `IXH` and `IXL`.
- Implementation of standard Z80 8-bit subtraction (`SUB A, s`) for index components `IXH`, `IXL`, and indexed memory `(IX+d)`.
- Implementation of standard Z80 8-bit subtraction with carry (`SBC A, s`) for index components `IXH`, `IXL`, and indexed memory `(IX+d)`.
- Implementation of standard Z80 I/O instructions: `IN A, (n)` and `OUT (n), A`.
- Implementation of `0xDD` and `0xFD` opcodes as `NOP` instructions in the main opcode table.
- Implementation of Z80N extended instructions: `ADD HL, A`, `ADD DE, A`, `ADD BC, A`, `ADD HL, nn`, `ADD DE, nn`, `ADD BC, nn`, and `PUSH nn`.
- Organized instructions and tests into group-specific files: `op_ld.go`/`op_ld_scenarios_test.go`, `op_add.go`/`op_add_scenarios_test.go`, `op_sub.go`/`op_sub_scenarios_test.go`, `op_logic.go`/`op_logic_scenarios_test.go`, `op_misc.go`, and `op_pushpop.go`/`op_pushpop_scenarios_test.go`.
- Implementation of standard Z80 `PUSH rr` and `POP rr` for all register pairs (BC, DE, HL, AF).
- Implementation of `PUSH IX/IY` and `POP IX/IY` using DD/FD prefixes.
- Implementation of standard Z80 8-bit increment (`INC r`) and decrement (`DEC r`) instructions for registers A, B, C, D, E, H, and L.
- Implementation of standard Z80 16-bit increment (`INC rr`) and decrement (`DEC rr`) instructions for register pairs BC, DE, HL, SP, and index registers IX/IY.
- Implementation of standard Z80 8-bit increment (`INC r`) and decrement (`DEC r`) instructions for index components `IXH`, `IXL`, `IYH`, and `IYL`.
- Specialized flag update utilities for 8-bit increment (`UpdateFlagsInc8`) and decrement (`UpdateFlagsDec8`).
- Updated `op_core.go` to register new `INC/DEC` instructions.
- DSL-based unit tests for all implemented Z80 and Z80N instructions.

### Changed
- Modified the `release` command in `magefile.go` to automatically compile Z80 assembly examples before packaging assets.
- Added a new `compileExamples` function in `magefile.go` to perform incremental compilation of Z80 examples using `vasmz80_std.exe`.
- Updated `magefile.go` to isolate Z80 exercisers in the `mage zex` command.
- Refactored `internal/cpu/z80/op_core.go` into group-specific files for better maintainability.
- Updated `op_core.go` to serve as a central registry for all instruction groups.
