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

func (g *Game) AddCardToDealer(ctx context.Context, telegramChatID int64, card core.Card) error {
	filter := bson.M{"telegram_chat_id": telegramChatID}
	update := bson.M{"$push": bson.M{"active_game.dealer_hand": card}}
	if err := g.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

func (g *Game) AddCardToPlayerHand(ctx context.Context, telegramChatID int64, username string, card core.Card) error {
	filter := bson.M{"$and": bson.A{
		bson.M{"telegram_chat_id": telegramChatID},
		bson.M{"active_game.players_hands.username": username},
	}}
	update := bson.M{"$push": bson.M{"active_game.players_hands.$.cards": card}}
	if err := g.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

func (g *Game) AddNewPlayerHand(ctx context.Context, telegramChatID int64, playerHand core.PlayerHand) error {
	filter := bson.M{"telegram_chat_id": telegramChatID}
	update := bson.M{"$push": bson.M{"active_game.players_hands": playerHand}}
	if err := g.collection.FindOneAndUpdate(ctx, filter, update).Err(); err == nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

// todo test for concurrency
func (g *Game) DrawCard(ctx context.Context, telegramChatID int64) (core.Card, error) {
	var card core.Card

	filter := bson.M{"telegram_chat_id": telegramChatID}
	update := bson.M{"$pop": bson.M{"active_game.deck": 1}}

	res := g.collection.FindOneAndUpdate(ctx, filter, update)
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

	return chat.ActiveGame.Deck.DrawCard()
}
