package poker
type Hand []Card

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
	return high - low + 1 == len(h)
}

func (h Hand) IsFlush() bool {
	suit := h[0].Suit()
	for _, c := range h {
		val := c.Suit()
		if val != suit {
			return false
		}
	}
	return true
}

func (h Hand) HighCard() Card {
	return h[0]
}
