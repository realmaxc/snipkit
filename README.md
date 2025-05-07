# snipkit

A fast terminal snippet manager. Save shell one-liners and code snippets, then
fuzzy-search and paste them back without ever leaving the terminal. Everything
lives in a local file — no accounts, no sync server, no telemetry.

## Why

I kept re-deriving the same `ffmpeg`, `jq`, and `git` incantations, or digging
through shell history hoping I'd typed one recently. snipkit is the little tool
I wanted: type three letters, get the command, paste it. That's it.

## Features

- **Fuzzy search** with an acronym-friendly ranker — `gco` finds `git checkout`.
- **Two local backends**: a human-readable JSON document, or SQLite when you
  want indexing. Switch with a flag or an environment variable.
- **Clipboard integration** on macOS (`pbcopy`), Wayland (`wl-copy`), X11
  (`xclip`/`xsel`), and Windows (`clip`).
- **Interactive picker** that uses [`fzf`](https://github.com/junegunn/fzf) when
  it's available and falls back to a numbered prompt otherwise.
- **Shell key bindings** for bash, zsh, and fish — bind a key and drop a snippet
  straight onto the command line.
- **Import / export** as plain JSON, so your snippets are yours to move around.
- **Zero network.** snipkit never opens a socket.

## Install

```console
$ go install github.com/realmaxc/snipkit/cmd/snipkit@latest
```

Or build from source:

```console
$ git clone https://github.com/realmaxc/snipkit
$ cd snipkit
$ make build      # produces ./snipkit
```

## Quick start

```console
$ echo 'git log --oneline --graph --decorate --all' | \
    snipkit add -title "git graph" -lang sh -tag git

$ snipkit list
ID        LANG  TAGS  TITLE
k3f9a1b2  sh    git   git graph

$ snipkit search gitgraph
ID        LANG  TAGS  TITLE
k3f9a1b2  sh    git   git graph

$ snipkit copy k3f      # id prefixes work too
copied k3f9a1b2 to clipboard
```

## Storage & configuration

By default snipkit keeps a JSON store at
`$XDG_DATA_HOME/snipkit/snippets.json` (falling back to
`~/.local/share/snipkit/snippets.json`).

| Setting | Env var | Flag | Default |
| ------- | ------- | ---- | ------- |
| Store path | `SNIPKIT_DB` | `--db` | `.../snipkit/snippets.json` |
| Backend | `SNIPKIT_BACKEND` | `--backend` | `json` |
| Editor | `SNIPKIT_EDITOR`, `VISUAL`, `EDITOR` | — | `vi` |

```console
$ snipkit --backend sqlite --db ~/snips.db list
```

## Shell integration

Bind <kbd>Ctrl</kbd>+<kbd>S</kbd> to the interactive picker — see
[`examples/shell/`](examples/shell/). For bash:

```bash
source /path/to/snipkit/examples/shell/bash-widget.sh
```

## Documentation

- [Usage guide](docs/usage.md)
- [CLI reference](docs/api-reference.md)
- [Architecture](docs/architecture.md)
- [Design notes](docs/design-notes.md)

## License

[MIT](LICENSE) © Max Collins
