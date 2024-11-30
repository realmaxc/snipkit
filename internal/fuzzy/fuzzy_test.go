package fuzzy

import "testing"

func TestMatch(t *testing.T) {
	cases := []struct {
		pattern, text string
		want          bool
	}{
		{"gco", "git checkout", true},
		{"", "anything", true},
		{"xyz", "git checkout", false},
		{"GCO", "git checkout", true}, // case-insensitive
		{"tac", "cat", false},         // order matters
	}
	for _, c := range cases {
		if ok, _ := Match(c.pattern, c.text); ok != c.want {
			t.Errorf("Match(%q, %q) = %v, want %v", c.pattern, c.text, ok, c.want)
		}
	}
}

func TestScorePrefersContiguousAndBoundary(t *testing.T) {
	contiguous, ok := Score("chk", "checkout")
	if !ok {
		t.Fatal("expected a match")
	}
	scattered, ok := Score("chk", "cache kilo")
	if !ok {
		t.Fatal("expected a match")
	}
	if contiguous <= scattered {
		t.Errorf("contiguous score %d should beat scattered score %d", contiguous, scattered)
	}
}
