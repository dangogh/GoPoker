package poker

/*
In the card game poker, a hand consists of five cards and are ranked, from
lowest to highest, in the following way:

    High Card: Highest value card.
    One Pair: Two cards of the same value.
    Two Pairs: Two different pairs.
    Three of a Kind: Three cards of the same value.
    Straight: All cards are consecutive values.
    Flush: All cards of the same suit.
    Full House: Three of a kind and a pair.
    Four of a Kind: Four cards of the same value.
    Straight Flush: All cards are consecutive values of same suit.
    Royal Flush: Ten, Jack, Queen, King, Ace, in same suit.

The cards are valued in the order:
2, 3, 4, 5, 6, 7, 8, 9, 10, Jack, Queen, King, Ace.
*/

import (
	"fmt"
	"sort"
	"strings"
)

type Hand struct {
	cards []Card
}

// implement sort interface
func (h Hand) Len() int {
	return len(h.cards)
}

// sort cards by significance for the hand type
func (h Hand) Less(i, j int) bool {
	a, b := h.cards[i], h.cards[j]
	arank, brank := a.RankIndex(), b.RankIndex()
	retval := true
	if arank == brank {
		retval = a.SuitIndex() < b.SuitIndex()
	} else {
		// count number of a, b ranks in this hand
		acount, bcount := 0, 0
		for _, c := range h.cards {
			if a.Rank() == c.Rank() {
				acount++
			}
			if b.Rank() == c.Rank() {
				bcount++
			}
		}
		if acount == bcount {
			// same -- compare rank
			retval = arank < brank
		} else {
			// pick higher count
			retval = acount > bcount
		}
	}
	return retval
}

func (h Hand) Swap(i, j int) {
	h.cards[i], h.cards[j] = h.cards[j], h.cards[i]
}

// Hand ranks
const (
	HighCard = iota
	Pair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

////////////////////////////////////////////////////////////////
// Helper methods
// helper method -- return map of ranks to counts
func (h Hand) countRanks() map[rune][]Card {
	ranks := make(map[rune][]Card)
	for _, c := range h.cards {
		rank := c.Rank()
		ranks[rank] = append(ranks[rank], c)
	}
	return ranks
}

// helper method -- check if n of same rank
func (h Hand) isNOfAKind(n int) bool {
	isNOfAKind := false
	for _, ary := range h.countRanks() {
		if len(ary) == n {
			isNOfAKind = true
			break
		}
	}
	return isNOfAKind
}

////////////////////////////////////////////////////////////////
// Exported methods
// IsStraight returns bool and high card
func (h Hand) IsStraight() bool {
	var ranks []int
	for _, c := range h.cards {
		ranks = append(ranks, c.RankIndex())
	}
	sort.Ints(ranks)
	isStraight := true
	prev, ranks := ranks[0], ranks[1:]
	for _, r := range ranks {
		if prev+1 != r {
			isStraight = false
			break
		}
		prev = r
	}
	return isStraight
}

func (h Hand) IsFlush() bool {
	suit := h.cards[0].Suit()
	isFlush := true

	for _, c := range h.cards {
		val := c.Suit()
		if val != suit {
			isFlush = false
			break
		}
	}
	return isFlush
}

func (h Hand) IsRoyalFlush() bool {
	return h.HighCard().Rank() == 'A' && h.IsStraightFlush()
}

func (h Hand) IsStraightFlush() bool {
	return h.IsStraight() && h.IsFlush()
}

func (h Hand) IsFourOfAKind() bool {
	return h.isNOfAKind(4)
}

func (h Hand) IsFullHouse() bool {
	foundTwo, foundThree := false, false
	for _, ary := range h.countRanks() {
		switch len(ary) {
		case 3:
			foundThree = true
		case 2:
			foundTwo = true
		}
	}
	return foundThree && foundTwo
}

func (h Hand) IsTwoPair() bool {
	numPairs := 0
	for _, ary := range h.countRanks() {
		if len(ary) == 2 {
			numPairs++
		}
	}
	return numPairs == 2
}

func (h Hand) IsThreeOfAKind() bool {
	return h.isNOfAKind(3)
}
func (h Hand) IsPair() bool {
	return h.isNOfAKind(2)
}

func (h Hand) HighCard() Card {
	return h.cards[0]
}

func (h Hand) Rank() int {
	rank := HighCard
	switch {
	case h.IsRoyalFlush():
		rank = RoyalFlush
	case h.IsStraightFlush():
		rank = StraightFlush
	case h.IsFourOfAKind():
		rank = FourOfAKind
	case h.IsFullHouse():
		rank = FullHouse
	case h.IsFlush():
		rank = Flush
	case h.IsStraight():
		rank = Straight
	case h.IsThreeOfAKind():
		rank = ThreeOfAKind
	case h.IsTwoPair():
		rank = TwoPair
	case h.IsPair():
		rank = Pair
	default:
		rank = HighCard
	}
	return rank
}

func (h Hand) Beats(other Hand) bool {
	rank := h.Rank()
	otherRank := other.Rank()
	wins := false
	switch {
	case rank > otherRank:
		wins = true
	case otherRank > rank:
		wins = false
	default:
		{
			// else pick highest ranking card
			for ii := 0; ii < len(h.cards); ii++ {
				switch {
				case h.cards[ii].Rank() > other.cards[ii].Rank():
					wins = true
				case h.cards[ii].Rank() < other.cards[ii].Rank():
					wins = false
				default:
					continue
				}
				// determined..
				break
			}
		}
	}
	return wins
}

func (h Hand) String() string {
	s := make([]string, len(h.cards))
	for _, c := range h.cards {
		cs := string(c)
		s = append(s, string(c))
		fmt.Printf(" card is %s\n", c)
	}
	return strings.Join(s, " ")
}
