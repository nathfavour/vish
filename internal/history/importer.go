package history

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ImportExternalHistory looks for .zsh_history or .bash_history and imports the latest one
func (h *History) ImportExternalHistory() error {
	home, _ := os.UserHomeDir()
	
	// Check if we've already imported from this user
	var count int
	h.db.QueryRow("SELECT COUNT(*) FROM history WHERE cwd = 'imported'").Scan(&count)
	if count > 0 {
		return nil
	}

	zshPath := filepath.Join(home, ".zsh_history")
	bashPath := filepath.Join(home, ".bash_history")

	zshStat, zshErr := os.Stat(zshPath)
	bashStat, bashErr := os.Stat(bashPath)

	var targetPath string
	if zshErr == nil && bashErr == nil {
		if zshStat.ModTime().After(bashStat.ModTime()) {
			targetPath = zshPath
		} else {
			targetPath = bashPath
		}
	} else if zshErr == nil {
		targetPath = zshPath
	} else if bashErr == nil {
		targetPath = bashPath
	}

	if targetPath == "" {
		return nil
	}

	file, err := os.Open(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	tx, err := h.db.Begin()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Clean zsh history format ": 1234567890:0;cmd"
		if strings.HasPrefix(line, ": ") {
			parts := strings.SplitN(line, ";", 2)
			if len(parts) > 1 {
				line = parts[1]
			}
		}
		if strings.TrimSpace(line) != "" {
			_, _ = tx.Exec("INSERT INTO history (command, cwd, exit_code) VALUES (?, ?, ?)", line, "imported", 0)
		}
	}
	return tx.Commit()
}
