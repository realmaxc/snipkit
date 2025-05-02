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

	for _, arg := range args {
		id, err := resolveID(st, arg)
		if err != nil {
			return err
		}
		if err := st.Delete(id); err != nil {
			return err
		}
		fmt.Fprintf(ctx.stdout, "removed %s\n", id)
	}
	return nil
}
