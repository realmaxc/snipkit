package cli

import (
	"fmt"
	"time"
)

func init() {
	register(command{name: "edit", summary: "edit a snippet in $EDITOR", run: runEdit})
}

func runEdit(ctx *appContext, args []string) error {
	if len(args) != 1 {
		fmt.Fprintln(ctx.stderr, "Usage: snipkit edit <id>")
		return errSilent
	}
	st, cfg, err := openStore()
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	sn, err := st.Get(args[0])
	if err != nil {
		return err
	}
	edited, err := editBuffer(cfg.Editor, sn.Body)
	if err != nil {
		return err
	}
	sn.Body = edited
	sn.UpdatedAt = time.Now()
	sn.Normalize()
	if err := sn.Validate(); err != nil {
		return err
	}
	if err := st.Update(sn); err != nil {
		return err
	}
	fmt.Fprintf(ctx.stdout, "updated %s\n", sn.ID)
	return nil
}
