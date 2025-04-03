package cli

import (
	"strings"
	"testing"
)

func TestSearchFindsFuzzyMatch(t *testing.T) {
	ctx, out, _ := newTestEnv(t, "")
	add := func(title, body string) {
		ctx.stdin = strings.NewReader(body)
		if code := run(ctx, []string{"add", "-title", title, "-lang", "sh"}); code != exitOK {
			t.Fatalf("add %q exit", title)
		}
	}
	add("git checkout", "git checkout $1")
	add("list containers", "docker ps -a")

	ctx.stdin = strings.NewReader("")
	out.Reset()
	if code := run(ctx, []string{"search", "gco"}); code != exitOK {
		t.Fatalf("search exit")
	}
	got := out.String()
	if !strings.Contains(got, "git checkout") {
		t.Errorf("search gco should surface git checkout:\n%s", got)
	}
	if strings.Contains(got, "list containers") {
		t.Errorf("search gco should not surface unrelated snippet:\n%s", got)
	}
}
