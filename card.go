package poker

import (
	"strings"
	"errors"
	"fmt"
	)

type Card struct {
	name       string
	rank, suit rune
}

var (
	ranks = "23456789TJQKA"
	suits = "HDSC"
)

func Suits() []rune {
  var runes []rune
  for _, r := range strings.Split(suits, "") {
  	runes = append(runes, rune(r[0]))
  }
  return runes
}

func Ranks() []rune {
  var runes []rune
  for _, r := range strings.Split(ranks, "") {
  	runes = append(runes, rune(r[0]))
  }
  return runes
}

func (c Card) extrude() (err error) {
	isuit, irank := -1, -1
	if len(c.name) == 2 {
		for _, ch := range c.name {
			if isuit != -1 && irank != -1 {
				break
			}
			if isuit == -1 {
				isuit = strings.Index(suits, string(ch))
				fmt.Printf("isuit is %v\n", isuit)
			}
			if irank == -1 {
				irank = strings.Index(ranks, string(ch))
				fmt.Printf("irank is %v\n", irank)
			}
		}
	}
	if irank == -1 || isuit == -1 {
		return errors.New("Error creating card from " + c.name)
	}

	c.rank, c.suit = rune(ranks[irank]), rune(suits[isuit])
	return nil
}

func (c Card) Rank() rune {
	if c.rank == 0 {
		c.extrude()
	}
	return c.rank
}

func (c Card) Suit() rune {
	if c.suit == 0 {
		c.extrude()
	}
	return c.suit
}

func (c Card) String() string {
	ret := string(c.Rank()) + string(c.Suit())
	return ret
}
