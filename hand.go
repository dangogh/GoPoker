package poker
type Hand []Card

////////////////////////////////////////////////////////////////
// Helper methods
// helper method -- return map of ranks to counts
func (h Hand) countRanks() map[rune]int {
	ranks := make(map[rune]int)
	for _, c := range h {
		ranks[c.Rank()]++
	}
	return ranks
}

// helper method -- check if n of same rank
func (h Hand) isNOfAKind(n int) bool {
	isNOfAKind := false
	for _, c := range h.countRanks() {
		if c >= n {
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
	low, high := 999, -1
	for _, c := range h {
		val := c.RankIndex()
		if val < low {
			low = val
		}
		if val > high {
			high = val
		}
	}
	isStraight := (high - low + 1 == len(h))
	return isStraight
}

func (h Hand) IsFlush() bool {
	suit := h[0].Suit()
	isFlush := true
	
	for _, c := range h {
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
	for _, count := range h.countRanks() {
		switch count {
			case 3 : foundThree = true
			case 2 : foundTwo = true
		}
	}
	return foundThree && foundTwo
}

func (h Hand) IsTwoPair() bool {
	numPairs := 0
	for _, count := range h.countRanks() {
		if count >= 2 {
			numPairs++
		}
	}
	return numPairs >= 2
}

func (h Hand) IsThreeOfAKind() bool {
	return h.isNOfAKind(3)
}
func (h Hand) IsTwoOfAKind() bool {
	return h.isNOfAKind(2)
}

func (h Hand) HighCard() Card {
	return h[0]
}

