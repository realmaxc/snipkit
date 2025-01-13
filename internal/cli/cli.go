// Package cli implements snipkit's command-line interface and subcommand
// dispatch.
package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/realmaxc/snipkit/internal/store"
	"github.com/realmaxc/snipkit/internal/version"
)

// Process exit codes returned by Run.
const (
	exitOK       = 0
	exitError    = 1
	exitUsage    = 2
	exitNotFound = 3
)

// errSilent is returned by a subcommand that has already written its own error
// and usage output; the dispatcher then exits with a usage code and prints
// nothing further.
var errSilent = errors.New("usage error already reported")

// command is a single snipkit subcommand.
type command struct {
	name    string
	summary string
	run     func(ctx *appContext, args []string) error
}

// appContext carries shared IO into subcommands so they can be exercised in
// tests without touching the real terminal.
type appContext struct {
	stdout io.Writer
	stderr io.Writer
	stdin  io.Reader
}

var commands = map[string]command{}

func register(c command) { commands[c.name] = c }

func sortedCommandNames() []string {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Run dispatches a single invocation and returns a process exit code.
func Run(args []string) int {
	ctx := &appContext{stdout: os.Stdout, stderr: os.Stderr, stdin: os.Stdin}
	return run(ctx, args)
}

func run(ctx *appContext, args []string) int {
	if len(args) == 0 {
		usage(ctx.stderr)
		return exitUsage
	}
	switch args[0] {
	case "-h", "--help", "help":
		usage(ctx.stdout)
		return exitOK
	case "-v", "--version", "version":
		fmt.Fprintln(ctx.stdout, version.String())
		return exitOK
	}
	cmd, ok := commands[args[0]]
	if !ok {
		fmt.Fprintf(ctx.stderr, "snipkit: unknown command %q\n", args[0])
		usage(ctx.stderr)
		return exitUsage
	}
	if err := cmd.run(ctx, args[1:]); err != nil {
		return reportError(ctx.stderr, err)
	}
	return exitOK
}

func reportError(w io.Writer, err error) int {
	if errors.Is(err, errSilent) {
		return exitUsage
	}
	fmt.Fprintln(w, "snipkit:", err)
	if errors.Is(err, store.ErrNotFound) {
		return exitNotFound
	}
	return exitError
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "snipkit — a fast terminal snippet manager")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  snipkit <command> [args]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	for _, name := range sortedCommandNames() {
		fmt.Fprintf(w, "  %-9s %s\n", name, commands[name].summary)
	}
	fmt.Fprintln(w, "  version   print the snipkit version")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Run 'snipkit <command> -h' for command-specific help.")
}
