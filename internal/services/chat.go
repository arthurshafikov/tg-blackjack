package services

import (
	"context"
	"errors"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/repository"
)

type ChatService struct {
	logger Logger
	repo   repository.Chats
}

func NewChatService(logger Logger, repo repository.Chats) *ChatService {
	return &ChatService{
		logger: logger,
		repo:   repo,
	}
}

func (c *ChatService) CheckChatExists(ctx context.Context, telegramChatID int64) error {
	if err := c.repo.CheckChatExists(ctx, telegramChatID); err != nil {
		if !errors.Is(err, core.ErrNotFound) {
			c.logger.Error(err)

			return core.ErrServerError
		}

		return core.ErrNotFound
	}

	return nil
}

func (c *ChatService) RegisterChat(ctx context.Context, telegramChatID int64) error {
	if err := c.CheckChatExists(ctx, telegramChatID); err == nil {
		return core.ErrChatRegistered
	}

	if err := c.repo.RegisterChat(ctx, telegramChatID); err != nil {
		c.logger.Error(err)

		return core.ErrNotFound
	}

	return nil
}
