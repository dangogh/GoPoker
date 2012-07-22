package poker

import "testing"
import "fmt"

func TestCard(t *testing.T) {
	c := Card(0)
	if c.Rank() != '2' || c.Suit() != 'H' {
		t.Errorf("Expected 2C, got %c%c", c.Rank(), c.Suit())
	}
}

func ExampleCard() {
	for ii := 0; ii < 52; ii++ {
		c := Card(ii)
		rank, suit := c.Rank(), c.Suit()
		fmt.Printf("%c%c\n", rank, suit)
	}
	// Output:
	// 2H
	// 3H
	// 4H
	// 5H
	// 6H
	// 7H
	// 8H
	// 9H
	// TH
	// JH
	// QH
	// KH
	// AH
	// 2D
	// 3D
	// 4D
	// 5D
	// 6D
	// 7D
	// 8D
	// 9D
	// TD
	// JD
	// QD
	// KD
	// AD
	// 2S
	// 3S
	// 4S
	// 5S
	// 6S
	// 7S
	// 8S
	// 9S
	// TS
	// JS
	// QS
	// KS
	// AS
	// 2C
	// 3C
	// 4C
	// 5C
	// 6C
	// 7C
	// 8C
	// 9C
	// TC
	// JC
	// QC
	// KC
	// AC
}
