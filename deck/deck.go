package deck

import (
	"fmt"
	"math/rand"

	"github.com/dangogh/GoPoker/cards"
)

type Deck struct {
	cards []cards.Card
}

func NewDeck() *Deck {
	cs := make([]cards.Card, 0, 52)
	for s := cards.Clubs; s <= cards.Spades; s++ {
		for r := cards.Two; r <= cards.Ace; r++ {
			cs = append(cs, cards.NewCard(s, r))
		}
	}
	return &Deck{cards: cs}
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Len() int { return len(d.cards) }

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

// RemoveCards removes the first occurrence of each card in toRemove from the deck's remaining cards.
// It returns the number of cards actually removed.
// Implementation: single-pass filtering using a frequency map and in-place slice reuse for efficiency.
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
