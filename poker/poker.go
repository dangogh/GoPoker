package main

import (
	"io"
	"os"
	"bufio"
	"fmt"
	"log"
	"poker"
       )

func dealCards(line string) []Hand {
	// split line into cards
	hands := make([]Hand)
	hands = append(hands, hand := Hand{})
	for ii, c := strings.split(line, " ") {
		hand = append(hand, card.NewCard(c))
		if ii % 5 == 4 {
			hands = append(hands, hand := Hand{})
		}
	}
	return hands
}

func pickWinner([]Hand hands) int {
	best := card.HighCard
	for ii, hand := range hands {
		switch 
	}
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	winners := make([]int)
	for {
		line, _, lineerr := reader.ReadLine()
		if lineerr != nil {
			if lineerr != io.EOF {
				log.Fatal(lineerr)
			}
			break
		}

		hands := dealCards(line)
		// Pick the winning hand
		winners = append(winners, pickWinner(hands))
	}
}
