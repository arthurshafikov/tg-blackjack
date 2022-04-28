package repository

import (
	"context"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Chats interface {
	CheckChatExists(ctx context.Context, telegramChatID int64) error
	RegisterChat(ctx context.Context, telegramChatID int64) error
	GetStatistics(ctx context.Context, telegramChatID int64) (core.UsersStatistics, error)
}

type Repository struct {
	Chats
}

func NewRepository(db *mongo.Client) *Repository {
	collection := db.Database("homestead").Collection("chats")

	return &Repository{
		Chats: mongodb.NewChat(collection),
	}
}
