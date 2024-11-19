// Package fuzzy provides a small subsequence matcher used to rank snippets
// against a query typed at the terminal.
package fuzzy

import "strings"

// Match reports whether every rune of pattern appears in text in order,
// case-insensitively. When it matches it also returns the rune positions (into
// the lower-cased text) where each pattern rune landed.
func Match(pattern, text string) (bool, []int) {
	if pattern == "" {
		return true, nil
	}
	p := []rune(strings.ToLower(pattern))
	t := []rune(strings.ToLower(text))
	positions := make([]int, 0, len(p))
	ti := 0
	for _, pr := range p {
		for ti < len(t) && t[ti] != pr {
			ti++
		}
		if ti == len(t) {
			return false, nil
		}
		positions = append(positions, ti)
		ti++
	}
	return true, positions
}
