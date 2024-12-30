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

func TestScoreRewardsContiguousRuns(t *testing.T) {
	// A tight run with no gaps should beat the same letters split by
	// non-boundary characters.
	tight, ok := Score("abc", "abcdef")
	if !ok {
		t.Fatal("expected a match")
	}
	gappy, ok := Score("abc", "axbxcx")
	if !ok {
		t.Fatal("expected a match")
	}
	if tight <= gappy {
		t.Errorf("tight run %d should beat gappy run %d", tight, gappy)
	}
}

func TestScoreRewardsWordInitials(t *testing.T) {
	// Acronym-style queries should favour word initials: "gco" is exactly the
	// initials of "git checkout".
	initials, ok := Score("gco", "git checkout")
	if !ok {
		t.Fatal("expected a match")
	}
	buried, ok := Score("gco", "goconut")
	if !ok {
		t.Fatal("expected a match")
	}
	if initials <= buried {
		t.Errorf("word-initial score %d should beat buried score %d", initials, buried)
	}
}
