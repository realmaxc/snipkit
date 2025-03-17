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

// Load resolves configuration from built-in defaults layered with the
// SNIPKIT_BACKEND, SNIPKIT_DB and editor environment variables.
func Load() (Config, error) {
	base, err := dataDir()
	if err != nil {
		return Config{}, err
	}
	cfg := Config{Backend: BackendJSON, Editor: defaultEditor()}
	if b := strings.TrimSpace(os.Getenv("SNIPKIT_BACKEND")); b != "" {
		cfg.Backend = Backend(b)
	}
	if p := strings.TrimSpace(os.Getenv("SNIPKIT_DB")); p != "" {
		cfg.Path = p
	} else {
		cfg.Path = filepath.Join(base, defaultFile(cfg.Backend))
	}
	return cfg, nil
}

func defaultFile(b Backend) string {
	if b == BackendSQLite {
		return "snippets.db"
	}
	return "snippets.json"
}

func defaultEditor() string {
	for _, key := range []string{"SNIPKIT_EDITOR", "VISUAL", "EDITOR"} {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
	}
	return "vi"
}
