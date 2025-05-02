package cli

import (
	"fmt"

	"github.com/realmaxc/snipkit/internal/clip"
)

func init() {
	register(command{name: "copy", summary: "copy a snippet body to the clipboard", run: runCopy})
}

func runCopy(ctx *appContext, args []string) error {
	if len(args) != 1 {
		fmt.Fprintln(ctx.stderr, "Usage: snipkit copy <id|prefix>")
		return errSilent
	}
	st, _, err := openStore()
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	id, err := resolveID(st, args[0])
	if err != nil {
		return err
	}
	sn, err := st.Get(id)
	if err != nil {
		return err
	}
	if err := clip.Copy(sn.Body); err != nil {
		return err
	}
	fmt.Fprintf(ctx.stderr, "copied %s to clipboard\n", sn.ID)
	return nil
}
