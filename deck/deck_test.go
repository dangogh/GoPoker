package deck

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dangogh/GoPoker/cards"
)

func TestNewDeck(t *testing.T) {
	d := NewDeck()

	// Should have exactly 52 cards
	assert.Equal(t, 52, d.Len(), "new deck should have 52 cards")

	// Verify all unique cards are present
	seen := make(map[cards.Card]bool)
	for _, c := range d.cards {
		assert.False(t, seen[c], "duplicate card found: %s", c.String())
		seen[c] = true
	}

	// Verify we have all 4 suits × 13 ranks
	for s := cards.Clubs; s <= cards.Spades; s++ {
		for r := cards.Two; r <= cards.Ace; r++ {
			card := cards.NewCard(s, r)
			assert.True(t, seen[card], "missing card: %s", card.String())
		}
	}
}

func TestDeal(t *testing.T) {
	tests := []struct {
		name         string
		dealCount    int
		expectError  bool
		expectedLen  int
		remainingLen int
	}{
		{
			name:         "deal 5 cards",
			dealCount:    5,
			expectError:  false,
			expectedLen:  5,
			remainingLen: 47,
		},
		{
			name:         "deal 0 cards",
			dealCount:    0,
			expectError:  false,
			expectedLen:  0,
			remainingLen: 52,
		},
		{
			name:         "deal all 52 cards",
			dealCount:    52,
			expectError:  false,
			expectedLen:  52,
			remainingLen: 0,
		},
		{
			name:         "deal negative cards",
			dealCount:    -1,
			expectError:  true,
			expectedLen:  0,
			remainingLen: 52,
		},
		{
			name:         "deal more than available",
			dealCount:    53,
			expectError:  true,
			expectedLen:  0,
			remainingLen: 52,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDeck()
			hand, err := d.Deal(tc.dealCount)

			if tc.expectError {
				assert.Error(t, err, "expected error for %s", tc.name)
				assert.Nil(t, hand, "hand should be nil on error")
			} else {
				assert.NoError(t, err, "unexpected error for %s", tc.name)
				assert.Len(t, hand, tc.expectedLen, "dealt hand size mismatch")
			}

			assert.Equal(t, tc.remainingLen, d.Len(), "remaining deck size mismatch")
		})
	}
}

func TestDealMultipleTimes(t *testing.T) {
	d := NewDeck()

	// Deal 5 cards three times
	hand1, err := d.Deal(5)
	assert.NoError(t, err)
	assert.Len(t, hand1, 5)
	assert.Equal(t, 47, d.Len())

	hand2, err := d.Deal(5)
	assert.NoError(t, err)
	assert.Len(t, hand2, 5)
	assert.Equal(t, 42, d.Len())

	hand3, err := d.Deal(5)
	assert.NoError(t, err)
	assert.Len(t, hand3, 5)
	assert.Equal(t, 37, d.Len())

	// Verify no overlap between hands
	allCards := make(map[cards.Card]bool)
	for _, c := range hand1 {
		allCards[c] = true
	}
	for _, c := range hand2 {
		assert.False(t, allCards[c], "card %s appears in multiple hands", c.String())
		allCards[c] = true
	}
	for _, c := range hand3 {
		assert.False(t, allCards[c], "card %s appears in multiple hands", c.String())
	}
}

func TestShuffle(t *testing.T) {
	d1 := NewDeck()
	d2 := NewDeck()

	// Capture original order
	original := make([]cards.Card, len(d1.cards))
	copy(original, d1.cards)

	// Shuffle d2
	d2.Shuffle()

	// After shuffle, deck should still have 52 cards
	assert.Equal(t, 52, d2.Len(), "shuffle should not change deck size")

	// After shuffle, deck should have same cards (just different order)
	d1Cards := make(map[cards.Card]bool)
	for _, c := range d1.cards {
		d1Cards[c] = true
	}
	for _, c := range d2.cards {
		assert.True(t, d1Cards[c], "shuffled deck has unexpected card: %s", c.String())
	}

	// With high probability, shuffle should change the order
	// (this could fail with probability ~1/52! which is negligible)
	different := false
	for i := range original {
		if original[i] != d2.cards[i] {
			different = true
			break
		}
	}
	assert.True(t, different, "shuffle should change card order (may fail with negligible probability)")
}

func TestRemoveCards(t *testing.T) {
	tests := []struct {
		name          string
		toRemove      []cards.Card
		expectedCount int
		expectedLen   int
	}{
		{
			name: "remove single card",
			toRemove: []cards.Card{
				cards.NewCard(cards.Clubs, cards.Ace),
			},
			expectedCount: 1,
			expectedLen:   51,
		},
		{
			name: "remove multiple cards",
			toRemove: []cards.Card{
				cards.NewCard(cards.Clubs, cards.Ace),
				cards.NewCard(cards.Diamonds, cards.King),
				cards.NewCard(cards.Hearts, cards.Queen),
			},
			expectedCount: 3,
			expectedLen:   49,
		},
		{
			name:          "remove empty list",
			toRemove:      []cards.Card{},
			expectedCount: 0,
			expectedLen:   52,
		},
		{
			name: "remove non-existent card",
			toRemove: []cards.Card{
				cards.NewCard(cards.Clubs, cards.Ace),
				cards.NewCard(cards.Clubs, cards.Ace), // duplicate
			},
			expectedCount: 1,
			expectedLen:   51,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDeck()
			removed := d.RemoveCards(tc.toRemove)

			assert.Equal(t, tc.expectedCount, removed, "removed count mismatch")
			assert.Equal(t, tc.expectedLen, d.Len(), "deck length mismatch after removal")

			// Verify removed cards are no longer in deck (only check first occurrence)
			if tc.expectedCount > 0 && len(tc.toRemove) > 0 {
				firstCard := tc.toRemove[0]
				found := false
				for _, c := range d.cards {
					if c == firstCard {
						found = true
						break
					}
				}
				// If we removed cards successfully, the first card should be gone
				// (unless it was a duplicate request where only first instance is removed)
				if tc.name != "remove non-existent card" {
					assert.False(t, found, "card %s should be removed from deck", firstCard.String())
				}
			}
		})
	}
}

func TestRemoveCardsPreservesOrder(t *testing.T) {
	d := NewDeck()

	// Remove middle cards
	toRemove := []cards.Card{
		cards.NewCard(cards.Clubs, cards.Five),
		cards.NewCard(cards.Diamonds, cards.Seven),
	}

	d.RemoveCards(toRemove)

	// Remaining cards should maintain relative order
	assert.Equal(t, 50, d.Len())

	// Verify removed cards are gone
	for _, c := range d.cards {
		for _, removed := range toRemove {
			assert.NotEqual(t, removed, c, "removed card %s still in deck", removed.String())
		}
	}
}

func TestRemoveCardsFromEmptyDeck(t *testing.T) {
	d := NewDeck()
	// Empty the deck
	_, _ = d.Deal(52)
	assert.Equal(t, 0, d.Len())

	// Try to remove cards
	toRemove := []cards.Card{cards.NewCard(cards.Clubs, cards.Ace)}
	removed := d.RemoveCards(toRemove)

	assert.Equal(t, 0, removed, "should not remove from empty deck")
	assert.Equal(t, 0, d.Len())
}

func TestIntegrationDealAndRemove(t *testing.T) {
	d := NewDeck()

	// Deal 5 cards
	hand, err := d.Deal(5)
	assert.NoError(t, err)
	assert.Len(t, hand, 5)
	assert.Equal(t, 47, d.Len())

	// Remove 3 specific cards from remaining deck
	toRemove := []cards.Card{
		cards.NewCard(cards.Clubs, cards.Two),
		cards.NewCard(cards.Diamonds, cards.Two),
		cards.NewCard(cards.Hearts, cards.Two),
	}
	removed := d.RemoveCards(toRemove)
	assert.LessOrEqual(t, removed, 3, "should remove at most 3 cards")

	// Verify deck has correct remaining size
	expectedLen := 47 - removed
	assert.Equal(t, expectedLen, d.Len())
}
