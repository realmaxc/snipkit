// Package store defines the snippet persistence contract and its backends.
package store

import (
	"errors"

	"github.com/realmaxc/snipkit/internal/snippet"
)

// ErrNotFound is returned when a snippet id does not exist in the store.
var ErrNotFound = errors.New("snippet not found")

// Store is the persistence contract every backend implements. Implementations
// only need to be safe for sequential use by a single CLI process; snipkit
// never shares a store across concurrent writers.
type Store interface {
	// Add inserts a new snippet and returns it with store-populated fields.
	Add(s snippet.Snippet) (snippet.Snippet, error)
	// Get returns the snippet with the given id, or ErrNotFound.
	Get(id string) (snippet.Snippet, error)
	// All returns every snippet, newest first.
	All() ([]snippet.Snippet, error)
	// Update replaces the stored snippet carrying the same id.
	Update(s snippet.Snippet) error
	// Delete removes the snippet with the given id, or returns ErrNotFound.
	Delete(id string) error
	// Close flushes and releases any resources held by the store.
	Close() error
}

// Open opens the named backend at path. Supported backends are "json" and
// "sqlite".
func Open(backend, path string) (Store, error) {
	switch backend {
	case "json":
		return openJSON(path)
	case "sqlite":
		return openSQLite(path)
	default:
		return nil, errUnsupportedBackend(backend)
	}
}

func errUnsupportedBackend(name string) error {
	return &unsupportedBackendError{name: name}
}

type unsupportedBackendError struct{ name string }

func (e *unsupportedBackendError) Error() string {
	return "unsupported store backend: " + e.name
}
