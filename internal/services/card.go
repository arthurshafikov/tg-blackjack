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
) (*core.PlayerHand, error) {
	playerHand, err := c.repo.GetPlayerHand(ctx, telegramChatID, username)
	if err != nil {
		if !errors.Is(err, core.ErrNotFound) {
			c.logger.Error(err)

			return playerHand, core.ErrServerError
		}

		return c.createNewPlayerHand(ctx, telegramChatID, username)
	}

	if playerHand.Stop {
		return playerHand, core.ErrCantDraw
	}

	card, err := c.drawCardFromDeckToUser(ctx, telegramChatID, username)
	if err != nil {
		c.logger.Error(err)

		return playerHand, core.ErrServerError
	}

	playerHand.Cards = append(playerHand.Cards, card)

	if playerHand.Cards.CountValue() >= 21 {
		playerHand.Stop = true
		if err := c.StopDrawing(ctx, telegramChatID, username); err != nil {
			c.logger.Error(err)

			return nil, core.ErrServerError
		}

		if playerHand.Cards.CountValue() > 21 {
			return playerHand, core.ErrMoreThan21
		}
	}

	return playerHand, nil
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
	username string,
) error {
	if err := c.repo.StopDrawing(ctx, telegramChatID, username); err != nil {
		if !errors.Is(err, core.ErrNotFound) {
			c.logger.Error(err)

			return core.ErrServerError
		}

		return core.ErrNotFound
	}

	return nil
}

func (c *CardService) createNewPlayerHand(
	ctx context.Context,
	telegramChatID int64,
	username string,
) (*core.PlayerHand, error) {
	playerCards, err := c.repo.DrawCards(ctx, telegramChatID, 2)
	if err != nil {
		c.logger.Error(err)

		return nil, core.ErrServerError
	}

	playerHand := &core.PlayerHand{
		Username: username,
		Cards:    playerCards,
	}
	if err := c.repo.AddNewPlayerHand(ctx, telegramChatID, *playerHand); err != nil {
		c.logger.Error(err)

		return nil, core.ErrServerError
	}

	return playerHand, nil
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
	if err := c.repo.AddCardToPlayerHand(ctx, telegramChatID, username, card); err != nil {
		return card, err
	}

	return card, nil
}
