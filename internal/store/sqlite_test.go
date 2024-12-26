package store

import (
	"path/filepath"
	"testing"
)

func TestSQLitePersistsAcrossReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "snips.db")
	st, err := openSQLite(path)
	if err != nil {
		t.Fatalf("openSQLite: %v", err)
	}
	if _, err := st.Add(sample("keep", "persist me")); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	st2, err := openSQLite(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer func() { _ = st2.Close() }()
	got, err := st2.Get("keep")
	if err != nil {
		t.Fatalf("Get after reopen: %v", err)
	}
	if got.Title != "persist me" {
		t.Fatalf("Title = %q, want 'persist me'", got.Title)
	}
}
