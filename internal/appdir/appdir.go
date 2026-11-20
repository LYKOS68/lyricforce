// Package appdir resolves the per-user data root of the application:
// %APPDATA%\LykosReimarchiv\ holding the database and runtime logs.
package appdir

import (
	"fmt"
	"os"
	"path/filepath"
)

const dirName = "LykosReimarchiv"

// Root returns the per-user data directory (not created).
func Root() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(dir, dirName), nil
}

// DatabasePath returns the path of the SQLite archive.
func DatabasePath() (string, error) {
	root, err := Root()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "reimarchiv.db"), nil
}

// LogPath returns the path of the application log file.
func LogPath() (string, error) {
	root, err := Root()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "logs", "app.log"), nil
}
