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

type Cards []string

func NewCards() Cards {
	cards := Cards{}

	for v := range CardValues {
		for _, s := range CardSymbols {
			cards = append(cards, s+v)
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })

	return cards
}

func (cards Cards) CountValue() int {
	var value int

	for _, v := range cards {
		value += CardValues[trimLeftChars(v, 1)]
	}

	return value
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
