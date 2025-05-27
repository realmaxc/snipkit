# Contributing

Thanks for taking a look! snipkit is small and I'd like to keep it that way, but
bug reports, fixes, and focused features are very welcome.

## Getting set up

You'll need Go 1.22 or newer.

```console
$ git clone https://github.com/realmaxc/snipkit
$ cd snipkit
$ make build
$ make test
```

## Before you open a PR

- `make test` passes (it runs with the race detector).
- `make lint` is clean (`golangci-lint`).
- `gofmt` is happy — `make fmt` will sort it out.
- New behaviour comes with a test. The store has a shared conformance suite;
  new backends should join it rather than duplicate it.

## Style

- Keep commands small: one file per subcommand, registered in an `init`.
- Errors are lower-case and don't end with punctuation.
- Prefer the standard library. New third-party dependencies need a good reason.

## Reporting bugs

Open an issue with your `snipkit version`, OS, and the exact command you ran.
There's a template to make that easy.
