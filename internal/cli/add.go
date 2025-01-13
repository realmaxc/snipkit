package cli

import (
	"flag"
	"fmt"
	"io"
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
	body := fs.String("body", "", "snippet body; read from stdin when omitted")
	var tags multiFlag
	fs.Var(&tags, "tag", "tag to attach (repeatable)")
	fs.Usage = usageFor(ctx.stderr, "snipkit add [-title T] [-lang L] [-tag t]... [-body TEXT]", fs)
	if err := fs.Parse(args); err != nil {
		return errSilent
	}

	text := *body
	if text == "" {
		b, err := io.ReadAll(ctx.stdin)
		if err != nil {
			return err
		}
		text = string(b)
	}

	st, _, err := openStore()
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

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
