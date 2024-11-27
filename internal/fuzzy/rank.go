package fuzzy

// score is a first cut at rating a set of matched positions: it simply sums
// the gaps between consecutive matches so tighter runs rank better.
// TODO: revisit scoring weights once real snippets are in the store.
func score(positions []int) int {
	if len(positions) == 0 {
		return 0
	}
	total := positions[0]
	for i := 1; i < len(positions); i++ {
		total += positions[i] - positions[i-1] - 1
	}
	return total
}
