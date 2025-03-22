package config

import (
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	t.Setenv("SNIPKIT_BACKEND", "")
	t.Setenv("SNIPKIT_DB", "")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Backend != BackendJSON {
		t.Errorf("Backend = %q, want json", cfg.Backend)
	}
	if filepath.Base(cfg.Path) != "snippets.json" {
		t.Errorf("Path = %q, want a snippets.json", cfg.Path)
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	t.Setenv("SNIPKIT_BACKEND", "sqlite")
	t.Setenv("SNIPKIT_DB", "")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Backend != BackendSQLite {
		t.Errorf("Backend = %q, want sqlite", cfg.Backend)
	}
	if filepath.Base(cfg.Path) != "snippets.db" {
		t.Errorf("Path = %q, want a snippets.db", cfg.Path)
	}

	t.Setenv("SNIPKIT_DB", "/tmp/custom.db")
	cfg, err = Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Path != "/tmp/custom.db" {
		t.Errorf("Path = %q, want the explicit SNIPKIT_DB", cfg.Path)
	}
}

func TestDefaultEditorFallsBack(t *testing.T) {
	t.Setenv("SNIPKIT_EDITOR", "")
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "nano")
	if got := defaultEditor(); got != "nano" {
		t.Errorf("defaultEditor() = %q, want nano", got)
	}
}
