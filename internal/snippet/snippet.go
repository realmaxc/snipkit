// Package snippet defines the core snippet value type shared across snipkit.
package snippet

import "time"

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
