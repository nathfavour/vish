# GEMINI.md - Context for vish

## Project Overview
**vish** (The Vibe Coders Shell) is a modern, security-focused shell written in **Go**. Its core philosophy centers on "Security by Default," "Instant Performance," and "AI-Integrated Intent Analysis."

### Key Technologies
- **Language**: Go 1.21 (strictly pinned).
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).
- **Parser**: [mvdan.cc/sh](https://mvdan.cc/sh) for POSIX/Bash syntax compatibility.
- **Database**: `modernc.org/sqlite` for immutable history logging.
- **Distribution**: Integrated with [Anyisland](https://github.com/nathfavour/anyisland) for decentralized package management.

### Architecture
The project follows a modular internal structure:
- `cmd/vish`: Entry point and main TUI loop.
- `internal/ui`: Custom prompt rendering and interface logic.
- `internal/parser`: AST generation and command expansion.
- `internal/auditor`: Multi-tier security engine (Speed, Logic, Intent tiers).
- `internal/executor`: Process management and sandboxing (using `os/exec` and `creack/pty`).
- `internal/history`: SQLite-backed audit trail and search.
- `internal/ecosystem`: Integration with Anyisland (Pulse handshake and Daemon registration).

---

## Building and Running

### Commands
- **Build**: `mkdir -p bin && go build -o bin/vish ./cmd/vish`
- **Run**: `./bin/vish`
- **Test**: `go test ./...`
- **Tidy Dependencies**: `go mod tidy && go mod edit -go=1.21`

### CI/CD
- GitHub Actions workflow at `.github/workflows/go.yml` enforces Go 1.21 and builds into the `bin/` directory.

---

## Development Conventions

### Strict Rules (from AGENTS.md)
1. **Build Artifacts**: All binaries MUST be placed in the `bin/` directory. NEVER output them to the project root.
2. **Go Version**: Always use Go 1.21. Verify `go.mod` and `.go-version` before tasks.
3. **Project Structure**: Entry points in `cmd/`, logic in `internal/`. Avoid bloating `main.go`.

### Anyisland Integration
- **Manifest**: `anyisland.json` controls the build process and Pulse status.
- **Pulse**: vish performs a handshake with the Anyisland daemon via Unix Domain Socket (`~/.anyisland/anyisland.sock`).
- **Registration**: Auto-registers with the local Anyisland daemon via UDP port 1995 on startup.

---

## Key Files
- `ARCHITECTURE.json`: Detailed blueprint of the system modules and responsibilities.
- `anyisland.json`: Manifest for distribution and OTA updates.
- `AGENTS.md`: Governance rules for AI agents.
- `.go-version`: Pinned Go version for environment managers.
- `internal/ecosystem/anyisland.go`: Logic for ecosystem integration.
