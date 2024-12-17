package store

// schemaSQL is executed when a SQLite store is opened. It is idempotent so it
// doubles as the (very small) migration for a brand-new database.
const schemaSQL = `
CREATE TABLE IF NOT EXISTS snippets (
	id         TEXT PRIMARY KEY,
	title      TEXT NOT NULL,
	body       TEXT NOT NULL,
	lang       TEXT NOT NULL DEFAULT '',
	tags       TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_snippets_created_at ON snippets (created_at);
`
