package core

type Deck struct {
	cards Cards
}

func NewDeck() *Deck {
	return &Deck{
		cards: NewCards(),
	}
}

func (d *Deck) IsEmpty() bool {
	return len(d.cards) < 1
}

// todo test for concurrency
func (d *Deck) DrawCards(amount int) (Cards, error) {
	var drawedCards Cards
	deckLength := len(d.cards)
	if deckLength < 1 {
		return drawedCards, ErrDeckEmpty
	}

	if deckLength < amount {
		amount = deckLength
	}

	drawedCards = d.cards[deckLength-amount:]
	d.cards = d.cards[:deckLength-amount]

	return drawedCards, nil
}

func (d *Deck) DrawCard() (Card, error) {
	var card Card

	cards, err := d.DrawCards(1)
	if err != nil {
		return card, nil
	}
	card = cards[len(cards)]

	return card, nil
}
