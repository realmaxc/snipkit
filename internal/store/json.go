package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/realmaxc/snipkit/internal/snippet"
)

// jsonStore persists snippets as a single JSON document on disk. The whole
// document is held in memory, which is plenty for the human-scale collections
// snipkit targets.
type jsonStore struct {
	path string
	mu   sync.Mutex
	data []snippet.Snippet
}

type jsonDocument struct {
	Version  int               `json:"version"`
	Snippets []snippet.Snippet `json:"snippets"`
}

func openJSON(path string) (*jsonStore, error) {
	s := &jsonStore{path: path}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *jsonStore) load() error {
	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		s.data = nil
		return nil
	}
	if err != nil {
		return err
	}
	var doc jsonDocument
	if err := json.Unmarshal(b, &doc); err != nil {
		return err
	}
	s.data = doc.Snippets
	return nil
}

func (s *jsonStore) save() error {
	doc := jsonDocument{Version: 1, Snippets: s.data}
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	// Write to a sibling temp file and rename so a crash never leaves a
	// half-written store behind.
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// Close implements Store; the JSON backend keeps nothing open between calls.
func (s *jsonStore) Close() error { return nil }
