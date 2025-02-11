package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/realmaxc/snipkit/internal/fuzzy"
	"github.com/realmaxc/snipkit/internal/snippet"
)

func init() {
	register(command{name: "search", summary: "fuzzy-search snippets", run: runSearch})
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
	matches := rankSnippets(all, query)
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

// rankSnippets returns snippets that fuzzy-match query, best first. An empty
// query keeps the store's own newest-first ordering.
func rankSnippets(in []snippet.Snippet, query string) []snippet.Snippet {
	if strings.TrimSpace(query) == "" {
		return in
	}
	cands := make([]fuzzy.Candidate, len(in))
	for i, s := range in {
		cands[i] = fuzzy.Candidate{Index: i, Text: haystack(s)}
	}
	ranked := fuzzy.Rank(query, cands)
	out := make([]snippet.Snippet, 0, len(ranked))
	for _, r := range ranked {
		out = append(out, in[r.Index])
	}
	return out
}

func haystack(s snippet.Snippet) string {
	return s.Title + " " + s.Body + " " + strings.Join(s.Tags, " ")
}
