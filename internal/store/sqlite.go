package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/realmaxc/snipkit/internal/snippet"
	_ "modernc.org/sqlite"
)

// sqlStore persists snippets in a local SQLite database via the pure-Go
// modernc.org/sqlite driver (no cgo required).
type sqlStore struct {
	db *sql.DB
}

func openSQLite(path string) (*sqlStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schemaSQL); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &sqlStore{db: db}, nil
}

const selectColumns = `id, title, body, lang, tags, created_at, updated_at`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanSnippet(r rowScanner) (snippet.Snippet, error) {
	var (
		sn                     snippet.Snippet
		tags, created, updated string
	)
	if err := r.Scan(&sn.ID, &sn.Title, &sn.Body, &sn.Lang, &tags, &created, &updated); err != nil {
		return snippet.Snippet{}, err
	}
	if tags != "" {
		if err := json.Unmarshal([]byte(tags), &sn.Tags); err != nil {
			return snippet.Snippet{}, err
		}
	}
	sn.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	sn.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return sn, nil
}

// Add inserts a new snippet.
func (s *sqlStore) Add(sn snippet.Snippet) (snippet.Snippet, error) {
	tags, err := json.Marshal(sn.Tags)
	if err != nil {
		return snippet.Snippet{}, err
	}
	_, err = s.db.Exec(
		`INSERT INTO snippets (id, title, body, lang, tags, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		sn.ID, sn.Title, sn.Body, sn.Lang, string(tags),
		sn.CreatedAt.Format(time.RFC3339Nano), sn.UpdatedAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		return snippet.Snippet{}, err
	}
	return sn, nil
}

// Get returns the snippet with the given id.
func (s *sqlStore) Get(id string) (snippet.Snippet, error) {
	row := s.db.QueryRow(`SELECT `+selectColumns+` FROM snippets WHERE id = ?`, id)
	sn, err := scanSnippet(row)
	if errors.Is(err, sql.ErrNoRows) {
		return snippet.Snippet{}, ErrNotFound
	}
	return sn, err
}

// All returns every snippet, newest first.
func (s *sqlStore) All() ([]snippet.Snippet, error) {
	rows, err := s.db.Query(`SELECT ` + selectColumns + ` FROM snippets ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []snippet.Snippet
	for rows.Next() {
		sn, err := scanSnippet(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sn)
	}
	return out, rows.Err()
}

// Close closes the underlying database handle.
func (s *sqlStore) Close() error { return s.db.Close() }

// Update replaces the stored snippet carrying the same id.
func (s *sqlStore) Update(sn snippet.Snippet) error {
	tags, err := json.Marshal(sn.Tags)
	if err != nil {
		return err
	}
	res, err := s.db.Exec(
		`UPDATE snippets SET title = ?, body = ?, lang = ?, tags = ?, updated_at = ? WHERE id = ?`,
		sn.Title, sn.Body, sn.Lang, string(tags), sn.UpdatedAt.Format(time.RFC3339Nano), sn.ID,
	)
	if err != nil {
		return err
	}
	return affectedOrNotFound(res)
}

// Delete removes the snippet with the given id.
func (s *sqlStore) Delete(id string) error {
	res, err := s.db.Exec(`DELETE FROM snippets WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return affectedOrNotFound(res)
}

func affectedOrNotFound(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
