package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
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
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// Add appends a snippet and persists the store.
func (s *jsonStore) Add(sn snippet.Snippet) (snippet.Snippet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, sn)
	if err := s.save(); err != nil {
		return snippet.Snippet{}, err
	}
	return sn, nil
}

// Get returns the snippet with the given id.
func (s *jsonStore) Get(id string) (snippet.Snippet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sn := range s.data {
		if sn.ID == id {
			return sn, nil
		}
	}
	return snippet.Snippet{}, ErrNotFound
}

// All returns every snippet, newest first.
func (s *jsonStore) All() ([]snippet.Snippet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]snippet.Snippet, len(s.data))
	copy(out, s.data)
	// Newest first, breaking ties on id so the order is stable across runs.
	sort.SliceStable(out, func(i, j int) bool {
		if !out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].CreatedAt.After(out[j].CreatedAt)
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

// Close implements Store; the JSON backend keeps nothing open between calls.
func (s *jsonStore) Close() error { return nil }

// Update replaces the stored snippet carrying the same id.
func (s *jsonStore) Update(sn snippet.Snippet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.data {
		if s.data[i].ID == sn.ID {
			s.data[i] = sn
			return s.save()
		}
	}
	return ErrNotFound
}

// Delete removes the snippet with the given id.
func (s *jsonStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.data {
		if s.data[i].ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return s.save()
		}
	}
	return ErrNotFound
}
