//package card

package poker

import "fmt"
import "time"
import "math/rand"

type Deck []Card

func Shuffle() {
	rand.Seed(time.Now().Unix())
	for ii, r := range rand.Perm(52) {
		fmt.Printf("%d is %d\n", ii, r)
	}
}

// implements Sort interface
func (self Deck) Len() int {
	return len(self)
}

func (self Deck) Less(ii, jj int) bool {
	return self[ii].Rank() < self[jj].Rank()
}

func (self Deck) Swap(ii, jj int) {
	self[ii], self[jj] = self[jj], self[ii]
}
