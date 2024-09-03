// Package snippet defines the core snippet value type shared across snipkit.
package snippet

import (
	"errors"
	"sort"
	"strings"
	"time"
)

// ErrEmptyBody is returned when a snippet has no body text to store.
var ErrEmptyBody = errors.New("snippet body is empty")

// Snippet is a single stored piece of text — a shell one-liner, a command, or
// a short block of code — together with the metadata snipkit indexes on.
type Snippet struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Lang      string    `json:"lang,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Normalize cleans incidental whitespace, lower-cases the language, and
// de-duplicates tags. When no title is set it derives one from the first line
// of the body so every snippet has something searchable to show.
func (s *Snippet) Normalize() {
	s.Title = strings.TrimSpace(s.Title)
	s.Body = strings.TrimRight(s.Body, "\n")
	s.Lang = strings.ToLower(strings.TrimSpace(s.Lang))
	s.Tags = normalizeTags(s.Tags)
	if s.Title == "" {
		s.Title = firstLine(s.Body)
	}
}

// Validate reports whether the snippet is well-formed enough to be stored.
func (s Snippet) Validate() error {
	if strings.TrimSpace(s.Body) == "" {
		return ErrEmptyBody
	}
	return nil
}

// HasTag reports whether the snippet carries the given tag (case-insensitive).
func (s Snippet) HasTag(tag string) bool {
	tag = strings.ToLower(strings.TrimSpace(tag))
	for _, t := range s.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func normalizeTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		t = strings.ToLower(strings.TrimSpace(t))
		if t == "" {
			continue
		}
		if _, dup := seen[t]; dup {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return s
}
