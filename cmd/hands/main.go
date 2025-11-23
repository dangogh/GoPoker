package main

import (
	"flag"
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

// performDraw takes current cards, asks hand.RecommendDiscards for up to maxDiscard indices,
// draws replacements from the deck and returns the updated cards, the cards that were discarded,
// and the cards that were drawn.
// maxDiscard is computed based on 5-card draw rules: if Ace is kept, allow 4 discards; else 3.
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

func printCards(cs []cards.Card) {
	for i, c := range cs {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(c.String())
	}
	fmt.Println()
}

func run(players int) error {
	if players <= 0 {
		return fmt.Errorf("players must be > 0")
	}

	d := deck.NewDeck()
	d.Shuffle()

	// initial deal
	hands := make([][]cards.Card, players)
	for i := 0; i < players; i++ {
		cs, err := d.Deal(5)
		if err != nil {
			return fmt.Errorf("deal error: %w", err)
		}
		hands[i] = cs
	}

	fmt.Println("Initial hands:")
	for i := 0; i < players; i++ {
		fmt.Printf("Player %d:\n", i+1)
		printCards(hands[i])
	}

	// Draw phase for each player
	for i := 0; i < players; i++ {
		// Compute max discard based on 5-card draw rules
		maxDisc := hand.ComputeMaxDiscard(hand.Hand{Cards: hands[i]})

		cs, discarded, drew, err := performDraw(d, hands[i], maxDisc)
		if err != nil {
			return fmt.Errorf("draw error: %w", err)
		}
		hands[i] = cs
		if len(discarded) > 0 {
			fmt.Printf("Player %d discarded: ", i+1)
			for j, c := range discarded {
				if j > 0 {
					fmt.Print(" ")
				}
				fmt.Print(c.String())
			}
			fmt.Print(" and drew: ")
			for j, c := range drew {
				if j > 0 {
					fmt.Print(" ")
				}
				fmt.Print(c.String())
			}
			fmt.Println()
		} else {
			fmt.Printf("Player %d stood pat.\n", i+1)
		}
	}

	// Evaluate final hands and find winner(s)
	evals := make([]hand.EvaluatedHand, players)
	for i := range players {
		h := hand.Hand{Cards: hands[i]}
		evals[i] = hand.Evaluate(h)
	}

	fmt.Println()
	fmt.Println("Final hands:")
	for i := 0; i < players; i++ {
		fmt.Printf("Player %d: %s\n", i+1, categoryName(evals[i].Category))
		printCards(hands[i])
	}

	// determine best hand(s)
	bestIdxs := []int{0}
	for i := 1; i < players; i++ {
		cmp := hand.Compare(evals[i], evals[bestIdxs[0]])
		if cmp > 0 {
			bestIdxs = []int{i}
		} else if cmp == 0 {
			bestIdxs = append(bestIdxs, i)
		}
	}

	if len(bestIdxs) == 1 {
		fmt.Printf("Winner: Player %d\n", bestIdxs[0]+1)
	} else {
		fmt.Printf("Result: Tie among players")
		for _, idx := range bestIdxs {
			fmt.Printf(" %d", idx+1)
		}
		fmt.Println()
	}

	return nil
}

func main() {
	players := flag.Int("players", 5, "number of players (each dealt 5 cards)")
	flag.Parse()

	if err := run(*players); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
