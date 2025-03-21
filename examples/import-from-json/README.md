# Importing from JSON

`snipkit import` reads a JSON array of snippets from stdin or a file. Only
`title` and `body` are required — ids and timestamps are filled in for you.

```console
$ snipkit import -f examples/import-from-json/sample.json
imported 2 snippets

# or pipe it in
$ cat examples/import-from-json/sample.json | snipkit import
```

`snipkit export` produces the same shape, so exports round-trip back through
`import`.
