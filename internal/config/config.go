package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Backend identifies which store implementation snipkit should use.
type Backend string

const (
	// BackendJSON keeps snippets in a single JSON document.
	BackendJSON Backend = "json"
	// BackendSQLite keeps snippets in a local SQLite database.
	BackendSQLite Backend = "sqlite"
)

// Config holds the resolved runtime settings for a snipkit invocation.
type Config struct {
	Backend Backend // which store implementation to open
	Path    string  // filesystem path to the store
	Editor  string  // command used for interactive edits
}

// Load resolves configuration from built-in defaults. Environment overrides
// are layered on top in a later change.
func Load() (Config, error) {
	base, err := dataDir()
	if err != nil {
		return Config{}, err
	}
	cfg := Config{
		Backend: BackendJSON,
		Path:    filepath.Join(base, "snippets.json"),
		Editor:  defaultEditor(),
	}
	return cfg, nil
}

func defaultEditor() string {
	for _, key := range []string{"SNIPKIT_EDITOR", "VISUAL", "EDITOR"} {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
	}
	return "vi"
}
