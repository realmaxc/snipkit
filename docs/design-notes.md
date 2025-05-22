# Design notes

A few decisions worth writing down, mostly so I remember why later.

## Hand-rolled command dispatch

snipkit uses a tiny registry (`register`/`command`) instead of a CLI framework.
The command set is small and stable, the flag needs are modest, and avoiding a
framework keeps startup instant and the dependency graph almost empty. Each
subcommand registers itself in an `init`, so adding a command is one new file.

## Two backends, one interface

The JSON backend loads the whole document into memory. For the human-scale
collections snipkit targets (hundreds, maybe low thousands of snippets) that is
simpler and faster than round-tripping a database, and the file is easy to read,
diff, and back up. SQLite exists for people who want a real index or already
keep everything in databases. A shared conformance suite guarantees the two
behave identically.

## Short ids

Ids are five random bytes rendered in a Crockford-style base32 alphabet (no
`i`, `l`, `o`, `u`), giving eight unambiguous characters. Commands accept any
unique prefix, so `k3f` usually works where `k3f9a1b2` would.

## Fuzzy scoring

The ranker is a subsequence matcher with a small scoring pass. Matches score
higher when they are contiguous and when they land on word boundaries, which is
what makes acronym queries like `gco` → `git checkout` feel natural. The weights
are deliberately simple; there's a `TODO` to revisit them with real usage data.

## No network, ever

There is no update check, no analytics, no sync. snipkit only touches the store
file and, for the clipboard and picker, well-known local utilities.
