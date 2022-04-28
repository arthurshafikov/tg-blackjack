package services

import (
	"context"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/repository"
)

type Chats interface {
	CheckChatExists(ctx context.Context, telegramChatID int64) error
	RegisterChat(ctx context.Context, telegramChatID int64) error
	GetStatistics(ctx context.Context, telegramChatID int64) (core.UsersStatistics, error)
}

type Logger interface {
	Error(err error)
}

type Services struct {
	Chats
}

type Deps struct {
	Repository *repository.Repository
	Logger
}

func NewServices(deps Deps) *Services {
	return &Services{
		Chats: NewChatService(deps.Logger, deps.Repository.Chats),
	}
}
