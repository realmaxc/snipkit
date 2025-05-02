// Package cli implements snipkit's command-line interface and subcommand
// dispatch.
package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

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

// overrides collects the global flags that sit in front of the subcommand and
// win over configuration and the environment.
var overrides struct {
	dbPath  string
	backend string
}

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
	rest, err := parseGlobalFlags(args)
	if err != nil {
		fmt.Fprintln(ctx.stderr, "snipkit:", err)
		return exitUsage
	}
	args = rest
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

// parseGlobalFlags consumes any leading --db/--backend options and returns the
// remaining command and arguments.
func parseGlobalFlags(args []string) ([]string, error) {
	i := 0
	for i < len(args) {
		arg := args[i]
		switch {
		case arg == "--db" || arg == "-db":
			if i+1 >= len(args) {
				return nil, errors.New("--db requires a path")
			}
			overrides.dbPath = args[i+1]
			i += 2
		case strings.HasPrefix(arg, "--db="):
			overrides.dbPath = strings.TrimPrefix(arg, "--db=")
			i++
		case arg == "--backend" || arg == "-backend":
			if i+1 >= len(args) {
				return nil, errors.New("--backend requires a name")
			}
			overrides.backend = args[i+1]
			i += 2
		case strings.HasPrefix(arg, "--backend="):
			overrides.backend = strings.TrimPrefix(arg, "--backend=")
			i++
		default:
			return args[i:], nil
		}
	}
	return nil, nil
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
	fmt.Fprintln(w, "  snipkit [--db PATH] [--backend json|sqlite] <command> [args]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	for _, name := range sortedCommandNames() {
		fmt.Fprintf(w, "  %-9s %s\n", name, commands[name].summary)
	}
	fmt.Fprintln(w, "  version   print the snipkit version")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Run 'snipkit <command> -h' for command-specific help.")
}

// resolveID turns a full id or an unambiguous id prefix into a full snippet id.
// TODO: fall back to fuzzy id matching when a prefix is ambiguous.
func resolveID(st store.Store, prefix string) (string, error) {
	if prefix == "" {
		return "", errors.New("empty id")
	}
	if _, err := st.Get(prefix); err == nil {
		return prefix, nil // exact hit, no scan needed
	}
	all, err := st.All()
	if err != nil {
		return "", err
	}
	var matches []string
	for _, s := range all {
		if strings.HasPrefix(s.ID, prefix) {
			matches = append(matches, s.ID)
		}
	}
	switch len(matches) {
	case 0:
		return "", store.ErrNotFound
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("id prefix %q is ambiguous (%d matches)", prefix, len(matches))
	}
}
