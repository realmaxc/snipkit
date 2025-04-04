package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportImportRoundTrip(t *testing.T) {
	ctx, out, _ := newTestEnv(t, "echo one")
	if code := run(ctx, []string{"add", "-title", "one", "-lang", "sh", "-tag", "x"}); code != exitOK {
		t.Fatalf("add exit")
	}
	ctx.stdin = strings.NewReader("")
	out.Reset()
	if code := run(ctx, []string{"export"}); code != exitOK {
		t.Fatalf("export exit")
	}
	exported := out.String()
	if !strings.Contains(exported, "\"one\"") {
		t.Fatalf("export missing snippet:\n%s", exported)
	}

	// Import the export into a fresh store.
	t.Setenv("SNIPKIT_DB", filepath.Join(t.TempDir(), "second.json"))
	overrides.dbPath = ""
	var out2, err2 bytes.Buffer
	ctx2 := &appContext{stdout: &out2, stderr: &err2, stdin: strings.NewReader(exported)}
	if code := run(ctx2, []string{"import"}); code != exitOK {
		t.Fatalf("import exit; stderr=%s", err2.String())
	}
	out2.Reset()
	ctx2.stdin = strings.NewReader("")
	if code := run(ctx2, []string{"list"}); code != exitOK {
		t.Fatalf("list exit")
	}
	if !strings.Contains(out2.String(), "one") {
		t.Errorf("imported store missing snippet:\n%s", out2.String())
	}
}
