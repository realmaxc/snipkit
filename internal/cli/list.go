package cli

import (
	"flag"
	"fmt"

	"github.com/realmaxc/snipkit/internal/snippet"
)

func init() {
	register(command{name: "list", summary: "list stored snippets", run: runList})
}

func runList(ctx *appContext, args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	tag := fs.String("tag", "", "only show snippets carrying this tag")
	lang := fs.String("lang", "", "only show snippets in this language")
	fs.Usage = usageFor(ctx.stderr, "snipkit list [-tag t] [-lang l]", fs)
	if err := fs.Parse(args); err != nil {
		return errSilent
	}

	st, _, err := openStore()
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	all, err := st.All()
	if err != nil {
		return err
	}
	rows := filterSnippets(all, *tag, *lang)
	if len(rows) == 0 {
		fmt.Fprintln(ctx.stdout, "no snippets yet — add one with 'snipkit add'")
		return nil
	}
	renderTable(ctx.stdout, rows)
	return nil
}

func filterSnippets(in []snippet.Snippet, tag, lang string) []snippet.Snippet {
	if tag == "" && lang == "" {
		return in
	}
	out := make([]snippet.Snippet, 0, len(in))
	for _, s := range in {
		if tag != "" && !s.HasTag(tag) {
			continue
		}
		if lang != "" && s.Lang != lang {
			continue
		}
		out = append(out, s)
	}
	return out
}
