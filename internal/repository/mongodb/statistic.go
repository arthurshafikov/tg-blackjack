package mongodb

import (
	"context"
	"errors"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Statistic struct {
	collection *mongo.Collection
}

func NewStatistic(collection *mongo.Collection) *Statistic {
	return &Statistic{
		collection: collection,
	}
}

func (s *Statistic) GetStatistics(ctx context.Context, telegramChatID int64) (core.UsersStatistics, error) {
	filter := bson.M{"telegram_chat_id": telegramChatID}
	res := s.collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, core.ErrNotFound
		}

		return nil, err
	}

	var chat core.Chat
	if err := res.Decode(&chat); err != nil {
		return nil, err
	}

	return chat.Statistics, nil
}

func (s *Statistic) SetStatistics(
	ctx context.Context,
	telegramChatID int64,
	stats core.UsersStatistics,
) error {
	filter := bson.M{"telegram_chat_id": telegramChatID}
	update := bson.M{"statistics": stats}
	res := s.collection.FindOneAndUpdate(ctx, filter, update)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}
