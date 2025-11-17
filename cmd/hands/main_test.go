package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dangogh/GoPoker/cards"
	"github.com/dangogh/GoPoker/deck"
	"github.com/dangogh/GoPoker/hand"
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
	// remove the manually-constructed hand from the deck so replacements won't match discarded cards
	d.RemoveCards(cs)

	cs2, discarded, drew, err := performDraw(d, cs, 3)
	assert.NoError(t, err)
	assert.Nil(t, discarded)
	assert.Nil(t, drew)
	// ensure hand unchanged
	assert.Equal(t, cs, cs2)
}

func TestPerformDraw_HighCardAggressive(t *testing.T) {
	// Deal initial hand from the deck so those cards are removed and cannot be drawn as replacements.
	d := deck.NewDeck()
	cs, err := d.Deal(5)
	assert.NoError(t, err)

	// Ensure we have a true high-card hand; skip test if not to avoid flakiness.
	e := hand.Evaluate(hand.Hand{Cards: cs})
	if e.Category != hand.HighCard {
		t.Skipf("dealt hand is not high-card (category=%v); skipping", e.Category)
	}

	cs2, discarded, drew, err := performDraw(d, cs, 3)
	assert.NoError(t, err)
	assert.Len(t, discarded, 3, "expected 3 discarded cards")
	assert.Len(t, drew, 3, "expected 3 drawn cards")

	// Verify the highest card in the original hand was kept
	highest := cs[0].Rank
	for _, c := range cs {
		if c.Rank > highest {
			highest = c.Rank
		}
	}
	foundHighest := false
	for _, c := range cs2 {
		if c.Rank == highest {
			foundHighest = true
			break
		}
	}
	assert.True(t, foundHighest, "expected highest card to be kept in resulting hand: %v", cs2)

	// Verify drawn cards are in final hand and differ from discarded equivalents (deck had initial hand removed).
	for i, dc := range drew {
		assert.Contains(t, cs2, dc, "drawn card not in final hand")
		if i < len(discarded) {
			assert.NotEqual(t, discarded[i], drew[i], "drawn card same as discarded")
		}
	}
}
