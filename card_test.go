package poker

import (
	"testing"
	"fmt"
	)

func TestCard(t *testing.T) {
	//c := Card{name: "KS"}
	for _, suit := range "CSDH" {
		for _, rank := range "23456789TJQKA" {
			c := Card{suit: rune(suit), rank: rune(rank)}
			if c.Rank() != rank || c.Suit() != suit {
				t.Errorf("Expected KS, got %c%c", c.Rank(), c.Suit())
			}
		}
	}
}

func ExampleCard() {
	for _, suit := range "CSDH" {
		for _, rank := range "23456789TJQKA" {
			c := Card{suit: rune(suit), rank: rune(rank)}
			fmt.Printf("%s\n", c);
		}
	}
	// Output:
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
}
