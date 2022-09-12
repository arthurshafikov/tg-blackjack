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

func (g *Game) GetActiveGame(ctx context.Context, telegramChatID int64) (*core.Game, error) {
	var chat core.Chat

	filter := bson.M{core.TelegramChatIDField: telegramChatID}
	res := g.collection.FindOne(ctx, filter)
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

func (g *Game) NullActiveGame(ctx context.Context, telegramChatID int64) error {
	filter := bson.M{core.TelegramChatIDField: telegramChatID}
	update := bson.M{"$set": bson.M{core.ActiveGameField: nil}}
	if err := g.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

func (g *Game) SetActiveGame(ctx context.Context, telegramChatID int64, game core.Game) error {
	var chat core.Chat

	filter := bson.M{core.TelegramChatIDField: telegramChatID}
	res := g.collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	if err := res.Decode(&chat); err != nil {
		return err
	}

	if chat.ActiveGame != nil {
		return core.ErrActiveGame
	}

	update := bson.M{"$set": bson.M{core.ActiveGameField: game}}

	return g.collection.FindOneAndUpdate(ctx, filter, update).Err()
}
