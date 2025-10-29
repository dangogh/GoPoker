package cards

import "fmt"

type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

type Rank int

const (
	Two   Rank = 2
	Three Rank = 3
	Four  Rank = 4
	Five  Rank = 5
	Six   Rank = 6
	Seven Rank = 7
	Eight Rank = 8
	Nine  Rank = 9
	Ten   Rank = 10
	Jack  Rank = 11
	Queen Rank = 12
	King  Rank = 13
	Ace   Rank = 14
)

type Card struct {
	Suit Suit
	Rank Rank
}

func NewCard(s Suit, r Rank) Card {
	return Card{Suit: s, Rank: r}
}

var rankNames = map[Rank]string{
	Two:   "2",
	Three: "3",
	Four:  "4",
	Five:  "5",
	Six:   "6",
	Seven: "7",
	Eight: "8",
	Nine:  "9",
	Ten:   "10",
	Jack:  "J",
	Queen: "Q",
	King:  "K",
	Ace:   "A",
}

var suitNames = map[Suit]string{
	Clubs:    "♣",
	Diamonds: "♦",
	Hearts:   "♥",
	Spades:   "♠",
}

func (c Card) String() string {
	r, okR := rankNames[c.Rank]
	s, okS := suitNames[c.Suit]
	if !okR || !okS {
		return fmt.Sprintf("Card(%d,%d)", c.Rank, c.Suit)
	}
	return r + s
}
