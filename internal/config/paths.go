// Package config resolves snipkit's runtime settings and data locations.
package config

import (
	"os"
	"path/filepath"
)

// dataDir returns the base directory snipkit stores its data in, honouring
// XDG_DATA_HOME and otherwise falling back to ~/.local/share/snipkit.
func dataDir() (string, error) {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return filepath.Join(dir, "snipkit"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "snipkit"), nil
}
