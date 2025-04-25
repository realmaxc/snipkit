package store

import (
	"errors"
	"os"
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

func TestJSONStoreStableOrderForEqualTimes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "s.json")
	st, err := openJSON(path)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	for _, id := range []string{"ccc", "aaa", "bbb"} {
		if _, err := st.Add(snippet.Snippet{ID: id, Body: "x", CreatedAt: now, UpdatedAt: now}); err != nil {
			t.Fatal(err)
		}
	}
	all, err := st.All()
	if err != nil {
		t.Fatal(err)
	}
	got := []string{all[0].ID, all[1].ID, all[2].ID}
	want := []string{"aaa", "bbb", "ccc"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order = %v, want %v (ties broken by id)", got, want)
		}
	}
}

func TestJSONStoreToleratesEmptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "s.json")
	if err := os.WriteFile(path, []byte("   \n"), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := openJSON(path)
	if err != nil {
		t.Fatalf("openJSON on empty file: %v", err)
	}
	all, err := st.All()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 0 {
		t.Fatalf("All() = %d snippets, want 0", len(all))
	}
}
