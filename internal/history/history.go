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

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			command TEXT,
			cwd TEXT,
			exit_code INTEGER
		)
	`)
	if err != nil {
		return nil, err
	}

	return &History{db: db}, nil
}

func (h *History) GetLastCommandLike(prefix string) (string, error) {
	if prefix == "" {
		return "", nil
	}
	var command string
	err := h.db.QueryRow("SELECT command FROM history WHERE command LIKE ? || '%' ORDER BY timestamp DESC LIMIT 1", prefix).Scan(&command)
	if err != nil {
		return "", err
	}
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
