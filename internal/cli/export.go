package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func init() {
	register(command{name: "export", summary: "export snippets as JSON", run: runExport})
}

func runExport(ctx *appContext, args []string) error {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	out := fs.String("o", "", "write to this file instead of stdout")
	fs.Usage = usageFor(ctx.stderr, "snipkit export [-o file.json]", fs)
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
	b, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	if *out != "" {
		if err := os.WriteFile(*out, b, 0o600); err != nil {
			return err
		}
		fmt.Fprintf(ctx.stderr, "exported %d snippets to %s\n", len(all), *out)
		return nil
	}
	_, err = ctx.stdout.Write(b)
	return err
}
