package fuzzy

import "strings"

// Score rates how well pattern matches text. The boolean is false when there
// is no match at all. Higher scores are better: contiguous runs, matches at
// word boundaries, and matches near the start of the text all score higher.
func Score(pattern, text string) (int, bool) {
	ok, positions := Match(pattern, text)
	if !ok {
		return 0, false
	}
	if len(positions) == 0 {
		// An empty pattern matches everything, but only weakly.
		return 1, true
	}
	runes := []rune(strings.ToLower(text))
	total := 0
	prev := -2
	for _, pos := range positions {
		pts := 1 // base reward for a matched rune
		if pos == prev+1 {
			pts += 3 // contiguous with the previous match
		}
		if isBoundary(runes, pos) {
			pts += 4 // start of a word
		}
		total += pts
		prev = pos
	}
	total -= positions[0]                      // penalise a late first match
	total -= (len(runes) - len(positions)) / 8 // mild penalty for long text
	return total, true
}

func isBoundary(runes []rune, i int) bool {
	if i == 0 {
		return true
	}
	return isSeparator(runes[i-1])
}

func isSeparator(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '-', '_', '.', '/', ':':
		return true
	}
	return false
}
