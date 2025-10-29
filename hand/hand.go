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
