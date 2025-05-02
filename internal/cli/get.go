package cli

import "fmt"

func init() {
	register(command{name: "get", summary: "print a snippet's body", run: runGet})
}

func runGet(ctx *appContext, args []string) error {
	if len(args) != 1 {
		fmt.Fprintln(ctx.stderr, "Usage: snipkit get <id|prefix>")
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
	fmt.Fprintln(ctx.stdout, sn.Body)
	return nil
}
