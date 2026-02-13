package history

import (
	"database/sql"
	"path/filepath"
	"vish/internal/crypto"
	_ "modernc.org/sqlite"
)

type History struct {
	db *sql.DB
}

func NewHistory(configDir string) (*History, error) {
	dbPath := filepath.Join(configDir, "history.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table and index if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			command TEXT,
			cwd TEXT,
			exit_code INTEGER
		);
		CREATE INDEX IF NOT EXISTS idx_history_command ON history(command);
	`)
	if err != nil {
		return nil, err
	}

	return &History{db: db}, nil
}

func (h *History) Add(command, cwd string, exitCode int, key []byte) error {
	finalCmd := command
	if key != nil {
		encrypted, err := crypto.Encrypt(key, command)
		if err == nil {
			finalCmd = "ENC:" + encrypted
		}
	}
	_, err := h.db.Exec("INSERT INTO history (command, cwd, exit_code) VALUES (?, ?, ?)", finalCmd, cwd, exitCode)
	return err
}

var (
	hintCache     = make(map[string]string)
	hintCacheSize = 100
)

func (h *History) GetLastCommandLike(prefix string) (string, error) {
	if prefix == "" {
		return "", nil
	}

	// Check cache
	if cmd, ok := hintCache[prefix]; ok {
		return cmd, nil
	}

	var command string
	err := h.db.QueryRow("SELECT command FROM history WHERE command LIKE ? || '%' ORDER BY timestamp DESC LIMIT 1", prefix).Scan(&command)
	if err != nil {
		return "", err
	}

	// Simple cache management
	if len(hintCache) > hintCacheSize {
		// Clear cache if it gets too big (simple approach)
		hintCache = make(map[string]string)
	}
	hintCache[prefix] = command

	return command, nil
}

func (h *History) GetAll(limit int) ([]string, error) {
	rows, err := h.db.Query("SELECT command FROM history ORDER BY timestamp DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cmds []string
	for rows.Next() {
		var cmd string
		if err := rows.Scan(&cmd); err == nil {
			cmds = append(cmds, cmd)
		}
	}
	return cmds, nil
}
