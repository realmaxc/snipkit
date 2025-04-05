package snippet

import (
	"strings"
	"testing"
)

const idAlphabet = "0123456789abcdefghjkmnpqrstvwxyz"

func TestNewIDShapeAndUniqueness(t *testing.T) {
	seen := make(map[string]bool, 2000)
	for i := 0; i < 2000; i++ {
		id := NewID()
		if len(id) != 8 {
			t.Fatalf("NewID() = %q, want length 8", id)
		}
		if seen[id] {
			t.Fatalf("duplicate id generated: %q", id)
		}
		seen[id] = true
		for _, r := range id {
			if !strings.ContainsRune(idAlphabet, r) {
				t.Fatalf("id %q has out-of-alphabet rune %q", id, r)
			}
		}
	}
}
