# AI Agent Instructions for vish

This document serves as the primary instructional context and governance for all AI agents interacting with the `vish` codebase.

## 1. Project Overview & Context
**vish** (The Vibe Coders Shell) is a modern, security-focused shell written in **Go**. 
- **Core Philosophy**: Security by Default, Instant Performance, and AI-Integrated Intent Analysis.
- **Key Technologies**: Go 1.21, Bubble Tea & Lip Gloss (TUI), `mvdan.cc/sh` (Parser), Modernc SQLite (History).
- **Distribution**: Integrated with [Anyisland](https://github.com/nathfavour/anyisland) for decentralized, sovereign management.

## 2. Mandatory Rules (Non-Negotiable)

### 2.1 Build Artifacts
- **Rule**: All compiled binaries and build artifacts MUST be placed in the `bin/` directory.
- **Enforcement**: NEVER output binaries directly into the project root. Check `anyisland.json` and CI workflows to ensure they respect this path.

### 2.2 Go Version Pinning
- **Rule**: The project is strictly pinned to **Go 1.21**.
- **Enforcement**: Always check `go.mod` and `.go-version` before performing Go operations. If `go mod tidy` bumps the version, immediately revert it using `go mod edit -go=1.21`.

### 2.3 Directory Structure & Style
- **Rule**: Follow the established modular structure.
- **Entry Points**: Located in `cmd/vish/`. Keep `main.go` lean.
- **Logic**: Located in `internal/`.
- **Imports**: Use relative imports for internal packages (e.g., `vish/internal/...`).

### 2.4 Anyisland Integration
- **Rule**: Maintain the `anyisland.json` manifest and Pulse integration.
- **Pulse**: vish must perform a handshake via `~/.anyisland/anyisland.sock`.
- **Registration**: Auto-register with the local Anyisland daemon via UDP port 1995 on startup.

## 3. Architecture Blueprint
- `internal/ui`: Custom prompt and interface logic.
- `internal/parser`: AST generation and command expansion.
- `internal/auditor`: Multi-tier security engine (Speed, Logic, Intent tiers).
- `internal/executor`: Process management and sandboxing.
- `internal/history`: SQLite-backed audit trail.
- `internal/ecosystem`: Anyisland Pulse/Daemon integration.

## 4. Operational Commands
- **Build**: `mkdir -p bin && go build -o bin/vish ./cmd/vish`
- **Run**: `./bin/vish`
- **Test**: `go test ./...`
- **Maintenance**: `go mod tidy && go mod edit -go=1.21`

## 5. Reference Files
- `ARCHITECTURE.json`: Detailed functional specifications for each module.
- `anyisland.json`: Manifest for distribution and OTA updates.
- `.go-version`: The authoritative Go version.
- `.github/workflows/go.yml`: CI/CD enforcement logic.
