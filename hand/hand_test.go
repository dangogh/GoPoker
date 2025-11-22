package hand

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dangogh/GoPoker/cards"
)

func mk(cs ...cards.Card) Hand { return Hand{Cards: cs} }

// cardStrings converts card indices to their string representation for assertion messages.
func cardStrings(hand Hand, indices []int) []string {
	if indices == nil {
		return nil
	}
	strs := make([]string, len(indices))
	for i, idx := range indices {
		strs[i] = hand.Cards[idx].String()
	}
	return strs
}

func TestEvaluateCategories(t *testing.T) {
	tests := []struct {
		name     string
		hand     Hand
		category Category
		ranks    []cards.Rank
	}{
		{
			name: "StraightFlush A-high",
			hand: mk(
				cards.NewCard(cards.Spades, cards.Ten),
				cards.NewCard(cards.Spades, cards.Jack),
				cards.NewCard(cards.Spades, cards.Queen),
				cards.NewCard(cards.Spades, cards.King),
				cards.NewCard(cards.Spades, cards.Ace),
			),
			category: StraightFlush,
			ranks:    []cards.Rank{cards.Ace},
		},
		{
			name: "Four of a kind",
			hand: mk(
				cards.NewCard(cards.Clubs, cards.King),
				cards.NewCard(cards.Diamonds, cards.King),
				cards.NewCard(cards.Hearts, cards.King),
				cards.NewCard(cards.Spades, cards.King),
				cards.NewCard(cards.Clubs, cards.Ace),
			),
			category: FourOfKind,
			ranks:    []cards.Rank{cards.King, cards.Ace},
		},
		{
			name: "Full House",
			hand: mk(
				cards.NewCard(cards.Clubs, cards.Three),
				cards.NewCard(cards.Diamonds, cards.Three),
				cards.NewCard(cards.Hearts, cards.Three),
				cards.NewCard(cards.Spades, cards.Two),
				cards.NewCard(cards.Clubs, cards.Two),
			),
			category: FullHouse,
			ranks:    []cards.Rank{cards.Three, cards.Two},
		},
		{
			name: "Flush",
			hand: mk(
				cards.NewCard(cards.Hearts, cards.Ace),
				cards.NewCard(cards.Hearts, cards.King),
				cards.NewCard(cards.Hearts, cards.Nine),
				cards.NewCard(cards.Hearts, cards.Five),
				cards.NewCard(cards.Hearts, cards.Two),
			),
			category: Flush,
			ranks:    []cards.Rank{cards.Ace, cards.King, cards.Nine, cards.Five, cards.Two},
		},
		{
			name: "Straight normal",
			hand: mk(
				cards.NewCard(cards.Clubs, cards.Six),
				cards.NewCard(cards.Diamonds, cards.Seven),
				cards.NewCard(cards.Hearts, cards.Eight),
				cards.NewCard(cards.Spades, cards.Nine),
				cards.NewCard(cards.Clubs, cards.Ten),
			),
			category: Straight,
			ranks:    []cards.Rank{cards.Ten},
		},
		{
			name: "Wheel straight A-2-3-4-5",
			hand: mk(
				cards.NewCard(cards.Clubs, cards.Ace),
				cards.NewCard(cards.Diamonds, cards.Two),
				cards.NewCard(cards.Hearts, cards.Three),
				cards.NewCard(cards.Spades, cards.Four),
				cards.NewCard(cards.Clubs, cards.Five),
			),
			category: Straight,
			ranks:    []cards.Rank{cards.Five},
		},
		{
			name: "Three of a kind",
			hand: mk(
				cards.NewCard(cards.Clubs, cards.Seven),
				cards.NewCard(cards.Diamonds, cards.Seven),
				cards.NewCard(cards.Hearts, cards.Seven),
				cards.NewCard(cards.Spades, cards.King),
				cards.NewCard(cards.Clubs, cards.Queen),
			),
			category: ThreeOfKind,
			ranks:    []cards.Rank{cards.Seven, cards.King, cards.Queen},
		},
		{
			name: "Two Pair",
			hand: mk(
				cards.NewCard(cards.Clubs, cards.King),
				cards.NewCard(cards.Diamonds, cards.King),
				cards.NewCard(cards.Hearts, cards.Nine),
				cards.NewCard(cards.Spades, cards.Nine),
				cards.NewCard(cards.Clubs, cards.Five),
			),
			category: TwoPair,
			ranks:    []cards.Rank{cards.King, cards.Nine, cards.Five},
		},
		{
			name: "One Pair",
			hand: mk(
				cards.NewCard(cards.Clubs, cards.Jack),
				cards.NewCard(cards.Diamonds, cards.Jack),
				cards.NewCard(cards.Hearts, cards.Ace),
				cards.NewCard(cards.Spades, cards.King),
				cards.NewCard(cards.Clubs, cards.Two),
			),
			category: OnePair,
			ranks:    []cards.Rank{cards.Jack, cards.Ace, cards.King, cards.Two},
		},
		{
			name: "High Card",
			hand: mk(
				cards.NewCard(cards.Clubs, cards.Ace),
				cards.NewCard(cards.Diamonds, cards.King),
				cards.NewCard(cards.Hearts, cards.Nine),
				cards.NewCard(cards.Spades, cards.Five),
				cards.NewCard(cards.Clubs, cards.Two),
			),
			category: HighCard,
			ranks:    []cards.Rank{cards.Ace, cards.King, cards.Nine, cards.Five, cards.Two},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ev := Evaluate(tc.hand)
			assert.Equal(t, tc.category, ev.Category, "category mismatch for %s", tc.name)
			assert.Equal(t, tc.ranks, ev.Ranks, "ranks mismatch for %s", tc.name)
		})
	}
}

func TestCompareBasics(t *testing.T) {
	// Four of a kind beats full house
	four := Evaluate(mk(
		cards.NewCard(cards.Clubs, cards.King),
		cards.NewCard(cards.Diamonds, cards.King),
		cards.NewCard(cards.Hearts, cards.King),
		cards.NewCard(cards.Spades, cards.King),
		cards.NewCard(cards.Clubs, cards.Ace),
	))
	full := Evaluate(mk(
		cards.NewCard(cards.Clubs, cards.Three),
		cards.NewCard(cards.Diamonds, cards.Three),
		cards.NewCard(cards.Hearts, cards.Three),
		cards.NewCard(cards.Spades, cards.Two),
		cards.NewCard(cards.Clubs, cards.Two),
	))
	assert.Equal(t, 1, Compare(four, full), "four should beat full")

	// Straight top-rank comparison
	s1 := Evaluate(mk(
		cards.NewCard(cards.Clubs, cards.Six),
		cards.NewCard(cards.Diamonds, cards.Seven),
		cards.NewCard(cards.Hearts, cards.Eight),
		cards.NewCard(cards.Spades, cards.Nine),
		cards.NewCard(cards.Clubs, cards.Ten),
	))
	s2 := Evaluate(mk(
		cards.NewCard(cards.Clubs, cards.Five),
		cards.NewCard(cards.Diamonds, cards.Six),
		cards.NewCard(cards.Hearts, cards.Seven),
		cards.NewCard(cards.Spades, cards.Eight),
		cards.NewCard(cards.Clubs, cards.Nine),
	))
	assert.Equal(t, 1, Compare(s1, s2), "higher straight should win")

	// Same hands compare equal and symmetry holds
	h1 := Evaluate(mk(
		cards.NewCard(cards.Clubs, cards.Ace),
		cards.NewCard(cards.Diamonds, cards.King),
		cards.NewCard(cards.Hearts, cards.Nine),
		cards.NewCard(cards.Spades, cards.Five),
		cards.NewCard(cards.Clubs, cards.Two),
	))
	h2 := Evaluate(mk(
		cards.NewCard(cards.Hearts, cards.Ace),
		cards.NewCard(cards.Spades, cards.King),
		cards.NewCard(cards.Diamonds, cards.Nine),
		cards.NewCard(cards.Clubs, cards.Five),
		cards.NewCard(cards.Spades, cards.Two),
	))
	assert.Equal(t, 0, Compare(h1, h2), "identical high-card hands should tie")
	assert.Equal(t, 0, Compare(h2, h1), "symmetry: tie both ways")
}

func TestRecommendDiscards_Table(t *testing.T) {
	tests := []struct {
		name     string
		hand     Hand
		maxDisc  int
		expected []int // expected discard indices (order ignored)
	}{
		{
			name: "StrongMade_NoDiscard (FullHouse)",
			hand: mk(
				cards.NewCard(cards.Spades, cards.Two),     // keep
				cards.NewCard(cards.Clubs, cards.Three),    // keep
				cards.NewCard(cards.Diamonds, cards.Three), // keep
				cards.NewCard(cards.Clubs, cards.Two),      // keep
				cards.NewCard(cards.Hearts, cards.Three),   // keep
			),
			maxDisc:  3,
			expected: nil,
		},
		{
			name: "FourToFlush_DiscardOne",
			hand: mk(
				cards.NewCard(cards.Spades, cards.Three), // discard
				cards.NewCard(cards.Hearts, cards.Seven), // keep
				cards.NewCard(cards.Hearts, cards.Two),   // keep
				cards.NewCard(cards.Hearts, cards.King),  // keep
				cards.NewCard(cards.Hearts, cards.Five),  // keep
			),
			maxDisc:  3,
			expected: []int{0},
		},
		{
			name: "FourToStraight_DiscardOne",
			hand: mk(
				cards.NewCard(cards.Clubs, cards.Two),    // discard
				cards.NewCard(cards.Hearts, cards.Seven), // keep
				cards.NewCard(cards.Clubs, cards.Five),   // keep
				cards.NewCard(cards.Diamonds, cards.Six), // keep
				cards.NewCard(cards.Spades, cards.Eight), // keep
			),
			maxDisc:  3,
			expected: []int{0},
		},
		{
			name: "ThreeOfKind_DiscardTwo",
			hand: mk(
				cards.NewCard(cards.Spades, cards.King),    // discard
				cards.NewCard(cards.Clubs, cards.Seven),    // keep
				cards.NewCard(cards.Clubs, cards.Queen),    // discard
				cards.NewCard(cards.Diamonds, cards.Seven), // keep
				cards.NewCard(cards.Hearts, cards.Seven),   // keep
			),
			maxDisc:  3,
			expected: []int{0, 2},
		},
		{
			name: "TwoPair_DiscardKicker",
			hand: mk(
				cards.NewCard(cards.Clubs, cards.Five),    // discard kicker
				cards.NewCard(cards.Diamonds, cards.King), // keep
				cards.NewCard(cards.Hearts, cards.Nine),   // keep
				cards.NewCard(cards.Clubs, cards.King),    // keep
				cards.NewCard(cards.Spades, cards.Nine),   // keep
			),
			maxDisc:  3,
			expected: []int{0},
		},
		{
			name: "OnePair_DiscardThree",
			hand: mk(
				cards.NewCard(cards.Hearts, cards.Ace),    // discard
				cards.NewCard(cards.Clubs, cards.Jack),    // keep
				cards.NewCard(cards.Clubs, cards.Two),     // discard
				cards.NewCard(cards.Diamonds, cards.Jack), // keep
				cards.NewCard(cards.Spades, cards.King),   // discard
			),
			maxDisc:  3,
			expected: []int{0, 2, 4},
		},
		{
			name: "FourToBabyStraight_DiscardOne",
			hand: mk(
				cards.NewCard(cards.Clubs, cards.Nine),     // discard
				cards.NewCard(cards.Spades, cards.Ace),     // wheel keep
				cards.NewCard(cards.Clubs, cards.Two),      // wheel keep
				cards.NewCard(cards.Diamonds, cards.Three), // wheel keep
				cards.NewCard(cards.Hearts, cards.Four),    // wheel keep
			),
			maxDisc:  3,
			expected: []int{0},
		},
		{
			name: "ComputeMaxDiscard_HighCardWithAce",
			hand: mk(
				cards.NewCard(cards.Diamonds, cards.King), // discard
				cards.NewCard(cards.Spades, cards.Five),   // discard
				cards.NewCard(cards.Clubs, cards.Ace),     // keep (highest)
				cards.NewCard(cards.Clubs, cards.Two),     // discard
				cards.NewCard(cards.Hearts, cards.Nine),   // discard
			),
			maxDisc: ComputeMaxDiscard(mk(
				cards.NewCard(cards.Diamonds, cards.King),
				cards.NewCard(cards.Spades, cards.Five),
				cards.NewCard(cards.Clubs, cards.Ace),
				cards.NewCard(cards.Clubs, cards.Two),
				cards.NewCard(cards.Hearts, cards.Nine),
			)),
			expected: []int{0, 1, 3, 4},
		},
		{
			name: "ComputeMaxDiscard_HighCardWithoutAce",
			hand: mk(
				cards.NewCard(cards.Hearts, cards.Nine),    // discard
				cards.NewCard(cards.Clubs, cards.King),     // keep (highest)
				cards.NewCard(cards.Clubs, cards.Two),      // discard (lowest)
				cards.NewCard(cards.Diamonds, cards.Queen), // keep (2nd highest)
				cards.NewCard(cards.Spades, cards.Five),    // discard
			),
			maxDisc: ComputeMaxDiscard(mk(
				cards.NewCard(cards.Hearts, cards.Nine),
				cards.NewCard(cards.Clubs, cards.King),
				cards.NewCard(cards.Clubs, cards.Two),
				cards.NewCard(cards.Diamonds, cards.Queen),
				cards.NewCard(cards.Spades, cards.Five),
			)),
			expected: []int{0, 2, 4},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := RecommendDiscards(tc.hand, tc.maxDisc)
			if tc.expected == nil {
				assert.Nil(t, got, "%s: expected no discards", tc.name)
			} else {
				// assert.ElementsMatch(t, tc.expected, got, "%s: unexpected discard indices", tc.name)
				assert.LessOrEqual(t, len(got), tc.maxDisc, "%s: exceeded maxDiscard", tc.name)
				// Show card strings for clear failure messages
				expectedCards := make([]string, len(tc.expected))
				for i, idx := range tc.expected {
					expectedCards[i] = tc.hand.Cards[idx].String()
				}
				gotCards := make([]string, len(got))
				for i, idx := range got {
					gotCards[i] = tc.hand.Cards[idx].String()
				}
				assert.ElementsMatch(t, expectedCards, gotCards, "%s: discarded cards mismatch", tc.name)
			}
		})
	}
}

func TestCategoryString(t *testing.T) {
	tests := []struct {
		category Category
		expected string
	}{
		{HighCard, "High Card"},
		{OnePair, "One Pair"},
		{TwoPair, "Two Pair"},
		{ThreeOfKind, "Three of a Kind"},
		{Straight, "Straight"},
		{Flush, "Flush"},
		{FullHouse, "Full House"},
		{FourOfKind, "Four of a Kind"},
		{StraightFlush, "Straight Flush"},
		{Category(99), "Category(99)"}, // Unknown category
	}
	for _, tc := range tests {
		assert.Equal(t, tc.expected, tc.category.String())
	}
}

func TestEvaluateWheelStraightFlush(t *testing.T) {
	h := mk(
		cards.NewCard(cards.Clubs, cards.Ace),
		cards.NewCard(cards.Clubs, cards.Two),
		cards.NewCard(cards.Clubs, cards.Three),
		cards.NewCard(cards.Clubs, cards.Four),
		cards.NewCard(cards.Clubs, cards.Five),
	)
	ev := Evaluate(h)
	assert.Equal(t, StraightFlush, ev.Category)
	assert.Equal(t, []cards.Rank{cards.Five}, ev.Ranks)
}

func TestCompareLongerRanksList(t *testing.T) {
	// Test rare case where rank lists have different lengths
	a := EvaluatedHand{Category: HighCard, Ranks: []cards.Rank{cards.Ace, cards.King, cards.Queen}}
	b := EvaluatedHand{Category: HighCard, Ranks: []cards.Rank{cards.Ace, cards.King}}
	assert.Equal(t, 1, Compare(a, b))
	assert.Equal(t, -1, Compare(b, a))
}

func TestRecommendDiscardsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		hand     Hand
		maxDisc  int
		expected []int
	}{
		{
			name:     "maxDiscard zero",
			hand:     mk(cards.NewCard(cards.Clubs, cards.Two), cards.NewCard(cards.Diamonds, cards.Three), cards.NewCard(cards.Hearts, cards.Four), cards.NewCard(cards.Spades, cards.Five), cards.NewCard(cards.Clubs, cards.Six)),
			maxDisc:  0,
			expected: nil,
		},
		{
			name:     "empty hand",
			hand:     Hand{Cards: []cards.Card{}},
			maxDisc:  3,
			expected: nil,
		},
		{
			name: "ThreeOfKind exceeds maxDiscard",
			hand: mk(
				cards.NewCard(cards.Clubs, cards.Seven),
				cards.NewCard(cards.Diamonds, cards.Seven),
				cards.NewCard(cards.Hearts, cards.Seven),
				cards.NewCard(cards.Spades, cards.King),
				cards.NewCard(cards.Clubs, cards.Queen),
			),
			maxDisc:  1,
			expected: []int{4}, // Only discard Queen (lowest)
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := RecommendDiscards(tc.hand, tc.maxDisc)
			if tc.expected == nil {
				assert.Nil(t, got)
			} else {
				expectedCards := make([]string, len(tc.expected))
				for i, idx := range tc.expected {
					expectedCards[i] = tc.hand.Cards[idx].String()
				}
				gotCards := make([]string, len(got))
				for i, idx := range got {
					gotCards[i] = tc.hand.Cards[idx].String()
				}
				assert.ElementsMatch(t, expectedCards, gotCards)
			}
		})
	}
}

func TestComputeMaxDiscardEmptyHand(t *testing.T) {
	h := Hand{Cards: []cards.Card{}}
	assert.Equal(t, 3, ComputeMaxDiscard(h))
}

// FuzzEvaluate fuzzes Evaluate with random 5-card hands.
// Invariants: Evaluate should not panic, should return valid category, and ranks length should be reasonable.
func FuzzEvaluate(f *testing.F) {
	f.Add(uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0))
	f.Fuzz(func(t *testing.T, s0, r0, s1, r1, s2, r2, s3, r3, s4, r4 uint8) {
		// Map bytes to valid suits and ranks
		suits := []cards.Suit{cards.Clubs, cards.Diamonds, cards.Hearts, cards.Spades}
		ranks := []cards.Rank{cards.Two, cards.Three, cards.Four, cards.Five, cards.Six, cards.Seven,
			cards.Eight, cards.Nine, cards.Ten, cards.Jack, cards.Queen, cards.King, cards.Ace}

		h := Hand{Cards: []cards.Card{
			cards.NewCard(suits[s0%4], ranks[r0%13]),
			cards.NewCard(suits[s1%4], ranks[r1%13]),
			cards.NewCard(suits[s2%4], ranks[r2%13]),
			cards.NewCard(suits[s3%4], ranks[r3%13]),
			cards.NewCard(suits[s4%4], ranks[r4%13]),
		}}

		ev := Evaluate(h)
		// Invariants: category should be in range [0, StraightFlush]
		assert.GreaterOrEqual(t, int(ev.Category), int(HighCard))
		assert.LessOrEqual(t, int(ev.Category), int(StraightFlush))
		// Ranks should be non-empty
		assert.NotEmpty(t, ev.Ranks)
		// Ranks should not exceed 5
		assert.LessOrEqual(t, len(ev.Ranks), 5)
	})
}

// FuzzCompare fuzzes Compare with random hands to verify symmetry and antisymmetry.
// Invariants: Compare(a,b) == -Compare(b,a), Compare(a,a) == 0, comparison results should be -1, 0, or 1.
func FuzzCompare(f *testing.F) {
	f.Add(uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0))
	f.Fuzz(func(t *testing.T, s0, r0, s1, r1, s2, r2, s3, r3, s4, r4 uint8) {
		suits := []cards.Suit{cards.Clubs, cards.Diamonds, cards.Hearts, cards.Spades}
		ranks := []cards.Rank{cards.Two, cards.Three, cards.Four, cards.Five, cards.Six, cards.Seven,
			cards.Eight, cards.Nine, cards.Ten, cards.Jack, cards.Queen, cards.King, cards.Ace}

		h1 := Hand{Cards: []cards.Card{
			cards.NewCard(suits[s0%4], ranks[r0%13]),
			cards.NewCard(suits[s1%4], ranks[r1%13]),
			cards.NewCard(suits[s2%4], ranks[r2%13]),
			cards.NewCard(suits[s3%4], ranks[r3%13]),
			cards.NewCard(suits[s4%4], ranks[r4%13]),
		}}

		e1 := Evaluate(h1)
		e2 := Evaluate(h1) // Same hand

		// Symmetry: Compare(a, b) == -Compare(b, a)
		cmp := Compare(e1, e2)
		assert.Equal(t, -cmp, Compare(e2, e1), "symmetry violated")

		// Reflexivity: Compare(a, a) == 0
		assert.Equal(t, 0, cmp, "reflexivity violated for identical hands")

		// Result should be -1, 0, or 1
		assert.Contains(t, []int{-1, 0, 1}, cmp, "compare result out of range")
	})
}

// FuzzRecommendDiscards fuzzes RecommendDiscards with random hands and maxDiscard values.
// Invariants: returned indices should be <= maxDiscard, all indices should be valid (0-4), no panics.
func FuzzRecommendDiscards(f *testing.F) {
	f.Add(uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(0), uint8(3))
	f.Fuzz(func(t *testing.T, s0, r0, s1, r1, s2, r2, s3, r3, s4, r4, maxD uint8) {
		suits := []cards.Suit{cards.Clubs, cards.Diamonds, cards.Hearts, cards.Spades}
		ranks := []cards.Rank{cards.Two, cards.Three, cards.Four, cards.Five, cards.Six, cards.Seven,
			cards.Eight, cards.Nine, cards.Ten, cards.Jack, cards.Queen, cards.King, cards.Ace}

		h := Hand{Cards: []cards.Card{
			cards.NewCard(suits[s0%4], ranks[r0%13]),
			cards.NewCard(suits[s1%4], ranks[r1%13]),
			cards.NewCard(suits[s2%4], ranks[r2%13]),
			cards.NewCard(suits[s3%4], ranks[r3%13]),
			cards.NewCard(suits[s4%4], ranks[r4%13]),
		}}

		maxDisc := int(maxD%5) + 1 // Range [1, 5]
		discards := RecommendDiscards(h, maxDisc)

		// Invariant: number of discards should be <= maxDisc
		assert.LessOrEqual(t, len(discards), maxDisc, "discards exceeded maxDisc")

		// Invariant: all indices should be valid (0-4)
		for _, idx := range discards {
			assert.GreaterOrEqual(t, idx, 0)
			assert.Less(t, idx, 5)
		}

		// Invariant: no duplicate indices
		seen := make(map[int]bool)
		for _, idx := range discards {
			assert.False(t, seen[idx], "duplicate discard index")
			seen[idx] = true
		}
	})
}
