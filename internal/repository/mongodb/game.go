package mongodb

import (
	"context"
	"errors"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Game struct {
	collection *mongo.Collection
}

func NewGame(collection *mongo.Collection) *Game {
	return &Game{
		collection: collection,
	}
}

func (g *Game) SetActiveGame(ctx context.Context, telegramChatID int64, game core.Game) error {
	filter := bson.M{"$and": bson.A{
		bson.M{"telegram_chat_id": telegramChatID},
		bson.M{"active_game": nil},
	}}
	if err := g.collection.FindOne(ctx, filter).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrActiveGame
		}

		return err
	}

	update := bson.M{"$set": bson.M{"active_game": game}}
	err := g.collection.FindOneAndUpdate(ctx, filter, update).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return core.ErrNotFound
	}

	return err
}

func (g *Game) FinishActiveGame(ctx context.Context, telegramChatID int64) (core.Game, error) {
	var chat core.Chat

	filter := bson.M{"telegram_chat_id": telegramChatID}
	update := bson.M{"$set": bson.M{"active_game": nil}}
	res := g.collection.FindOneAndUpdate(ctx, filter, update)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return chat.ActiveGame, core.ErrNotFound
		}
		return chat.ActiveGame, err
	}

	if err := res.Decode(&chat); err != nil {
		return chat.ActiveGame, err
	}

	return chat.ActiveGame, nil
}

func (g *Game) AddCardToPlayerHand(ctx context.Context, telegramChatID int64, username string, card string) error {
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

func (g *Game) AddCardToDealer(ctx context.Context, telegramChatID int64, card string) error {
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
