# MCS: Technical Manifesto & Architectural Guidelines

## 1. Context and Project Overview
**MCS - Multi CPUs System** is a system that allows you to emulate different CPUs (e.g., Z80, 6502...)

The goal of this project is to create a multi-CPU emulator. The first CPU to emulate is the Zilog Z80.

## 2. Technology Stack
*   **Expert Golang Development**: Go (Golang) v1.25+ - Maximum exploitation of performance and type safety.
*   **Logging System**: `log/slog` with customized infrastructure.
*   **Resources**: `embed` for binary integration of assets.

## 3. Development Conventions

### Code Standards & Formatting
*   **Gofmt**: Code formatting must strictly adhere to the official `gofmt` tool standard.
*   **No Deprecated Code**: The use of deprecated functions, types, or methods is strictly prohibited. Always prefer current, supported APIs. If a library update introduces deprecations, perform the necessary migration before proceeding with new features.
*   **Verification**: Every modification to the source code **must** be immediately followed by a compilation phase (`go build` or `mage build`) to detect any warnings or potential errors, such as missing imports or undefined symbols.
*   **Changelog**: The `CHANGELOG.md` file **must** be systematically and automatically updated with every code modification, addition, or fix.
    *   **Language**: All entries must be written in **English**.
    *   **Format**: The content must strictly adhere to the [Keep A Changelog](https://keepachangelog.com/en/1.1.0/) standard.
*   **Naming Conventions**:
    *   **CamelCase** mandatory for variables, structures, and interfaces.
    *   **Visibility**: Initial uppercase for `Public`, lowercase for `private`.
    *   **Packages**: Lowercase, singular, no underscores (e.g., `cpu`, `memory`, `ìo`).
    *   **Files**: Lowercase, `snake_case` allowed. The name must reflect responsibility (e.g., `op_arith16.go`, `accumulator.go`).
    *   **Language**: Code (variables, functions, types) is in English.
*   **Clean Code**: Avoid unnecessary Getters/Setters. Expose fields or use semantic methods.

### Error Handling
*   **Explicit**: Every error must be handled or propagated. No "silent failure".
*   **Logging**: 
    * All system logs must be in English for compatibility with international analysis tools. 
    * Use helpers to facilitate logging operations.
    * The log file is truncated at each startup (`os.O_TRUNC`).
    * Format: `YYYY-MM-DD HH:MM:SS,ms [LEVEL] AppName: Message, key1=val1, key2=val2, ...`.
    * Every major object creation (Star, Planet, etc.) or resource detection must be logged with its full attributes.
    * Never automatically delete a write instruction from the log file

### Documentation (Godoc Standards)
*   **Language**: All technical documentation and code comments must be in **English**.
*   **Mandatory Commenting**: Every package, function, and type defined in the source code **must** be accompanied by a descriptive comment in English. These comments must clearly explain the functionality, responsibility, and usage of the element.
*   **Package Doc**: Each package must include a header comment describing its responsibility.
*   **Multiline**: Multiline comments are authorized and encouraged for complex algorithms.
*   **Compatibility**: Documentation must be perfectly rendered by the `godoc` tool.

### Testing Strategy
*   **Table-Driven Tests**: Systematically use this pattern for unit tests to ensure exhaustive coverage of edge cases.
*   **DSL Pattern**: All unit tests **must** strictly adhere to the scenario-based pattern defined in `testutils/dsl/dsl.go`.
*   **Test File Separation**: Unit tests **must** be split into two separate files:
    *   `*_test.go`: Contains the test runners and execution logic.
    *   `*_scenarios_test.go`: Contains only the test data, structures, and scenario definitions (DSL).
*   **Assertion Helpers**: It is **mandatory** to use the dedicated assertion helpers found in `testutils/assert/assert.go` (e.g., `assert.Equal`, `assert.True`, `assert.NotEqual`). Using standard `t.Errorf`, `t.Fatalf`, or `reflect.DeepEqual` directly is strictly prohibited to ensure consistency and proper tracking.

## 4. Project Structure and Responsibilities

```text
C:\My Program Files\MCS\
├── cmd/
│   └── mcs/
│       ├── main.go                 # Entry point. System initialization and main loop.
│       ├── rsrc_windows_386.syso   # Windows resources (compiled).
│       └── rsrc_windows_amd64.syso # Windows resources (compiled).
├── testutils/
│   ├── assert/
│   │   └── assert.go               # Custom assertion helpers for tests.
│   ├── dsl/
│   │   └── dsl.go                  # Domain-specific language for scenario-based testing.
│   ├── capture_stdout.go           # Utility to capture stdout during tests.
│   ├── counter.go                  # Concurrent counter for test tracking.
│   ├── export.go                   # Global test export utilities.
│   ├── helpers.go                  # Shared test helper functions.
│   ├── runcases.go                 # Scenario runner implementation.
│   └── testmain.go                 # Global test suite entry point.
├── winres/
│   └── winres.json                 # Windows resource configuration.
├── config.json                     # Persistence for global configuration.
├── GEMINI.md                       # Project technical manifesto and guidelines.
├── go.mod                          # Go module definition.
├── go.sum                          # Go module checksums.
├── LICENSE                         # Project license.
├── magefile.go                     # Build automation and project task manager.
├── README.md                       # Project documentation and architectural overview.
├── mcs.code-workspace              # VS Code workspace configuration.
├── test_summary.go                 # Utility for test result aggregation.
├── TODO.md                         # Project roadmap and pending tasks.
├── tools.go                        # Tool dependencies for the project.
└── VERSION                         # Current semantic version string.
```

**IMPORTANT**:
- Never invent gameplay rules
- Always refer to these documents before implementing a system
- If information is missing → ask before doing anything

