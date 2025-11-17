package hand

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dangogh/GoPoker/cards"
)

func mk(cs ...cards.Card) Hand { return Hand{Cards: cs} }

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
