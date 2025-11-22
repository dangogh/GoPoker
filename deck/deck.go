package deck

import (
	"fmt"
	"math/rand"

	"github.com/dangogh/GoPoker/cards"
)

// Deck represents a mutable stack of playing cards.
// It models a standard 52-card deck and supports dealing cards in order,
// shuffling, inspecting length, and removing specific cards for test or gameplay purposes.
type Deck struct {
	cards []cards.Card
}

// NewDeck builds a new standard 52-card deck in a deterministic order
// (Clubs -> Diamonds -> Hearts -> Spades; ranks Two -> Ace).
func NewDeck() *Deck {
	cs := make([]cards.Card, 0, 52)
	for s := cards.Clubs; s <= cards.Spades; s++ {
		for r := cards.Two; r <= cards.Ace; r++ {
			cs = append(cs, cards.NewCard(s, r))
		}
	}
	return &Deck{cards: cs}
}

// Shuffle randomly shuffles the remaining cards in the deck.
// Note: randomness source is the package-level math/rand; callers/tests may seed or use a custom source
// if deterministic behavior is required.
func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Len() int { return len(d.cards) }

// Deal removes and returns the next n cards from the deck (top of the deck).
// Returns an error for negative n or if there aren't enough cards remaining.
func (d *Deck) Deal(n int) ([]cards.Card, error) {
	if n < 0 {
		return nil, fmt.Errorf("negative deal count")
	}
	if n > len(d.cards) {
		return nil, fmt.Errorf("not enough cards to deal: requested %d, remaining %d", n, len(d.cards))
	}
	hand := d.cards[:n]
	d.cards = d.cards[n:]
	return hand, nil
}

// RemoveCards removes the first occurrences of the provided cards from the remaining deck.
// It returns the number of cards actually removed.
// This is useful in tests and advanced gameplay logic to ensure certain cards are not available.
func (d *Deck) RemoveCards(toRemove []cards.Card) int {
	if len(toRemove) == 0 || len(d.cards) == 0 {
		return 0
	}

	need := make(map[cards.Card]int, len(toRemove))
	for _, c := range toRemove {
		need[c]++
	}

	removed := 0
	// reuse underlying array to avoid extra allocation
	out := d.cards[:0]
	for _, c := range d.cards {
		if need[c] > 0 {
			need[c]--
			removed++
			continue
		}
		out = append(out, c)
	}
	d.cards = out
	return removed
}
