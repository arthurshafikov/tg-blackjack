package core

import (
	"math/rand"
	"time"
)

var CardValues = map[string]int{
	"A":  1,
	"2":  2,
	"3":  3,
	"4":  4,
	"5":  5,
	"6":  6,
	"7":  7,
	"8":  8,
	"9":  9,
	"10": 10,
	"J":  10,
	"Q":  10,
	"K":  10,
}

var CardSymbols = []string{
	"♣",
	"♦",
	"♥",
	"♠",
}

const NumOfCardsInDeck = 52

type Card string

type Cards []Card

func NewCards(numOfDecks int) Cards {
	cards := Cards{}

	for i := 0; i < numOfDecks; i++ {
		for v := range CardValues {
			for _, s := range CardSymbols {
				cards = append(cards, Card(s+v))
			}
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })

	return cards
}

func (cards Cards) CountValue() int {
	possibleValues := []int{0}

	for _, card := range cards {
		cardLetter := trimLeftChars(string(card), 1)
		for i := range possibleValues {
			possibleValues[i] += CardValues[cardLetter]
		}
		if cardLetter == "A" {
			copyPossibleValues := make([]int, len(possibleValues))
			copy(copyPossibleValues, possibleValues)

			for i := range copyPossibleValues {
				copyPossibleValues[i] += 10
			}

			possibleValues = append(possibleValues, copyPossibleValues...)
		}
	}

	bestValue := possibleValues[0]
	for _, value := range possibleValues {
		if value > bestValue && value <= 21 {
			bestValue = value
		}
	}

	return bestValue
}

func (cards Cards) IsBlackjack() bool {
	return cards.CountValue() == 21 && len(cards) == 2
}

func (c Card) ToString() string {
	return "*" + string(c) + "*"
}

func trimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}

	return s[:0]
}
