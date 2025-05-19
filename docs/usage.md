# Usage guide

## Adding snippets

Pipe a command in, or let snipkit open your editor:

```console
$ echo 'docker system prune -af' | snipkit add -title "docker cleanup" -tag docker
$ snipkit add -title "scratch"          # opens $EDITOR for the body
```

Flags: `-title`, `-lang`, `-tag` (repeatable), `-body`.

If you don't pass `-title`, snipkit uses the first line of the body.

## Finding snippets

```console
$ snipkit list                 # everything, newest first
$ snipkit list -tag git        # filter by tag
$ snipkit list -lang sh        # filter by language
$ snipkit search dkclean       # fuzzy search across title, body, and tags
$ snipkit search -n 5 gco      # cap the number of results
```

## Using a snippet

```console
$ snipkit get k3f9a1b2         # print the body to stdout
$ snipkit copy k3f             # copy to the clipboard (id prefixes are fine)
$ snipkit pick                 # interactively choose one
$ snipkit pick -c              # ...and copy the choice to the clipboard
```

`get` prints the raw body, which makes it easy to compose:

```console
$ eval "$(snipkit get k3f)"
```

## Editing and removing

```console
$ snipkit edit k3f             # open the body in $EDITOR
$ snipkit rm k3f 9a1 7bd       # remove several at once
```

## Moving snippets around

```console
$ snipkit export -o backup.json
$ snipkit import -f backup.json
$ snipkit tags                 # list tags with counts
```
