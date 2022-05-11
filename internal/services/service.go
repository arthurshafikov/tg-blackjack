package services

import (
	"context"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
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
	NewGame(ctx context.Context, telegramChatID int64) (*core.Game, error)
	FinishGame(ctx context.Context, telegramChatID int64) (*core.Game, core.UsersStatistics, error)
	CheckIfGameShouldBeFinished(ctx context.Context, telegramChatID int64) (bool, error)
}

type Cards interface {
	DrawCardFromDeckToPlayer(ctx context.Context, telegramChatID int64, username string) (*core.Player, error)
	DrawCardFromDeckToDealer(ctx context.Context, telegramChatID int64) (core.Card, error)
}

type Players interface {
	StopDrawing(ctx context.Context, telegramChatID int64, player *core.Player) error
	GetPlayer(ctx context.Context, telegramChatID int64, username string) (*core.Player, error)
	AddNewPlayer(ctx context.Context, telegramChatID int64, player core.Player) error
}

type Logger interface {
	Error(err error)
}

type Services struct {
	Chats
	Statistics
	Games
	Cards
	Players
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
		Chats:      chats,
		Statistics: statistics,
		Games:      games,
		Cards:      cards,
		Players:    players,
	}
}
