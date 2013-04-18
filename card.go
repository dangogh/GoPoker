package poker

import (
	"fmt"
)

type Card struct {
	rank, suit string
}

var (
	// string of rank/suit values in increasing value order
	rankValues, suitValues string
	// map of rank/suit to index
	ranks, suits map[string]int
)

func init() {
	rankValues = "23456789TJQKA"
	ranks = make(map[string]int, len(rankValues))
	for ii, r := range rankValues {
		ranks[string(r)] = ii
	}
	suitValues = "CSDH"
	suits = make(map[string]int, len(suitValues))
	for ii, s := range suitValues {
		suits[string(s)] = ii
	}
}

func NewCard(name string) (c *Card, err bool) {
	// search for suit and rank in letters of name.
	// Order is not important
	var rank, suit string
	for _, let := range name {
		if _, found := ranks[string(let)]; found {
			rank = string(let)
		} else if _, found := suits[string(let)]; found {
			suit = string(let)
		}
	}
	c = new(Card)
	c.rank, c.suit = rank, suit
	return c, false
}

func (c Card) RankIndex() int {
	return ranks[c.rank]
}

func (c Card) SuitIndex() int {
	return suits[c.suit]
}

// Class methods
func Suits() []string {
	var suitAry []string
	for _, suit := range suitValues {
		suitAry = append(suitAry, string(suit))
	}
	return suitAry
}

func Ranks() []string {
	var rankAry []string
	for _, rank := range rankValues {
		rankAry = append(rankAry, string(rank))
	}
	return rankAry
}

func (c Card) Rank() string {
	return c.rank
}

func (c Card) Suit() string {
	return c.suit
}

func (c Card) String() string {
	s := fmt.Sprintf("%s%s", c.Rank(), c.Suit())
	return s
}
