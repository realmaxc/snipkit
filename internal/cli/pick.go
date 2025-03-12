package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/realmaxc/snipkit/internal/snippet"
)

func init() {
	register(command{name: "pick", summary: "interactively pick a snippet", run: runPick})
}

func runPick(ctx *appContext, args []string) error {
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
	if len(all) == 0 {
		fmt.Fprintln(ctx.stderr, "no snippets to pick from")
		return nil
	}
	sn, ok, err := pickSnippet(ctx, all)
	if err != nil {
		return err
	}
	if !ok {
		return nil // the user cancelled
	}
	fmt.Fprintln(ctx.stdout, sn.Body)
	return nil
}

// pickSnippet uses fzf when it is on PATH, otherwise it falls back to a plain
// numbered prompt so pick works anywhere.
func pickSnippet(ctx *appContext, snippets []snippet.Snippet) (snippet.Snippet, bool, error) {
	if path, err := exec.LookPath("fzf"); err == nil {
		return pickWithFzf(path, snippets)
	}
	return pickNumbered(ctx, snippets)
}

func pickWithFzf(path string, snippets []snippet.Snippet) (snippet.Snippet, bool, error) {
	var input strings.Builder
	for _, s := range snippets {
		fmt.Fprintf(&input, "%s\t%s\t%s\n", s.ID, dash(s.Lang), s.Title)
	}
	cmd := exec.Command(path, "--with-nth", "2..", "--delimiter", "\t", "--prompt", "snip> ")
	cmd.Stdin = strings.NewReader(input.String())
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) && ee.ExitCode() == 130 {
			return snippet.Snippet{}, false, nil // fzf exits 130 on cancel
		}
		return snippet.Snippet{}, false, err
	}
	id, _, _ := strings.Cut(strings.TrimSpace(string(out)), "\t")
	return findByID(snippets, id)
}

func pickNumbered(ctx *appContext, snippets []snippet.Snippet) (snippet.Snippet, bool, error) {
	for i, s := range snippets {
		fmt.Fprintf(ctx.stderr, "%3d  %-8s  %s\n", i+1, s.ID, s.Title)
	}
	fmt.Fprint(ctx.stderr, "select> ")
	line, err := bufio.NewReader(ctx.stdin).ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return snippet.Snippet{}, false, err
	}
	n, convErr := strconv.Atoi(line)
	if convErr != nil || n < 1 || n > len(snippets) {
		return snippet.Snippet{}, false, fmt.Errorf("invalid selection %q", line)
	}
	return snippets[n-1], true, nil
}

func findByID(snippets []snippet.Snippet, id string) (snippet.Snippet, bool, error) {
	for _, s := range snippets {
		if s.ID == id {
			return s, true, nil
		}
	}
	return snippet.Snippet{}, false, fmt.Errorf("selection %q not found", id)
}
