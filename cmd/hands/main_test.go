package main

import (
	"testing"

	"github.com/dangogh/GoPoker/cards"
	"github.com/dangogh/GoPoker/deck"
	"github.com/dangogh/GoPoker/hand"
	"github.com/stretchr/testify/assert"
)

func TestCategoryName(t *testing.T) {
	tests := map[hand.Category]string{
		hand.HighCard:      "High Card",
		hand.OnePair:       "One Pair",
		hand.TwoPair:       "Two Pair",
		hand.ThreeOfKind:   "Three of a Kind",
		hand.Straight:      "Straight",
		hand.Flush:         "Flush",
		hand.FullHouse:     "Full House",
		hand.FourOfKind:    "Four of a Kind",
		hand.StraightFlush: "Straight Flush",
	}
	for cat, want := range tests {
		got := categoryName(cat)
		assert.Equal(t, want, got, "categoryName(%v)", cat)
	}
}

func TestPerformDraw_NoDiscard(t *testing.T) {
	// Full house: keep (no discards)
	cs := []cards.Card{
		cards.NewCard(cards.Clubs, cards.King),
		cards.NewCard(cards.Diamonds, cards.King),
		cards.NewCard(cards.Hearts, cards.King),
		cards.NewCard(cards.Spades, cards.Ace),
		cards.NewCard(cards.Clubs, cards.Ace),
	}
	d := deck.NewDeck() // deterministic unshuffled deck
	cs2, discarded, drew, err := performDraw(d, cs, 3)
	assert.NoError(t, err)
	assert.Nil(t, discarded)
	assert.Nil(t, drew)
	// ensure hand unchanged
	assert.Equal(t, cs, cs2)
}

func TestPerformDraw_HighCardAggressive(t *testing.T) {
	// High-card hand: keep only top card (Ace), discard up to 3 lowest
	cs := []cards.Card{
		cards.NewCard(cards.Clubs, cards.Two),      // should be discarded
		cards.NewCard(cards.Diamonds, cards.Seven), // should be discarded
		cards.NewCard(cards.Hearts, cards.Four),    // should be discarded
		cards.NewCard(cards.Clubs, cards.Nine),     // should be kept
		cards.NewCard(cards.Spades, cards.Ace),     // should be kept (highest)
	}
	d := deck.NewDeck() // deterministic unshuffled deck; top cards differ from Spades suit -> replacements will be different
	cs2, discarded, drew, err := performDraw(d, cs, 3)
	assert.NoError(t, err)
	assert.Len(t, discarded, 3, "expected 3 discarded cards")
	assert.Len(t, drew, 3, "expected 3 drawn cards")

	// Verify exact cards that should be discarded (lowest 3)
	expectedDiscards := []cards.Card{
		cards.NewCard(cards.Clubs, cards.Two),
		cards.NewCard(cards.Hearts, cards.Four),
		cards.NewCard(cards.Diamonds, cards.Seven),
	}
	assert.ElementsMatch(t, expectedDiscards, discarded, "wrong cards were discarded")

	// Verify the two highest cards were kept
	assert.Contains(t, cs2, cards.NewCard(cards.Spades, cards.Ace), "Ace should be kept")
	assert.Contains(t, cs2, cards.NewCard(cards.Clubs, cards.Nine), "Nine should be kept")

	// Verify drawn cards are in final hand and different from discards
	for i, dc := range drew {
		assert.Contains(t, cs2, dc, "drawn card not in final hand")
		if i < len(discarded) {
			assert.NotEqual(t, discarded[i], drew[i], "drawn card same as discarded")
		}
	}
}
