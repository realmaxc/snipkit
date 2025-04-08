package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestUsageListsCommands(t *testing.T) {
	var b bytes.Buffer
	usage(&b)
	for _, want := range []string{"add", "list", "search", "pick", "export", "import", "copy", "edit", "rm", "tags"} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("usage output missing command %q:\n%s", want, b.String())
		}
	}
}

func TestUnknownCommandIsUsageError(t *testing.T) {
	ctx := &appContext{stdout: &bytes.Buffer{}, stderr: &bytes.Buffer{}, stdin: strings.NewReader("")}
	if code := run(ctx, []string{"definitely-not-a-command"}); code != exitUsage {
		t.Errorf("unknown command exit = %d, want %d", code, exitUsage)
	}
}
