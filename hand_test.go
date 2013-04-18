package poker

import "testing"
import "fmt"

type THand struct {
	cards         [5]string
	highcard      string
	pair          bool
	twopair       bool
	threeofakind  bool
	straight      bool
	flush         bool
	fullhouse     bool
	fourofakind   bool
	straightflush bool
	royalflush    bool
}

var hands = [...]THand{
	// pair
	{cards: [...]string{"2C", "2D", "4C", "3C", "5C"},
		highcard: "6C",
		pair:     true,
	},
	// two pair
	{cards: [...]string{"6D", "TS", "7C", "7D", "TC"},
		highcard: "6C",
		twopair:  true,
		pair:     true,
	},
	// three of a kind
	{cards: [...]string{"8H", "2S", "8C", "8D", "TC"},
		highcard:     "6C",
		threeofakind: true,
	},
	// four of a kind
	{cards: [...]string{"5H", "5S", "5C", "5D", "KC"},
		highcard:    "6C",
		fourofakind: true,
	},
	// straight flush
	{cards: [...]string{"2C", "6C", "4C", "3C", "5C"},
		highcard:      "6C",
		straight:      true,
		flush:         true,
		straightflush: true},
	// flush
	{cards: [...]string{"2D", "6D", "TD", "3D", "AD"},
		highcard: "AD",
		flush:    true,
	},
}

func testHand(t *testing.T, h Hand, expected, got bool, handType string) {
	if got && !expected {
		t.Errorf("%v should not be %s", h, handType)
	}
	if !got && expected {
		t.Errorf("%v should be %s", h, handType)
	}
}

func TestHands(t *testing.T) {
	for _, h := range hands {
		var cards []Card
		for _, c := range h.cards {
			card, _ := NewCard(c)
			cards = append(cards, *card)
		}
		hand := Hand{cards}
		fmt.Printf("Hand is %s\n", hand)
		if hand.Len() != 5 {
			t.Errorf("Expected hand length 5, got %d\n",
				hand.Len())
		}
		testHand(t, hand, h.pair, hand.IsPair(), "pair")
		testHand(t, hand, h.twopair, hand.IsTwoPair(), "twopair")
		testHand(t, hand, h.threeofakind, hand.IsThreeOfAKind(), "threeofakind")
		testHand(t, hand, h.straight, hand.IsStraight(), "straight")
		testHand(t, hand, h.flush, hand.IsFlush(), "flush")
		testHand(t, hand, h.fullhouse, hand.IsFullHouse(), "fullhouse")
		testHand(t, hand, h.fourofakind, hand.IsFourOfAKind(), "fourofakind")
		testHand(t, hand, h.straightflush, hand.IsStraightFlush(), "straightflush")
		testHand(t, hand, h.royalflush, hand.IsRoyalFlush(), "royalflush")
	}
}
