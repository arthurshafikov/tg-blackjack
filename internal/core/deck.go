package core

import "sync"

type Deck struct {
	mu    sync.Mutex
	Cards Cards `bson:"cards"`
}

func NewDeck(numOfDecks int) *Deck {
	return &Deck{
		Cards: NewCards(numOfDecks),
	}
}

func NewDeckWithCards(cards Cards) *Deck {
	return &Deck{
		Cards: cards,
	}
}

func (d *Deck) IsEmpty() bool {
	return len(d.Cards) < 1
}

func (d *Deck) DrawCard() (Card, error) {
	var card Card

	cards, err := d.drawCards(1)
	if err != nil {
		return card, err
	}
	card = cards[0]

	return card, nil
}

func (d *Deck) drawCards(amount int) (Cards, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var drawedCards Cards
	deckLength := len(d.Cards)
	if deckLength < 1 {
		return drawedCards, ErrDeckEmpty
	}

	if deckLength < amount {
		amount = deckLength
	}

	drawedCards = d.Cards[deckLength-amount:]
	d.Cards = d.Cards[:deckLength-amount]

	return drawedCards, nil
}
