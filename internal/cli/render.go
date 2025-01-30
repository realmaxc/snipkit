package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/realmaxc/snipkit/internal/snippet"
)

// renderTable prints snippets as an aligned id / lang / tags / title table.
func renderTable(w io.Writer, snippets []snippet.Snippet) {
	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tLANG\tTAGS\tTITLE")
	for _, s := range snippets {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", s.ID, dash(s.Lang), dash(strings.Join(s.Tags, ",")), truncate(s.Title, 60))
	}
	_ = tw.Flush()
}

func dash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n <= 1 {
		return string(r[:n])
	}
	return string(r[:n-1]) + "…"
}
