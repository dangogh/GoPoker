package hand

import (
	"fmt"
	"sort"

	"github.com/dangogh/GoPoker/cards"
)

// Categories use iota so values are sequential and easy to compare.
// Keep the ordering so higher Category value means stronger hand.
type Category int

const (
	HighCard Category = iota
	OnePair
	TwoPair
	ThreeOfKind
	Straight
	Flush
	FullHouse
	FourOfKind
	StraightFlush
)

type Hand struct {
	Cards []cards.Card // expected length 5 for standard 5-card evaluation
}

type EvaluatedHand struct {
	Category Category
	Ranks    []cards.Rank // tiebreaker ranks in descending priority
}

// Evaluate computes the category and tiebreaker ranks for the hand.
func Evaluate(h Hand) EvaluatedHand {
	// Count ranks and gather suits
	rankCount := map[cards.Rank]int{}
	suitCount := map[cards.Suit]int{}
	for _, c := range h.Cards {
		rankCount[c.Rank]++
		suitCount[c.Suit]++
	}

	// Unique ranks slice
	uniqRanks := make([]cards.Rank, 0, len(rankCount))
	for r := range rankCount {
		uniqRanks = append(uniqRanks, r)
	}
	sort.Slice(uniqRanks, func(i, j int) bool { return uniqRanks[i] > uniqRanks[j] })

	// Helper: sorted ranks descending including duplicates (for high-card comparisons)
	allRanks := make([]cards.Rank, 0, len(h.Cards))
	for _, c := range h.Cards {
		allRanks = append(allRanks, c.Rank)
	}
	sort.Slice(allRanks, func(i, j int) bool { return allRanks[i] > allRanks[j] })

	// Detect flush
	isFlush := false
	if len(suitCount) == 1 {
		isFlush = true
	}

	// Detect straight (including wheel A-2-3-4-5)
	isStraight := false
	topStraightRank := cards.Rank(0)
	if len(uniqRanks) == 5 {
		// sorted uniqRanks descending
		max := uniqRanks[0]
		min := uniqRanks[4]
		if int(max)-int(min) == 4 {
			isStraight = true
			topStraightRank = max
		} else {
			// check wheel: A,5,4,3,2 -> top rank treated as 5
			hasAce := false
			hasTwo := false
			hasThree := false
			hasFour := false
			hasFive := false
			for _, r := range uniqRanks {
				switch r {
				case cards.Ace:
					hasAce = true
				case cards.Two:
					hasTwo = true
				case cards.Three:
					hasThree = true
				case cards.Four:
					hasFour = true
				case cards.Five:
					hasFive = true
				}
			}
			if hasAce && hasTwo && hasThree && hasFour && hasFive {
				isStraight = true
				topStraightRank = cards.Five
			}
		}
	}

	// Group ranks by count (e.g., pairs, threes, fours)
	countToRanks := map[int][]cards.Rank{}
	for r, cnt := range rankCount {
		countToRanks[cnt] = append(countToRanks[cnt], r)
	}
	for cnt := range countToRanks {
		sort.Slice(countToRanks[cnt], func(i, j int) bool { return countToRanks[cnt][i] > countToRanks[cnt][j] })
	}

	// Evaluate category and build tiebreaker ranks
	// Priority: StraightFlush, FourOfKind, FullHouse, Flush, Straight, ThreeOfKind, TwoPair, OnePair, HighCard
	switch {
	case isStraight && isFlush:
		return EvaluatedHand{Category: StraightFlush, Ranks: []cards.Rank{topStraightRank}}
	case len(countToRanks[4]) == 1:
		quad := countToRanks[4][0]
		// kicker
		kicker := cards.Rank(0)
		for _, r := range uniqRanks {
			if r != quad {
				kicker = r
				break
			}
		}
		return EvaluatedHand{Category: FourOfKind, Ranks: []cards.Rank{quad, kicker}}
	case len(countToRanks[3]) == 1 && len(countToRanks[2]) == 1:
		trip := countToRanks[3][0]
		pair := countToRanks[2][0]
		return EvaluatedHand{Category: FullHouse, Ranks: []cards.Rank{trip, pair}}
	case isFlush:
		// tiebreakers: all ranks desc
		return EvaluatedHand{Category: Flush, Ranks: allRanks}
	case isStraight:
		return EvaluatedHand{Category: Straight, Ranks: []cards.Rank{topStraightRank}}
	case len(countToRanks[3]) == 1:
		trip := countToRanks[3][0]
		// kickers descending
		kickers := make([]cards.Rank, 0, 2)
		for _, r := range allRanks {
			if r != trip {
				kickers = append(kickers, r)
			}
		}
		return EvaluatedHand{Category: ThreeOfKind, Ranks: append([]cards.Rank{trip}, kickers...)}
	case len(countToRanks[2]) == 2:
		// two pairs -> high pair first
		pairs := countToRanks[2]
		highPair := pairs[0]
		lowPair := pairs[1]
		kicker := cards.Rank(0)
		for _, r := range allRanks {
			if r != highPair && r != lowPair {
				kicker = r
				break
			}
		}
		return EvaluatedHand{Category: TwoPair, Ranks: []cards.Rank{highPair, lowPair, kicker}}
	case len(countToRanks[2]) == 1:
		pair := countToRanks[2][0]
		kickers := make([]cards.Rank, 0, 3)
		for _, r := range allRanks {
			if r != pair {
				kickers = append(kickers, r)
			}
		}
		return EvaluatedHand{Category: OnePair, Ranks: append([]cards.Rank{pair}, kickers...)}
	default:
		return EvaluatedHand{Category: HighCard, Ranks: allRanks}
	}
}

// RecommendDiscards returns the indices of cards to discard (0-based) up to maxDiscard.
// Aggressive strategy:
// - Do not break strong made hands (straight, flush, straight flush, full house, four of a kind).
// - If 4-to-a-flush exists, keep those 4 and discard the other card.
// - If 4-to-a-straight exists, keep those 4 and discard the other card.
// - For three-of-a-kind: keep the trips, discard other 2.
// - For two-pair: keep pairs, discard kicker (1).
// - For one-pair: keep the pair, discard 3 (aggressive).
// - For high-card hands: keep only the top card, discard up to maxDiscard (aggressive).
func RecommendDiscards(h Hand, maxDiscard int) []int {
	if maxDiscard <= 0 || len(h.Cards) == 0 {
		return nil
	}

	// If hand already strong (made hand), don't discard
	e := Evaluate(h)
	switch e.Category {
	case StraightFlush, FourOfKind, FullHouse, Straight, Flush:
		return nil
	}

	// Build counts and helper maps
	rankCount := map[cards.Rank]int{}
	suitCount := map[cards.Suit]int{}
	rankIdxs := map[cards.Rank][]int{}
	suitIdxs := map[cards.Suit][]int{}
	for i, c := range h.Cards {
		rankCount[c.Rank]++
		suitCount[c.Suit]++
		rankIdxs[c.Rank] = append(rankIdxs[c.Rank], i)
		suitIdxs[c.Suit] = append(suitIdxs[c.Suit], i)
	}

	// Helper: collect all unique ranks and a set for quick lookup (handle Ace as low for sequences)
	rankSet := map[int]bool{}
	for r := range rankCount {
		rankSet[int(r)] = true
		if r == cards.Ace {
			rankSet[1] = true // ace as low for wheel detection
		}
	}

	// Detect 4-to-a-flush
	for s, cnt := range suitCount {
		if cnt == 4 {
			keep := make([]bool, len(h.Cards))
			for _, idx := range suitIdxs[s] {
				keep[idx] = true
			}
			discards := make([]int, 0, 3)
			for i := range h.Cards {
				if !keep[i] {
					discards = append(discards, i)
				}
			}
			if len(discards) > maxDiscard {
				sort.Slice(discards, func(i, j int) bool {
					return h.Cards[discards[i]].Rank < h.Cards[discards[j]].Rank
				})
				discards = discards[:maxDiscard]
			}
			return discards
		}
	}

	// Detect 4-to-a-straight (look for any sequence of 4 consecutive ranks)
	if len(rankSet) >= 4 {
		uniqueVals := make([]int, 0, len(rankSet))
		for v := range rankSet {
			uniqueVals = append(uniqueVals, v)
		}
		sort.Ints(uniqueVals)

		runStart := -1
		runLen := 0
		bestRun := []int{}
		for i := 0; i < len(uniqueVals); i++ {
			if i == 0 || uniqueVals[i] == uniqueVals[i-1]+1 {
				if runLen == 0 {
					runStart = uniqueVals[i]
				}
				runLen++
			} else {
				if runLen >= 4 {
					curr := make([]int, 0, runLen)
					for v := runStart; v < runStart+runLen; v++ {
						curr = append(curr, v)
					}
					bestRun = curr
					break
				}
				runLen = 1
				runStart = uniqueVals[i]
			}
		}
		if runLen >= 4 && len(bestRun) == 0 {
			curr := make([]int, 0, runLen)
			for v := runStart; v < runStart+runLen; v++ {
				curr = append(curr, v)
			}
			bestRun = curr
		}
		if len(bestRun) >= 4 {
			keep := make([]bool, len(h.Cards))
			for i, c := range h.Cards {
				rv := int(c.Rank)
				if rv == 14 && containsInt(bestRun, 1) { // Ace counted as 1 in run
					keep[i] = true
					continue
				}
				if containsInt(bestRun, rv) {
					keep[i] = true
				}
			}
			discards := make([]int, 0, 3)
			for i := range h.Cards {
				if !keep[i] {
					discards = append(discards, i)
				}
			}
			if len(discards) > maxDiscard {
				sort.Slice(discards, func(i, j int) bool {
					return h.Cards[discards[i]].Rank < h.Cards[discards[j]].Rank
				})
				discards = discards[:maxDiscard]
			}
			return discards
		}
	}

	// Three of a kind -> discard other two
	for r, cnt := range rankCount {
		if cnt == 3 {
			discards := make([]int, 0, 2)
			for i, c := range h.Cards {
				if c.Rank != r {
					discards = append(discards, i)
				}
			}
			if len(discards) > maxDiscard {
				sort.Slice(discards, func(i, j int) bool {
					return h.Cards[discards[i]].Rank < h.Cards[discards[j]].Rank
				})
				discards = discards[:maxDiscard]
			}
			return discards
		}
	}

	// Two pair -> discard kicker
	pairRanks := make([]cards.Rank, 0, 2)
	for r, cnt := range rankCount {
		if cnt == 2 {
			pairRanks = append(pairRanks, r)
		}
	}
	if len(pairRanks) == 2 {
		keep := make([]bool, len(h.Cards))
		for _, pr := range pairRanks {
			for _, idx := range rankIdxs[pr] {
				keep[idx] = true
			}
		}
		discards := make([]int, 0, 1)
		for i := range h.Cards {
			if !keep[i] {
				discards = append(discards, i)
			}
		}
		if len(discards) > maxDiscard {
			discards = discards[:maxDiscard]
		}
		return discards
	}

	// One pair -> discard three
	if len(pairRanks) == 1 {
		pr := pairRanks[0]
		keep := make([]bool, len(h.Cards))
		for _, idx := range rankIdxs[pr] {
			keep[idx] = true
		}
		discards := make([]int, 0, 3)
		for i := range h.Cards {
			if !keep[i] {
				discards = append(discards, i)
			}
		}
		if len(discards) > maxDiscard {
			sort.Slice(discards, func(i, j int) bool {
				return h.Cards[discards[i]].Rank < h.Cards[discards[j]].Rank
			})
			discards = discards[:maxDiscard]
		}
		return discards
	}

	// High card: keep only the top card, discard up to maxDiscard
	idxs := make([]int, len(h.Cards))
	for i := range idxs {
		idxs[i] = i
	}
	sort.Slice(idxs, func(i, j int) bool {
		return h.Cards[idxs[i]].Rank > h.Cards[idxs[j]].Rank
	})
	keep := make([]bool, len(h.Cards))
	keep[idxs[0]] = true // keep top card only
	discards := make([]int, 0, len(h.Cards)-1)
	for i := range h.Cards {
		if !keep[i] {
			discards = append(discards, i)
		}
	}
	sort.Slice(discards, func(i, j int) bool {
		return h.Cards[discards[i]].Rank < h.Cards[discards[j]].Rank
	})
	if len(discards) > maxDiscard {
		discards = discards[:maxDiscard]
	}
	return discards
}

// helper to check presence in int slice
func containsInt(a []int, v int) bool {
	for _, x := range a {
		if x == v {
			return true
		}
	}
	return false
}

// Compare compares two evaluated hands. Returns 1 if a > b, -1 if a < b, 0 if equal.
func Compare(a, b EvaluatedHand) int {
	if a.Category != b.Category {
		if a.Category > b.Category {
			return 1
		}
		return -1
	}
	// Compare ranks lexicographically
	la, lb := a.Ranks, b.Ranks
	m := len(la)
	if len(lb) < m {
		m = len(lb)
	}
	for i := 0; i < m; i++ {
		if la[i] > lb[i] {
			return 1
		} else if la[i] < lb[i] {
			return -1
		}
	}
	// If all compared equal, longer rank list is greater (rare)
	if len(la) > len(lb) {
		return 1
	} else if len(la) < len(lb) {
		return -1
	}
	return 0
}

// String returns a human-readable name for the Category.
func (c Category) String() string {
	switch c {
	case HighCard:
		return "High Card"
	case OnePair:
		return "One Pair"
	case TwoPair:
		return "Two Pair"
	case ThreeOfKind:
		return "Three of a Kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	default:
		return fmt.Sprintf("Category(%d)", c)
	}
}
