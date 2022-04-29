package core

import (
	"log"
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

type Card string

type Cards []Card

func NewCards() Cards {
	cards := Cards{}

	for v := range CardValues {
		for _, s := range CardSymbols {
			cards = append(cards, Card(s+v))
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })

	return cards
}

func (cards Cards) CountValue() int {
	var value int

	for _, card := range cards {
		value += card.GetValue()
	}

	return value
}

func (c Card) GetValue() int {
	value, ok := CardValues[trimLeftChars(string(c), 1)]
	if !ok {
		log.Printf("wrong value for card %s", c)
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
