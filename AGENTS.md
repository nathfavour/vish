# Rules for AI Agents

To maintain a clean project structure, all agents must adhere to the following rules:

1. **Build Artifacts**: All compiled binaries and build artifacts MUST be placed in the `bin/` directory. NEVER output binaries directly into the project root or source directories.
2. **Go Version**: Always ensure the project uses Go 1.21. Check `go.mod` and `.go-version` before performing any Go-related tasks.
3. **Anyisland Integration**: Maintain the `anyisland.json` manifest and ensure any new features integrate with the Anyisland Pulse/Daemon where appropriate.
4. **Style**: Follow the existing directory structure: `cmd/` for entry points and `internal/` for logic.
