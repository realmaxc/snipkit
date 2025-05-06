package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/realmaxc/snipkit/internal/snippet"
)

func init() {
	register(command{name: "import", summary: "import snippets from JSON", run: runImport})
}

func runImport(ctx *appContext, args []string) error {
	fs := flag.NewFlagSet("import", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	in := fs.String("f", "", "read from this file instead of stdin")
	fs.Usage = usageFor(ctx.stderr, "snipkit import [-f file.json]", fs)
	if err := fs.Parse(args); err != nil {
		return errSilent
	}

	reader := ctx.stdin
	if *in != "" {
		f, err := os.Open(*in)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()
		reader = f
	}
	raw, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	var incoming []snippet.Snippet
	if err := json.Unmarshal(raw, &incoming); err != nil {
		return fmt.Errorf("parsing import: %w", err)
	}

	st, _, err := openStore()
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	added := 0
	for _, sn := range incoming {
		if sn.ID == "" {
			sn.ID = snippet.NewID()
		}
		now := time.Now()
		if sn.CreatedAt.IsZero() {
			sn.CreatedAt = now
		}
		sn.UpdatedAt = now
		sn.Normalize()
		if err := sn.Validate(); err != nil {
			fmt.Fprintf(ctx.stderr, "skipping %q: %v\n", sn.Title, err)
			continue
		}
		if _, err := st.Add(sn); err != nil {
			return err
		}
		added++
	}
	fmt.Fprintf(ctx.stderr, "imported %d snippets\n", added)
	return nil
}
