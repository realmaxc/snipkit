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

// Candidate is an item to be ranked, pairing a caller-defined index with the
// text the query is matched against.
type Candidate struct {
	Index int
	Text  string
}

// Result is a candidate that matched a query, together with its score.
type Result struct {
	Index int
	Score int
}

// Rank scores every candidate against query and returns the matches sorted
// best-first. Candidates that do not match are dropped. An empty query keeps
// the candidates in their original order.
func Rank(query string, candidates []Candidate) []Result {
	results := make([]Result, 0, len(candidates))
	for _, c := range candidates {
		s, ok := Score(query, c.Text)
		if !ok {
			continue
		}
		results = append(results, Result{Index: c.Index, Score: s})
	}
	sortByScore(results)
	return results
}

func sortByScore(results []Result) {
	// Insertion sort keeps this dependency-free and is stable, which matters
	// for predictable ordering of equally-scored results.
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Score > results[j-1].Score; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
}
