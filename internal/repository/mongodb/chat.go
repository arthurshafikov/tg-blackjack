package mongodb

import (
	"context"
	"errors"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Chat struct {
	collection *mongo.Collection
}

func NewChat(collection *mongo.Collection) *Chat {
	return &Chat{
		collection: collection,
	}
}

func (c *Chat) CheckChatExists(ctx context.Context, telegramChatID int64) error {
	filter := bson.M{core.TelegramChatIDField: telegramChatID}
	if err := c.collection.FindOne(ctx, filter).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

func (c *Chat) RegisterChat(ctx context.Context, telegramChatID int64) error {
	filter := bson.M{
		core.TelegramChatIDField: telegramChatID,
	}
	chat := core.Chat{
		TelegramChatID: telegramChatID,
		Statistics:     core.UsersStatistics{},
		Deck:           core.NewDeckWithCards(core.Cards{}),
	}
	update := bson.M{
		"$setOnInsert": chat,
	}
	opts := options.Update().SetUpsert(true)

	_, err := c.collection.UpdateOne(ctx, filter, update, opts)

	return err
}
