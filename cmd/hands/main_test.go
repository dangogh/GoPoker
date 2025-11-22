package main

import (
	"bytes"
	"io"
	"os"
	"strings"
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

func TestPrintCards(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cs := []cards.Card{
		cards.NewCard(cards.Spades, cards.Ace),
		cards.NewCard(cards.Hearts, cards.King),
		cards.NewCard(cards.Diamonds, cards.Queen),
	}
	printCards(cs)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.Contains(t, output, "A♠", "output should contain Ace of Spades")
	assert.Contains(t, output, "K♥", "output should contain King of Hearts")
	assert.Contains(t, output, "Q♦", "output should contain Queen of Diamonds")
}

func TestPrintCardsEmpty(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printCards([]cards.Card{})

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.Equal(t, "\n", output, "empty card list should print only newline")
}

func TestPerformDrawDeckExhaustion(t *testing.T) {
	d := deck.NewDeck()
	// Deal most of the deck, leaving only 1 card
	_, _ = d.Deal(51)
	assert.Equal(t, 1, d.Len())

	// Create a high-card hand that will want to discard 3+ cards
	cs := []cards.Card{
		cards.NewCard(cards.Clubs, cards.Two),
		cards.NewCard(cards.Diamonds, cards.Three),
		cards.NewCard(cards.Hearts, cards.Four),
		cards.NewCard(cards.Spades, cards.Six),
		cards.NewCard(cards.Clubs, cards.Seven),
	}

	// Try to draw 3 cards when only 1 remains
	_, _, _, err := performDraw(d, cs, 3)
	assert.Error(t, err, "should error when deck exhausted during draw")
}

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		players     int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "invalid player count zero",
			players:     0,
			expectError: true,
			errorMsg:    "players must be > 0",
		},
		{
			name:        "invalid player count negative",
			players:     -1,
			expectError: true,
			errorMsg:    "players must be > 0",
		},
		{
			name:        "too many players exhausts deck",
			players:     15,
			expectError: true,
			errorMsg:    "deal error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Capture stdout to suppress output during tests
			old := os.Stdout
			_, w, _ := os.Pipe()
			os.Stdout = w

			err := run(tc.players)

			w.Close()
			os.Stdout = old

			if tc.expectError {
				assert.Error(t, err, "expected error for %s", tc.name)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg, "error message mismatch")
				}
			} else {
				assert.NoError(t, err, "unexpected error for %s", tc.name)
			}
		})
	}
}

func TestRunSuccess(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := run(2)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.NoError(t, err, "run should succeed with valid args")
	assert.Contains(t, output, "Initial hands:", "output should contain initial hands")
	assert.Contains(t, output, "Player 1:", "output should contain player 1")
	assert.Contains(t, output, "Player 2:", "output should contain player 2")
	assert.Contains(t, output, "Final hands:", "output should contain final hands")
	assert.Regexp(t, "Winner: Player [12]|Result: Tie", output, "output should contain winner or tie")
}

func TestRunTieBetweenPlayers(t *testing.T) {
	// This test is probabilistic since we use a shuffled deck
	// Run multiple times to increase chance of seeing a tie
	foundTie := false

	for attempt := 0; attempt < 10 && !foundTie; attempt++ {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := run(3)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		assert.NoError(t, err, "run should succeed")

		if strings.Contains(output, "Result: Tie") {
			foundTie = true
			assert.Contains(t, output, "Result: Tie among players", "should contain tie message")
			assert.Regexp(t, "Result: Tie among players( [0-9])+", output, "tie message should list player numbers")
		}
	}

	if !foundTie {
		t.Log("Note: No tie occurred in 10 attempts (expected with random shuffling)")
	}
}
