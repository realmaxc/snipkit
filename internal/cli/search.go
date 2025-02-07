package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/realmaxc/snipkit/internal/snippet"
)

func init() {
	register(command{name: "search", summary: "search snippets by text", run: runSearch})
}

func runSearch(ctx *appContext, args []string) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	limit := fs.Int("n", 20, "maximum number of results")
	fs.Usage = usageFor(ctx.stderr, "snipkit search [-n N] <query>", fs)
	if err := fs.Parse(args); err != nil {
		return errSilent
	}
	query := strings.Join(fs.Args(), " ")

	st, _, err := openStore()
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	all, err := st.All()
	if err != nil {
		return err
	}
	matches := substringSearch(all, query)
	if *limit > 0 && len(matches) > *limit {
		matches = matches[:*limit]
	}
	if len(matches) == 0 {
		fmt.Fprintln(ctx.stderr, "no matches")
		return nil
	}
	renderTable(ctx.stdout, matches)
	return nil
}

func substringSearch(in []snippet.Snippet, query string) []snippet.Snippet {
	if query == "" {
		return in
	}
	q := strings.ToLower(query)
	out := make([]snippet.Snippet, 0, len(in))
	for _, s := range in {
		if strings.Contains(strings.ToLower(haystack(s)), q) {
			out = append(out, s)
		}
	}
	return out
}

func haystack(s snippet.Snippet) string {
	return s.Title + " " + s.Body + " " + strings.Join(s.Tags, " ")
}
