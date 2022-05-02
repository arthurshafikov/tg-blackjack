package services

import (
	"context"
	"errors"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/repository"
)

type CardService struct {
	logger Logger
	repo   repository.Cards
}

func NewCardService(logger Logger, repo repository.Cards) *CardService {
	return &CardService{
		logger: logger,
		repo:   repo,
	}
}

func (c *CardService) DrawCard(
	ctx context.Context,
	telegramChatID int64,
	username string,
) (*core.Player, error) {
	player, err := c.repo.GetPlayer(ctx, telegramChatID, username)
	if err != nil {
		if !errors.Is(err, core.ErrNotFound) {
			c.logger.Error(err)

			return player, core.ErrServerError
		}

		return c.createNewPlayer(ctx, telegramChatID, username)
	}

	if player.Stop {
		return player, core.ErrCantDraw
	}

	card, err := c.drawCardFromDeckToUser(ctx, telegramChatID, username)
	if err != nil {
		c.logger.Error(err)

		return player, core.ErrServerError
	}

	player.Cards = append(player.Cards, card)

	if player.Cards.CountValue() >= 21 {
		player.Stop = true
		player.Busted = true

		if player.Cards.CountValue() == 21 {
			player.Busted = false
		}

		if err := c.StopDrawing(ctx, telegramChatID, player); err != nil {
			c.logger.Error(err)

			return nil, core.ErrServerError
		}

		if player.Busted {
			return player, core.ErrBusted
		}
	}

	return player, nil
}

func (c *CardService) DrawCardFromDeckToDealer(ctx context.Context, telegramChatID int64) (core.Card, error) {
	card, err := c.repo.DrawCard(ctx, telegramChatID)
	if err != nil {
		c.logger.Error(err)

		return card, core.ErrServerError
	}
	if err := c.repo.AddCardToDealer(ctx, telegramChatID, card); err != nil {
		c.logger.Error(err)

		return card, core.ErrServerError
	}

	return card, nil
}

func (c *CardService) StopDrawing(
	ctx context.Context,
	telegramChatID int64,
	player *core.Player,
) error {
	if err := c.repo.StopDrawing(ctx, telegramChatID, player); err != nil {
		if !errors.Is(err, core.ErrNotFound) {
			c.logger.Error(err)

			return core.ErrServerError
		}

		return core.ErrNotFound
	}

	return nil
}

func (c *CardService) createNewPlayer(
	ctx context.Context,
	telegramChatID int64,
	username string,
) (*core.Player, error) {
	playerCards, err := c.repo.DrawCards(ctx, telegramChatID, 2)
	if err != nil {
		c.logger.Error(err)

		return nil, core.ErrServerError
	}

	player := &core.Player{
		Username: username,
		Cards:    playerCards,
	}
	if err := c.repo.AddNewPlayer(ctx, telegramChatID, *player); err != nil {
		c.logger.Error(err)

		return nil, core.ErrServerError
	}

	return player, nil
}

func (c *CardService) drawCardFromDeckToUser(
	ctx context.Context,
	telegramChatID int64,
	username string,
) (core.Card, error) {
	card, err := c.repo.DrawCard(ctx, telegramChatID) // empty deck?
	if err != nil {
		return card, err
	}
	if err := c.repo.AddCardToPlayer(ctx, telegramChatID, username, card); err != nil {
		return card, err
	}

	return card, nil
}
