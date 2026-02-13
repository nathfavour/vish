package storage

import (
	"os"
	"path/filepath"
	"vish/internal/config"
	"vish/internal/crypto"
)

type RCManager struct {
	paths *config.Paths
}

func NewRCManager(paths *config.Paths) *RCManager {
	return &RCManager{paths: paths}
}

// GetEncryptionKey retrieves or generates the master key from the key dir
func (rm *RCManager) GetEncryptionKey() ([]byte, error) {
	keyPath := filepath.Join(rm.paths.KeyDir, "master.key")
	if data, err := os.ReadFile(keyPath); err == nil {
		return data, nil
	}

	// Generate and save new key
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(keyPath, key, 0600)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// SaveRC saves plain text data to .vishrc
func (rm *RCManager) SaveRC(content string) error {
	return os.WriteFile(rm.paths.RCPath, []byte(content), 0644)
}

// LoadRC loads plain text data from .vishrc
func (rm *RCManager) LoadRC() (string, error) {
	if _, err := os.Stat(rm.paths.RCPath); os.IsNotExist(err) {
		return "", nil
	}

	data, err := os.ReadFile(rm.paths.RCPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
