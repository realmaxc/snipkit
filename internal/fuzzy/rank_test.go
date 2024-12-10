package fuzzy

import "testing"

func TestRankOrdersByRelevance(t *testing.T) {
	cands := []Candidate{
		{Index: 0, Text: "list all branches"},
		{Index: 1, Text: "git checkout branch"},
		{Index: 2, Text: "clone bare repo"},
	}
	got := Rank("gcb", cands)
	if len(got) == 0 {
		t.Fatal("expected at least one match")
	}
	if got[0].Index != 1 {
		t.Errorf("best match index = %d, want 1", got[0].Index)
	}
}

func TestRankDropsNonMatches(t *testing.T) {
	cands := []Candidate{{Index: 0, Text: "hello"}, {Index: 1, Text: "world"}}
	got := Rank("zzz", cands)
	if len(got) != 0 {
		t.Errorf("Rank returned %d results, want 0", len(got))
	}
}
