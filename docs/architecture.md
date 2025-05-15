# Architecture

snipkit is a small, single-binary CLI. The packages are deliberately narrow so
each one is easy to test in isolation.

```
cmd/snipkit           process entrypoint; calls cli.Run(os.Args[1:])
internal/cli          argument parsing, subcommand dispatch, output
internal/store        Store interface + JSON and SQLite backends
internal/snippet      the Snippet value type, validation, id generation
internal/fuzzy        subsequence matcher and ranking
internal/clip         clipboard integration (shells out per-OS)
internal/config       resolves the store path, backend, and editor
internal/version      build/version metadata (set via -ldflags)
```

## Request flow

A typical `snipkit search gco` invocation flows like this:

1. `cli.Run` strips any global `--db`/`--backend` flags and dispatches to the
   `search` subcommand.
2. The subcommand loads configuration (`config.Load`), applies the global
   overrides, and opens the configured store (`store.Open`).
3. It reads every snippet (`Store.All`) and ranks them with `fuzzy.Rank`.
4. Results are rendered as an aligned table to stdout.

## The store abstraction

Every backend implements the same six-method `store.Store` interface, and a
single conformance test suite runs against all of them. That keeps the JSON and
SQLite backends behaviourally identical, so switching `--backend` never changes
what you get back — only how it's stored on disk.

## Why so few dependencies

The only third-party dependency is the pure-Go
[`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) driver, and it is
only pulled in for the SQLite backend. Everything else — argument parsing, fuzzy
matching, table rendering — is the standard library. That keeps builds fast and
the binary self-contained.
