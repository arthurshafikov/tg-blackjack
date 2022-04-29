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

type Games interface {
	SetActiveGame(ctx context.Context, telegramChatID int64, game core.Game) error
	FinishActiveGame(ctx context.Context, telegramChatID int64) (core.Game, error)
	AddCardToPlayerHand(ctx context.Context, telegramChatID int64, username string, card core.Card) error
	AddNewPlayerHand(ctx context.Context, telegramChatID int64, playerHand core.PlayerHand) error
	AddCardToDealer(ctx context.Context, telegramChatID int64, card core.Card) error
}

type Repository struct {
	Chats
	Games
}

func NewRepository(db *mongo.Client) *Repository {
	chatsCollection := db.Database("homestead").Collection("chats")

	return &Repository{
		Chats: mongodb.NewChat(chatsCollection),
		Games: mongodb.NewGame(chatsCollection),
	}
}
