package services

import (
	"context"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/repository"
)

type Chats interface {
	CheckChatExists(ctx context.Context, telegramChatID int64) error
	RegisterChat(ctx context.Context, telegramChatID int64) error
}

type Statistics interface {
	GetStatistics(ctx context.Context, telegramChatID int64) (core.UsersStatistics, error)
	IncrementStatistic(ctx context.Context, telegramChatID int64, gameResult core.UsersStatistics) error
}

type Games interface {
	NewGame(ctx context.Context, telegramChatID int64) error
	FinishGame(ctx context.Context, telegramChatID int64) (core.UsersStatistics, error)
}

type Cards interface {
	DrawCard(ctx context.Context, telegramChatID int64, username string) (*core.PlayerHand, error)
	StopDrawing(ctx context.Context, telegramChatID int64, username string) error
}

type Logger interface {
	Error(err error)
}

type Services struct {
	Chats
	Statistics
	Games
	Cards
}

type Deps struct {
	Repository *repository.Repository
	Logger
}

func NewServices(deps Deps) *Services {
	chats := NewChatService(deps.Logger, deps.Repository.Chats)
	statistics := NewStatisticService(deps.Logger, deps.Repository.Statistic)
	games := NewGameService(deps.Logger, deps.Repository.Games, statistics)
	cards := NewCardService(deps.Logger, deps.Repository.Cards)

	return &Services{
		Chats:      chats,
		Statistics: statistics,
		Games:      games,
		Cards:      cards,
	}
}
