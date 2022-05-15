package services

import (
	"context"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/repository"
)

type Cards interface {
	DrawCardFromDeckToDealer(ctx context.Context, telegramChatID int64) (core.Card, error)
	DrawCardFromDeckToPlayer(ctx context.Context, telegramChatID int64, username string) (*core.Player, error)
}

type Chats interface {
	CheckChatExists(ctx context.Context, telegramChatID int64) error
	RegisterChat(ctx context.Context, telegramChatID int64) error
}

type Games interface {
	CheckIfGameShouldBeFinished(ctx context.Context, telegramChatID int64) (bool, error)
	FinishGame(ctx context.Context, telegramChatID int64) (*core.Game, core.UsersStatistics, error)
	NewGame(ctx context.Context, telegramChatID int64) (*core.Game, error)
}

type Players interface {
	AddNewPlayer(ctx context.Context, telegramChatID int64, player core.Player) error
	GetPlayer(ctx context.Context, telegramChatID int64, username string) (*core.Player, error)
	StopDrawing(ctx context.Context, telegramChatID int64, player *core.Player) error
}

type Statistics interface {
	GetStatistics(ctx context.Context, telegramChatID int64) (core.UsersStatistics, error)
	IncrementStatistic(ctx context.Context, telegramChatID int64, gameResult core.UsersStatistics) error
}

type Logger interface {
	Error(err error)
}

type Services struct {
	Cards
	Chats
	Games
	Players
	Statistics
}

type Deps struct {
	Config     *config.Config
	Repository *repository.Repository
	Logger
}

func NewServices(deps Deps) *Services {
	chats := NewChatService(deps.Logger, deps.Repository.Chats)
	statistics := NewStatisticService(deps.Logger, deps.Repository.Statistic)
	players := NewPlayerService(deps.Logger, deps.Repository.Players)
	cards := NewCardService(deps.Logger, deps.Repository.Cards, players, deps.Config.App.NumOfDecks)
	games := NewGameService(deps.Logger, deps.Repository.Games, statistics, cards)

	return &Services{
		Cards:      cards,
		Chats:      chats,
		Games:      games,
		Players:    players,
		Statistics: statistics,
	}
}
