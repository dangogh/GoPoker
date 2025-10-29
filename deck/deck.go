package deck

import (
	"fmt"
	"math/rand"
	"time"

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
	rand.Seed(time.Now().UnixNano())
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
