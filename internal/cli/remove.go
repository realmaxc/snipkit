package cli

import "fmt"

func init() {
	register(command{name: "rm", summary: "remove one or more snippets", run: runRemove})
}

func runRemove(ctx *appContext, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(ctx.stderr, "Usage: snipkit rm <id|prefix>...")
		return errSilent
	}
	st, _, err := openStore()
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	failed := 0
	for _, arg := range args {
		id, err := resolveID(st, arg)
		if err != nil {
			fmt.Fprintf(ctx.stderr, "snipkit: %s: %v\n", arg, err)
			failed++
			continue
		}
		if err := st.Delete(id); err != nil {
			fmt.Fprintf(ctx.stderr, "snipkit: %s: %v\n", arg, err)
			failed++
			continue
		}
		fmt.Fprintf(ctx.stdout, "removed %s\n", id)
	}
	if failed > 0 {
		return fmt.Errorf("%d of %d ids could not be removed", failed, len(args))
	}
	return nil
}
