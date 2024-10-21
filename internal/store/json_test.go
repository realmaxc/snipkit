package store

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/realmaxc/snipkit/internal/snippet"
)

func TestJSONStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "snippets.json")
	st, err := openJSON(path)
	if err != nil {
		t.Fatalf("openJSON: %v", err)
	}
	now := time.Now()
	in := snippet.Snippet{ID: "abc", Title: "greet", Body: "echo hi", Lang: "sh", CreatedAt: now, UpdatedAt: now}
	if _, err := st.Add(in); err != nil {
		t.Fatalf("Add: %v", err)
	}

	// Reopen to prove persistence.
	st2, err := openJSON(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	got, err := st2.Get("abc")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Body != "echo hi" {
		t.Fatalf("Body = %q, want %q", got.Body, "echo hi")
	}
	if _, err := st2.Get("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get(missing) = %v, want ErrNotFound", err)
	}
}
