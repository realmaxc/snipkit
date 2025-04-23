package snippet

import "testing"

func TestNormalizeDerivesTitleAndCleansTags(t *testing.T) {
	s := Snippet{
		Body: "  echo hi\n\n",
		Tags: []string{"B", "a", "a", " c "},
		Lang: " SH ",
	}
	s.Normalize()
	if s.Title != "echo hi" {
		t.Errorf("Title = %q, want the derived first line", s.Title)
	}
	if s.Lang != "sh" {
		t.Errorf("Lang = %q, want lowercased", s.Lang)
	}
	want := []string{"a", "b", "c"}
	if len(s.Tags) != len(want) {
		t.Fatalf("Tags = %v, want %v", s.Tags, want)
	}
	for i := range want {
		if s.Tags[i] != want[i] {
			t.Errorf("Tags[%d] = %q, want %q", i, s.Tags[i], want[i])
		}
	}
}

func TestValidateRejectsEmptyBody(t *testing.T) {
	var s Snippet
	if err := s.Validate(); err == nil {
		t.Error("expected an empty body to be invalid")
	}
	s.Body = "ls"
	if err := s.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestHasTagIsCaseInsensitive(t *testing.T) {
	s := Snippet{Tags: []string{"git", "cli"}}
	if !s.HasTag("GIT") {
		t.Error("HasTag should match case-insensitively")
	}
	if s.HasTag("nope") {
		t.Error("HasTag should be false for a missing tag")
	}
}
