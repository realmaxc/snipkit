package cli

import (
	"fmt"
	"sort"
	"text/tabwriter"
)

func init() {
	register(command{name: "tags", summary: "list tags and their counts", run: runTags})
}

func runTags(ctx *appContext, args []string) error {
	_ = args
	st, _, err := openStore()
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	all, err := st.All()
	if err != nil {
		return err
	}
	counts := map[string]int{}
	for _, s := range all {
		for _, t := range s.Tags {
			counts[t]++
		}
	}
	if len(counts) == 0 {
		fmt.Fprintln(ctx.stdout, "no tags yet")
		return nil
	}
	names := make([]string, 0, len(counts))
	for n := range counts {
		names = append(names, n)
	}
	sort.Strings(names)

	tw := tabwriter.NewWriter(ctx.stdout, 0, 2, 2, ' ', 0)
	for _, n := range names {
		fmt.Fprintf(tw, "%s\t%d\n", n, counts[n])
	}
	_ = tw.Flush()
	return nil
}
