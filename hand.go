package poker

type Hand []Card

// implements Sort interface
func (self Hand) Len() int {
	return len(self)
}

func (self Hand) Less(ii, jj int) bool {
	return self[ii].Rank() < self[jj].Rank()
}

func (self Hand) Swap(ii, jj int) {
	self[ii], self[jj] = self[jj], self[ii]
}
