package mongodb

import (
	"context"
	"errors"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Card struct {
	collection *mongo.Collection
}

func NewCard(collection *mongo.Collection) *Card {
	return &Card{
		collection: collection,
	}
}

func (c *Card) AddCardToDealer(ctx context.Context, telegramChatID int64, card core.Card) error {
	filter := bson.M{"telegram_chat_id": telegramChatID}
	update := bson.M{"$push": bson.M{"active_game.dealer": card}}

	if err := c.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

func (c *Card) AddCardToPlayer(ctx context.Context, telegramChatID int64, username string, card core.Card) error {
	filter := bson.M{"$and": bson.A{
		bson.M{"telegram_chat_id": telegramChatID},
		bson.M{"active_game.players.username": username},
	}}
	update := bson.M{"$push": bson.M{"active_game.players.$.cards": card}}

	if err := c.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

func (c *Card) DrawCard(ctx context.Context, telegramChatID int64) (core.Card, error) {
	var card core.Card

	filter := bson.M{"telegram_chat_id": telegramChatID}
	update := bson.M{"$pop": bson.M{"deck.cards": 1}}

	res := c.collection.FindOneAndUpdate(ctx, filter, update)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return card, core.ErrNotFound
		}

		return card, err
	}

	var chat core.Chat
	if err := res.Decode(&chat); err != nil {
		return card, err
	}

	if chat.Deck.IsEmpty() {
		var err error
		chat.Deck, err = c.setNewDeck(ctx, telegramChatID)
		if err != nil {
			return card, err
		}
	}

	return chat.Deck.DrawCard()
}

func (c *Card) DrawCards(ctx context.Context, telegramChatID int64, amount int) (core.Cards, error) {
	cards := core.Cards{}
	for amount > 0 {
		card, err := c.DrawCard(ctx, telegramChatID)
		if err != nil {
			return nil, err
		}

		cards = append(cards, card)

		amount--
	}

	return cards, nil
}

func (c *Card) setNewDeck(ctx context.Context, telegramChatID int64) (*core.Deck, error) {
	deck := core.NewDeck()
	filter := bson.M{"telegram_chat_id": telegramChatID}
	update := bson.M{"$set": bson.M{"deck": deck}}

	if err := c.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return deck, core.ErrNotFound
		}

		return deck, err
	}

	return deck, nil
}
