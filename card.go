package poker

type Card struct {
	rank, suit rune
}

var (
	// string of rank/suit values in increasing value order
	rankValues, suitValues string
	// map of rank/suit to index
	ranks, suits map[rune]int
    )

func init() {
	rankValues = "23456789TJQKA"
	ranks = make(map[rune]int, len(rankValues))
	for ii, r := range rankValues {
		ranks[r] = ii
	}
	suitValues = "CSDH"
	suits = make(map[rune]int, len(suitValues))
	for ii, s := range suitValues {
		suits[s] = ii
	}
}

func NewCard(name string) (c *Card, err bool) {
	// search for suit and rank in letters of name.
	// Order is not important
	var rank, suit rune
	for _, let := range name {
		if _, found := ranks[let]; found {
			rank = let
		} else if _, found := suits[let]; found {
			suit = let
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
func Suits() []rune {
	var suitAry []rune
	for _, suit := range suitValues {
		suitAry = append(suitAry, suit)
	}
	return suitAry
}

func Ranks() []rune {
	var rankAry []rune
	for _, rank := range rankValues {
		rankAry = append(rankAry, rank)
	}
	return rankAry
}

func (c Card) Rank() rune {
	return c.rank
}

func (c Card) Suit() rune {
	return c.suit
}

func (c Card) String() string {
	return string(c.Rank()) + string(c.Suit())
}
