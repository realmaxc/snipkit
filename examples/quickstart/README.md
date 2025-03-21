# Quickstart

Seed a fresh store with a handful of everyday one-liners, then search them.

```console
$ snipkit import -f examples/quickstart/snippets.json
imported 4 snippets

$ snipkit search gitgraph
ID        LANG  TAGS  TITLE
xxxxxxxx  sh    git   show git graph

$ snipkit list -tag net
ID        LANG  TAGS         TITLE
xxxxxxxx  sh    net,process  kill process on port
```

Point snipkit at a scratch store so you don't touch your real one:

```console
$ snipkit --db /tmp/demo.json import -f examples/quickstart/snippets.json
```
