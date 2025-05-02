package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// newTestEnv points snipkit at a throwaway JSON store and returns an
// appContext wired to buffers plus the captured stdout/stderr.
func newTestEnv(t *testing.T, stdin string) (*appContext, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	t.Setenv("SNIPKIT_DB", filepath.Join(t.TempDir(), "snips.json"))
	t.Setenv("SNIPKIT_BACKEND", "json")
	overrides.dbPath = ""
	overrides.backend = ""
	var out, errb bytes.Buffer
	ctx := &appContext{stdout: &out, stderr: &errb, stdin: strings.NewReader(stdin)}
	return ctx, &out, &errb
}

func TestAddThenList(t *testing.T) {
	ctx, out, _ := newTestEnv(t, "echo hello")
	if code := run(ctx, []string{"add", "-title", "greet", "-lang", "sh", "-tag", "demo"}); code != exitOK {
		t.Fatalf("add exit = %d", code)
	}

	ctx.stdin = strings.NewReader("")
	out.Reset()
	if code := run(ctx, []string{"list"}); code != exitOK {
		t.Fatalf("list exit = %d", code)
	}
	if !strings.Contains(out.String(), "greet") {
		t.Errorf("list output missing snippet:\n%s", out.String())
	}
}

func TestGetMissingIsNotFound(t *testing.T) {
	ctx, _, _ := newTestEnv(t, "")
	if code := run(ctx, []string{"get", "nope"}); code != exitNotFound {
		t.Errorf("get missing exit = %d, want %d", code, exitNotFound)
	}
}

func TestResolveByIDPrefix(t *testing.T) {
	ctx, out, _ := newTestEnv(t, "echo hi")
	if code := run(ctx, []string{"add", "-title", "greet"}); code != exitOK {
		t.Fatalf("add exit")
	}
	id := strings.Fields(out.String())[1]

	out.Reset()
	ctx.stdin = strings.NewReader("")
	if code := run(ctx, []string{"get", id[:3]}); code != exitOK {
		t.Fatalf("get by prefix exit")
	}
	if !strings.Contains(out.String(), "echo hi") {
		t.Errorf("get by prefix returned %q", out.String())
	}
}
