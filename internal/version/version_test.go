package version

import (
	"strings"
	"testing"
)

func TestStringHasBannerShape(t *testing.T) {
	got := String()
	if !strings.HasPrefix(got, "snipkit ") {
		t.Fatalf("String() = %q, want prefix %q", got, "snipkit ")
	}
	if !strings.Contains(got, Version) {
		t.Fatalf("String() = %q, want it to contain version %q", got, Version)
	}
	if !strings.Contains(got, Commit) {
		t.Fatalf("String() = %q, want it to contain commit %q", got, Commit)
	}
}
