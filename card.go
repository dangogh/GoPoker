package poker

type Card int

var (
	ranks = "23456789TJQKA"
	suits = "HDSC"
)

func (c Card) Rank() byte {
	return ranks[int(c)%len(ranks)]
}

func (c Card) Suit() byte {
	return suits[int(c)/len(ranks)]
}

func (c Card) Deck() int {
	return int(c) / (len(ranks) * len(suits))
}
