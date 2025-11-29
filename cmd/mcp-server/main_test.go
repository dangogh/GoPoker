package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dangogh/GoPoker/cards"
)

func TestParseCard(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantSuit  cards.Suit
		wantRank  cards.Rank
		wantError bool
	}{
		{
			name:     "Ace of Spades",
			input:    "A spades",
			wantSuit: cards.Spades,
			wantRank: cards.Ace,
		},
		{
			name:     "King of Hearts with symbol",
			input:    "K ♥",
			wantSuit: cards.Hearts,
			wantRank: cards.King,
		},
		{
			name:     "Ten of Clubs",
			input:    "10 clubs",
			wantSuit: cards.Clubs,
			wantRank: cards.Ten,
		},
		{
			name:      "Invalid format - no space",
			input:     "Aspades",
			wantError: true,
		},
		{
			name:      "Invalid rank",
			input:     "X spades",
			wantError: true,
		},
		{
			name:      "Invalid suit",
			input:     "A invalid",
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			card, err := parseCard(tc.input)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantSuit, card.Suit)
				assert.Equal(t, tc.wantRank, card.Rank)
			}
		})
	}
}
