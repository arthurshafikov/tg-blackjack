package tests

import (
	"fmt"
	"sync"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *APITestSuite) TestConcurrentDrawFromDeck() {
	err := s.services.Chats.RegisterChat(s.ctx, telegramChatID)
	r.NoError(err)

	_, err = s.services.Games.NewGame(s.ctx, telegramChatID)
	r.NoError(err)

	deck := s.getDeck()
	initLength := 6*core.NumOfCardsInDeck - 2 // dealer has taken 2 cards
	r.Len(deck.Cards, initLength)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			for j := 0; j < 10; j++ {
				s.playerEntersTheGame(fmt.Sprintf("player%v%v", i, j))
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	deck = s.getDeck()
	r.Len(deck.Cards, initLength-200)
}

func (s *APITestSuite) TestExactCardsFromDeck() {
	err := s.services.Chats.RegisterChat(s.ctx, telegramChatID)
	r.NoError(err)

	_, err = s.services.Games.NewGame(s.ctx, telegramChatID)
	r.NoError(err)

	deck := s.getDeck()
	expected := deck.Cards[len(deck.Cards)-2:] // cards that player should draw

	player := s.playerEntersTheGame(username)
	r.ElementsMatch(expected, player.Cards)
}

func (s *APITestSuite) getDeck() *core.Deck {
	var chat core.Chat
	res := s.collection.FindOne(s.ctx, bson.M{core.TelegramChatIDField: telegramChatID})
	r.NoError(res.Err())
	r.NoError(res.Decode(&chat))

	return chat.Deck
}
