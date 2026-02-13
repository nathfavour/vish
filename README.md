# vish - The Vibe Coders Shell

vish is a modern shell designed for security, performance, and AI integration.

## Architecture

- **Prompt & Interface**: Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).
- **Parser Engine**: Uses [mvdan.cc/sh](https://mvdan.cc/sh) for POSIX-compliant parsing.
- **Execution Engine**: Manages processes and PTYs.
- **Security**: Multi-tier "Vibe Auditor" for intent analysis.

## Development

### Prerequisites

- Go 1.21+ (Pinned to 1.21 in `go.mod`)

### Building

```bash
go build -o vish ./cmd/vish
```

### Running

```bash
./vish
```
