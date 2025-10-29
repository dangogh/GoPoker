package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/dangogh/GoPoker/cards"
	"github.com/dangogh/GoPoker/deck"
	"github.com/dangogh/GoPoker/hand"
)

func categoryName(c hand.Category) string {
	switch c {
	case hand.HighCard:
		return "High Card"
	case hand.OnePair:
		return "One Pair"
	case hand.TwoPair:
		return "Two Pair"
	case hand.ThreeOfKind:
		return "Three of a Kind"
	case hand.Straight:
		return "Straight"
	case hand.Flush:
		return "Flush"
	case hand.FullHouse:
		return "Full House"
	case hand.FourOfKind:
		return "Four of a Kind"
	case hand.StraightFlush:
		return "Straight Flush"
	default:
		return "Unknown"
	}
}

func printCards(cs []cards.Card) {
	for i, c := range cs {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(c.String())
	}
	fmt.Println()
}

// performDraw takes current cards, asks hand.RecommendDiscards for up to maxDiscard indices,
// draws replacements from the deck and returns the updated cards, the cards that were discarded,
// and the cards that were drawn.
func performDraw(d *deck.Deck, cs []cards.Card, maxDiscard int) ([]cards.Card, []cards.Card, []cards.Card, error) {
	discardIdxs := hand.RecommendDiscards(hand.Hand{Cards: cs}, maxDiscard)
	if len(discardIdxs) == 0 {
		return cs, nil, nil, nil
	}
	// ensure deterministic mapping: sort discard indices ascending
	sort.Ints(discardIdxs)

	// collect the cards being discarded
	discarded := make([]cards.Card, len(discardIdxs))
	for i, idx := range discardIdxs {
		discarded[i] = cs[idx]
	}

	// draw replacements
	repl, err := d.Deal(len(discardIdxs))
	if err != nil {
		return cs, nil, nil, err
	}
	// apply replacements in index order
	for i, idx := range discardIdxs {
		cs[idx] = repl[i]
	}

	return cs, discarded, repl, nil
}

func main() {
	d := deck.NewDeck()
	d.Shuffle()

	h1cs, err := d.Deal(5)
	if err != nil {
		fmt.Fprintln(os.Stderr, "deal error:", err)
		return
	}
	h2cs, err := d.Deal(5)
	if err != nil {
		fmt.Fprintln(os.Stderr, "deal error:", err)
		return
	}

	fmt.Println("Initial hands:")
	fmt.Println("Player 1:")
	printCards(h1cs)
	fmt.Println("Player 2:")
	printCards(h2cs)

	// Draw phase: up to 3 cards
	h1cs, disc1, drew1, err := performDraw(d, h1cs, 3)
	if err != nil {
		fmt.Fprintln(os.Stderr, "draw error:", err)
		return
	}
	if len(disc1) > 0 {
		fmt.Print("Player 1 discarded: ")
		for i, c := range disc1 {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(c.String())
		}
		fmt.Print(" and drew: ")
		for i, c := range drew1 {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(c.String())
		}
		fmt.Println()
	} else {
		fmt.Println("Player 1 stood pat.")
	}

	h2cs, disc2, drew2, err := performDraw(d, h2cs, 3)
	if err != nil {
		fmt.Fprintln(os.Stderr, "draw error:", err)
		return
	}
	if len(disc2) > 0 {
		fmt.Print("Player 2 discarded: ")
		for i, c := range disc2 {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(c.String())
		}
		fmt.Print(" and drew: ")
		for i, c := range drew2 {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(c.String())
		}
		fmt.Println()
	} else {
		fmt.Println("Player 2 stood pat.")
	}

	// Evaluate final hands
	h1 := hand.Hand{Cards: h1cs}
	h2 := hand.Hand{Cards: h2cs}
	e1 := hand.Evaluate(h1)
	e2 := hand.Evaluate(h2)

	fmt.Println()
	fmt.Println("Final hands:")
	fmt.Println("Player 1:", categoryName(e1.Category))
	printCards(h1cs)
	fmt.Println("Player 2:", categoryName(e2.Category))
	printCards(h2cs)

	switch hand.Compare(e1, e2) {
	case 1:
		fmt.Println("Winner: Player 1")
	case -1:
		fmt.Println("Winner: Player 2")
	default:
		fmt.Println("Result: Tie")
	}
}
