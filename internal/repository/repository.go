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
}

type Statistic interface {
	GetStatistics(ctx context.Context, telegramChatID int64) (core.UsersStatistics, error)
	SetStatistics(ctx context.Context, telegramChatID int64, stats core.UsersStatistics) error
}

type Games interface {
	SetActiveGame(ctx context.Context, telegramChatID int64, game core.Game) error
	FinishActiveGame(ctx context.Context, telegramChatID int64) (core.Game, error)
	GetActiveGame(ctx context.Context, telegramChatID int64) (core.Game, error)
}

type Cards interface {
	AddCardToDealer(ctx context.Context, telegramChatID int64, card core.Card) error
	AddCardToPlayerHand(ctx context.Context, telegramChatID int64, username string, card core.Card) error
	AddNewPlayerHand(ctx context.Context, telegramChatID int64, playerHand core.PlayerHand) error
	DrawCard(ctx context.Context, telegramChatID int64) (core.Card, error)
	DrawCards(ctx context.Context, telegramChatID int64, amount int) (core.Cards, error)
	StopDrawing(ctx context.Context, telegramChatID int64, username string) error
	GetPlayerHand(ctx context.Context, telegramChatID int64, username string) (*core.PlayerHand, error)
}

type Repository struct {
	Chats
	Statistic
	Games
	Cards
}

func NewRepository(db *mongo.Client) *Repository {
	chatsCollection := db.Database("homestead").Collection("chats")

	return &Repository{
		Chats:     mongodb.NewChat(chatsCollection),
		Statistic: mongodb.NewStatistic(chatsCollection),
		Games:     mongodb.NewGame(chatsCollection),
		Cards:     mongodb.NewGame(chatsCollection),
	}
}
