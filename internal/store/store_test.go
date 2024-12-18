package store

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/realmaxc/snipkit/internal/snippet"
)

func newBackends(t *testing.T) map[string]Store {
	t.Helper()
	dir := t.TempDir()
	js, err := openJSON(filepath.Join(dir, "snippets.json"))
	if err != nil {
		t.Fatalf("openJSON: %v", err)
	}
	sq, err := openSQLite(filepath.Join(dir, "snippets.db"))
	if err != nil {
		t.Fatalf("openSQLite: %v", err)
	}
	t.Cleanup(func() {
		_ = js.Close()
		_ = sq.Close()
	})
	return map[string]Store{"json": js, "sqlite": sq}
}

func sample(id, title string) snippet.Snippet {
	now := time.Now()
	return snippet.Snippet{
		ID: id, Title: title, Body: "echo " + title,
		Lang: "sh", Tags: []string{"demo"}, CreatedAt: now, UpdatedAt: now,
	}
}

func TestStoreConformance(t *testing.T) {
	for name, st := range newBackends(t) {
		t.Run(name, func(t *testing.T) {
			if _, err := st.Add(sample("aaa", "one")); err != nil {
				t.Fatalf("Add: %v", err)
			}
			if _, err := st.Add(sample("bbb", "two")); err != nil {
				t.Fatalf("Add: %v", err)
			}

			got, err := st.Get("aaa")
			if err != nil {
				t.Fatalf("Get: %v", err)
			}
			if got.Title != "one" {
				t.Fatalf("Get title = %q, want one", got.Title)
			}

			all, err := st.All()
			if err != nil {
				t.Fatalf("All: %v", err)
			}
			if len(all) != 2 {
				t.Fatalf("All len = %d, want 2", len(all))
			}

			upd := got
			upd.Title = "updated"
			if err := st.Update(upd); err != nil {
				t.Fatalf("Update: %v", err)
			}
			got, _ = st.Get("aaa")
			if got.Title != "updated" {
				t.Fatalf("after update title = %q", got.Title)
			}

			if err := st.Delete("aaa"); err != nil {
				t.Fatalf("Delete: %v", err)
			}
			if _, err := st.Get("aaa"); !errors.Is(err, ErrNotFound) {
				t.Fatalf("Get after delete = %v, want ErrNotFound", err)
			}
		})
	}
}
