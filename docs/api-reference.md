# CLI reference

```
snipkit [--db PATH] [--backend json|sqlite] <command> [args]
```

Global options may appear before the command:

- `--db PATH` — override the store path (also `SNIPKIT_DB`).
- `--backend json|sqlite` — override the backend (also `SNIPKIT_BACKEND`).

## Commands

### `add`
Add a snippet. The body comes from stdin, or from `$EDITOR` when stdin is a
terminal.

```
snipkit add [-title T] [-lang L] [-tag t]... [-body TEXT]
```

### `list`
List stored snippets, newest first.

```
snipkit list [-tag t] [-lang l]
```

### `search`
Fuzzy-search across title, body, and tags. Best matches first.

```
snipkit search [-n N] <query>
```

### `get`
Print a snippet's body to stdout. Accepts a full id or a unique prefix.

```
snipkit get <id|prefix>
```

### `copy`
Copy a snippet's body to the system clipboard.

```
snipkit copy <id|prefix>
```

### `pick`
Interactively pick a snippet (via `fzf` when available) and print its body.

```
snipkit pick [-c]
```

`-c` also copies the selection to the clipboard.

### `edit`
Open a snippet's body in `$EDITOR` and save the result.

```
snipkit edit <id|prefix>
```

### `rm`
Remove one or more snippets.

```
snipkit rm <id|prefix>...
```

### `tags`
List every tag with a count of how many snippets carry it.

### `export` / `import`
Move snippets in and out as a JSON array.

```
snipkit export [-o file.json]
snipkit import [-f file.json]
```

### `version`
Print the snipkit version banner.

## Exit codes

| Code | Meaning |
| ---- | ------- |
| 0 | success |
| 1 | a runtime error |
| 2 | a usage error (bad flags or unknown command) |
| 3 | the requested snippet was not found |
