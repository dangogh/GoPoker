package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCard(t *testing.T) {
	tests := []struct {
		name string
		suit Suit
		rank Rank
	}{
		{"Ace of Spades", Spades, Ace},
		{"Two of Clubs", Clubs, Two},
		{"King of Hearts", Hearts, King},
		{"Queen of Diamonds", Diamonds, Queen},
		{"Jack of Clubs", Clubs, Jack},
		{"Ten of Spades", Spades, Ten},
		{"Five of Hearts", Hearts, Five},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			card := NewCard(tc.suit, tc.rank)
			assert.Equal(t, tc.suit, card.Suit, "suit mismatch for %s", tc.name)
			assert.Equal(t, tc.rank, card.Rank, "rank mismatch for %s", tc.name)
		})
	}
}

func TestCardString(t *testing.T) {
	tests := []struct {
		name     string
		card     Card
		expected string
	}{
		{"Ace of Spades", NewCard(Spades, Ace), "A♠"},
		{"King of Hearts", NewCard(Hearts, King), "K♥"},
		{"Queen of Diamonds", NewCard(Diamonds, Queen), "Q♦"},
		{"Jack of Clubs", NewCard(Clubs, Jack), "J♣"},
		{"Ten of Spades", NewCard(Spades, Ten), "10♠"},
		{"Nine of Hearts", NewCard(Hearts, Nine), "9♥"},
		{"Eight of Diamonds", NewCard(Diamonds, Eight), "8♦"},
		{"Seven of Clubs", NewCard(Clubs, Seven), "7♣"},
		{"Six of Spades", NewCard(Spades, Six), "6♠"},
		{"Five of Hearts", NewCard(Hearts, Five), "5♥"},
		{"Four of Diamonds", NewCard(Diamonds, Four), "4♦"},
		{"Three of Clubs", NewCard(Clubs, Three), "3♣"},
		{"Two of Spades", NewCard(Spades, Two), "2♠"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.card.String()
			assert.Equal(t, tc.expected, got, "string representation mismatch for %s", tc.name)
		})
	}
}

func TestCardStringInvalidRank(t *testing.T) {
	card := Card{Suit: Clubs, Rank: Rank(99)}
	str := card.String()
	assert.Equal(t, "Card(99,0)", str, "invalid rank should return fallback format")
}

func TestCardStringInvalidSuit(t *testing.T) {
	card := Card{Suit: Suit(99), Rank: Ace}
	str := card.String()
	assert.Equal(t, "Card(14,99)", str, "invalid suit should return fallback format")
}

func TestCardStringInvalidBoth(t *testing.T) {
	card := Card{Suit: Suit(99), Rank: Rank(88)}
	str := card.String()
	assert.Equal(t, "Card(88,99)", str, "invalid suit and rank should return fallback format")
}

func TestAllSuits(t *testing.T) {
	suits := []Suit{Clubs, Diamonds, Hearts, Spades}
	expected := []string{"♣", "♦", "♥", "♠"}

	for i, s := range suits {
		card := NewCard(s, Ace)
		str := card.String()
		assert.Contains(t, str, expected[i], "suit symbol missing for suit %d", s)
	}
}

func TestAllRanks(t *testing.T) {
	ranks := []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
	expected := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	for i, r := range ranks {
		card := NewCard(Clubs, r)
		str := card.String()
		assert.Contains(t, str, expected[i], "rank symbol missing for rank %d", r)
	}
}

func TestCardEquality(t *testing.T) {
	c1 := NewCard(Spades, Ace)
	c2 := NewCard(Spades, Ace)
	c3 := NewCard(Hearts, Ace)

	assert.Equal(t, c1, c2, "identical cards should be equal")
	assert.NotEqual(t, c1, c3, "cards with different suits should not be equal")
}

func TestRankValues(t *testing.T) {
	// Verify rank ordering
	assert.Less(t, int(Two), int(Three))
	assert.Less(t, int(Three), int(Four))
	assert.Less(t, int(Jack), int(Queen))
	assert.Less(t, int(Queen), int(King))
	assert.Less(t, int(King), int(Ace))
	assert.Equal(t, 2, int(Two))
	assert.Equal(t, 14, int(Ace))
}

func TestSuitValues(t *testing.T) {
	// Verify suit enum values
	assert.Equal(t, 0, int(Clubs))
	assert.Equal(t, 1, int(Diamonds))
	assert.Equal(t, 2, int(Hearts))
	assert.Equal(t, 3, int(Spades))
}
