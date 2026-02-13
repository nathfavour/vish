package config

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Paths struct {
	ConfigDir string
	DataDir   string
	KeyDir    string
	RCPath    string
}

func GetPaths() (*Paths, error) {
	userConfig, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	vishConfig := filepath.Join(userConfig, "vish")
	
	// Create the config dir if it doesn't exist
	if _, err := os.Stat(vishConfig); os.IsNotExist(err) {
		err = os.MkdirAll(vishConfig, 0700)
		if err != nil {
			return nil, err
		}
	}

	// Setup a unique key dir if it doesn't exist
	// We store the UUID name in a pointer file if we are using local keys
	keyPointer := filepath.Join(vishConfig, ".keydir")
	var keyDirName string
	if data, err := os.ReadFile(keyPointer); err == nil {
		keyDirName = string(data)
	} else {
		keyDirName = uuid.New().String()
		_ = os.WriteFile(keyPointer, []byte(keyDirName), 0600)
	}

	keyDir := filepath.Join(vishConfig, keyDirName)
	if _, err := os.Stat(keyDir); os.IsNotExist(err) {
		_ = os.MkdirAll(keyDir, 0700)
	}

	return &Paths{
		ConfigDir: vishConfig,
		DataDir:   filepath.Join(vishConfig, "data"),
		KeyDir:    keyDir,
		RCPath:    filepath.Join(home, ".vishrc"),
	}, nil
}
