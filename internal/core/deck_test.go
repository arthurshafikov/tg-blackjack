package core

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDrawCards(t *testing.T) {
	cards := Cards{"♣5", "♦Q", "♥K", "♠J", "♥10", "♠2", "♥A", "♠4"}
	deck := Deck{
		Cards: cards,
	}

	result, err := deck.drawCards(4)
	require.NoError(t, err)
	require.Equal(t, cards[len(cards)-4:], result)
	require.Equal(t, cards[:len(cards)-4], deck.Cards)
}

func TestDrawCardsMoreThanAmount(t *testing.T) {
	cards := Cards{"♣5", "♦Q"}
	deck := Deck{
		Cards: cards,
	}

	result, err := deck.drawCards(4)
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, cards, result)
	require.Equal(t, Cards{}, deck.Cards)
}

func TestDrawCardsEmptyDeck(t *testing.T) {
	var expected Cards
	deck := Deck{
		Cards: Cards{},
	}

	result, err := deck.drawCards(4)
	require.ErrorIs(t, err, ErrDeckEmpty)
	require.Equal(t, expected, result)
}

func TestConcurrencyDrawCard(t *testing.T) {
	cards := NewCards(6)
	deckLength := len(cards)
	deck := Deck{
		Cards: cards,
	}
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 5; j++ {
				_, err := deck.DrawCard()
				require.NoError(t, err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	require.Len(t, deck.Cards, deckLength-50)
}
