package history

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

type History struct {
	db *sql.DB
}

func NewHistory(path string) (*History, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return &History{db: db}, nil
}
