package main

import (
	"fmt"
	"os"

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

	h1 := hand.Hand{Cards: h1cs}
	h2 := hand.Hand{Cards: h2cs}
	e1 := hand.Evaluate(h1)
	e2 := hand.Evaluate(h2)

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
