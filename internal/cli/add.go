package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/realmaxc/snipkit/internal/config"
	"github.com/realmaxc/snipkit/internal/snippet"
	"github.com/realmaxc/snipkit/internal/store"
)

func init() {
	register(command{name: "add", summary: "add a new snippet", run: runAdd})
}

// multiFlag collects a repeatable string flag into a slice.
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }

func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

func usageFor(w io.Writer, synopsis string, fs *flag.FlagSet) func() {
	return func() {
		fmt.Fprintln(w, "Usage:", synopsis)
		fs.PrintDefaults()
	}
}

// openStore resolves configuration and opens the configured store backend.
func openStore() (store.Store, config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, config.Config{}, err
	}
	if overrides.backend != "" {
		cfg.Backend = config.Backend(overrides.backend)
	}
	if overrides.dbPath != "" {
		cfg.Path = overrides.dbPath
	}
	st, err := store.Open(string(cfg.Backend), cfg.Path)
	if err != nil {
		return nil, config.Config{}, err
	}
	return st, cfg, nil
}

func runAdd(ctx *appContext, args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	title := fs.String("title", "", "short title (defaults to the first body line)")
	lang := fs.String("lang", "", "language or shell, e.g. sh, go, py")
	body := fs.String("body", "", "snippet body; read from stdin or $EDITOR when omitted")
	var tags multiFlag
	fs.Var(&tags, "tag", "tag to attach (repeatable)")
	fs.Usage = usageFor(ctx.stderr, "snipkit add [-title T] [-lang L] [-tag t]... [-body TEXT]", fs)
	if err := fs.Parse(args); err != nil {
		return errSilent
	}

	st, cfg, err := openStore()
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	text := *body
	if text == "" {
		text, err = readBody(ctx, cfg.Editor)
		if err != nil {
			return err
		}
	}

	now := time.Now()
	sn := snippet.Snippet{
		ID:        snippet.NewID(),
		Title:     *title,
		Body:      text,
		Lang:      *lang,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
	sn.Normalize()
	if err := sn.Validate(); err != nil {
		return err
	}
	saved, err := st.Add(sn)
	if err != nil {
		return err
	}
	fmt.Fprintf(ctx.stdout, "added %s  %s\n", saved.ID, saved.Title)
	return nil
}

// readBody reads a snippet body from stdin when input is piped, otherwise it
// opens the configured editor on a scratch file.
func readBody(ctx *appContext, editor string) (string, error) {
	if f, ok := ctx.stdin.(*os.File); ok && isTerminal(f) {
		return editBuffer(editor, "")
	}
	b, err := io.ReadAll(ctx.stdin)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// editBuffer opens editor on a temp file seeded with initial and returns the
// edited contents.
func editBuffer(editor, initial string) (string, error) {
	tmp, err := os.CreateTemp("", "snipkit-*.txt")
	if err != nil {
		return "", err
	}
	defer func() { _ = os.Remove(tmp.Name()) }()
	if initial != "" {
		if _, err := tmp.WriteString(initial); err != nil {
			return "", err
		}
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	if err := runEditor(editor, tmp.Name()); err != nil {
		return "", err
	}
	out, err := os.ReadFile(tmp.Name())
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func runEditor(editor, path string) error {
	parts := strings.Fields(editor)
	if len(parts) == 0 {
		return errors.New("no editor configured; set $EDITOR")
	}
	cmd := exec.Command(parts[0], append(parts[1:], path)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}
