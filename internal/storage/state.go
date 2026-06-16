package storage

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type SessionState struct {
	ID        string
	Cwd       string
	Env       map[string]string
	ActiveCmd string
}

type Store struct {
	db *sql.DB
}

func NewStore(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		cwd TEXT,
		env TEXT,
		active_cmd TEXT
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) SaveSession(state SessionState) error {
	envData, _ := json.Marshal(state.Env)
	_, err := s.db.Exec(`
		INSERT INTO sessions (id, cwd, env, active_cmd)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			cwd = excluded.cwd,
			env = excluded.env,
			active_cmd = excluded.active_cmd;
	`, state.ID, state.Cwd, string(envData), state.ActiveCmd)
	return err
}

func (s *Store) LoadSession(id string) (*SessionState, error) {
	var state SessionState
	var envStr string
	err := s.db.QueryRow("SELECT id, cwd, env, active_cmd FROM sessions WHERE id = ?", id).
		Scan(&state.ID, &state.Cwd, &envStr, &state.ActiveCmd)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(envStr), &state.Env)
	return &state, nil
}
