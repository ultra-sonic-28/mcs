# MCS - Multi CPUs System

## Project Overview & Objectives
**MCS** is a high-performance emulator framework designed to support multiple CPU architectures. The primary goal of this project is to provide a clean, modular, and idiomatic Go implementation for various legacy processors.

The current focus is the **Zilog Z80**, used as the foundational core for this project.

## Technical Stack
*   **Language**: Go (Golang) v1.25+ - Leveraging the latest performance improvements and type safety.
*   **Logging**: `log/slog` for structured, high-performance diagnostic logging.
*   **Automation**: `magefile` for build and task orchestration.
*   **Testing**: Custom Scenario-based DSL (Domain Specific Language) for robust, table-driven unit testing.

## Architecture Overview
The project follows a modular, interface-driven architecture to ensure flexibility and ease of extension:

*   **Registers**: Encapsulated state management for 8-bit and 16-bit registers, including alternate (shadow) sets.
*   **CPU Core**: Manages instruction cycles, interrupt flip-flops (IFF1/IFF2), and state transitions.
*   **Instruction Set**: Extensible opcode mapping system supporting:
    *   Standard Z80 opcodes (Main, CB, ED, DD, FD prefixes).
    *   **Z80N (ZX Spectrum Next)** extensions (e.g., `MUL`, `NEXTREG`).
*   **Buses**: Decoupled Memory and I/O communication through abstract interfaces:
    *   `Memory`: 16-bit addressable space (64KB).
    *   `IO`: 16-bit addressable port space.
*   **Concurrency**: Designed to eventually support multi-core or peripheral synchronization.

## Development & Build

### Prerequisites
*   Go **v1.25** or higher.
*   [Mage](https://magefile.org/) (**Highly Recommended**): While the project can be managed using standard Go tools, using Mage is recommended for automated build tasks and consistency.

### How to Build
The recommended way to build the project is using Mage:
```powershell
mage build
```
Alternatively, you can use the standard Go compiler:
```powershell
go build ./...
```

### How to Test
The project uses a strict DSL for testing. The recommended way to run tests is using Mage:
```powershell
mage test
```
Or using the standard Go test runner:
```powershell
go test ./internal/cpu/z80/... -v
```

## Emulated CPUs and machines

| CPU | Machine       | Emulated |
| --- | ------------- | :------: |
| Z80 |               | ✅ |
|     | Spectrum 48K  | ✅ |
|     | Spectrum 128K | ✅ |

## Known bugs
- Some 128K game tape not played perfectly
- Some demos from the demoscene are not played at all or inconsistantly

## License
This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for the full text.

## Attributions
### Icons / Images
- <a href="https://www.flaticon.com/free-icons/ram" title="ram icons">Ram icons created by Smashicons - Flaticon</a>
- <a href="https://www.flaticon.com/free-icons/shutdown" title="shutdown icons">Shutdown icons created by Fathema Khanom - Flaticon</a>

### Tapes
| Type | Title / Link  | Machine | Status | Comments |
| ---- | ------------- | ------- | :----: | -------- |
| Game | [Tetris - 40th Anniversary Edition](https://spectrumcomputing.co.uk/entry/45451/ZX-Spectrum/Tetris-40th_Anniversary_Edition) | Spectrum 48K | ✅ |  |
| Game | [Addix](https://spectrumcomputing.co.uk/entry/45446/ZX-Spectrum/Addix) | Spectrum 48K | ✅ | |
| Game | [4K Race Refueled](https://spectrumcomputing.co.uk/entry/18318/ZX-Spectrum/4K_Race_Refueled) | Spectrum 128K | ✅ | |
| Demo | [Forever 20 Invitation](https://spectrumcomputing.co.uk/entry/34414/ZX-Spectrum/Forever_20_Invitation) | Spectrum 48K<br />Spectrum 128K | ✅<br />✅ |  |
| Demo | [Gemba Demo](https://spectrumcomputing.co.uk/entry/27476/ZX-Spectrum/Gemba_Demo) | Spectrum 48K<br />Spectrum 128K | ✅<br />✅ |  |
| Demo | [Chaos](https://spectrumcomputing.co.uk/entry/34398/ZX-Spectrum/Chaos) | Spectrum 48K<br />Spectrum 128K | ✅<br />✅ | <ul><li>The sounds are messy and inaudible, but this is the same for all tested emulators</li></ul> |
| Demo | [Synergy 2024 Invitation](https://spectrumcomputing.co.uk/entry/43018/ZX-Spectrum/Synergy_2024_Invitation) | Spectrum 48K<br />Spectrum 128K | ✅<br />✅ |  |
| Demo | [Action](https://spectrumcomputing.co.uk/entry/7233/ZX-Spectrum/Action) | Spectrum 128K | ⚠️ | <ul><li>Graphics hang up in rotating checkerboard, but music continues</li></ul> |
| Demo | [Scrollerparty](https://www.pouet.net/prod.php?which=104454) | Spectrum 128k | ✅ |  |
| Demo | [Another Lightspeed Release](https://www.pouet.net/prod.php?which=94430) | Spectrum 128k | ✅ |  |
| Demo | [Kebugaruha](https://spectrumcomputing.co.uk/entry/43014) | Spectrum 128k | ✅ |  |