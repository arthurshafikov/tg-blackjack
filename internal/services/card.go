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

	playerService Players
}

func NewCardService(logger Logger, repo repository.Cards, playerService Players) *CardService {
	return &CardService{
		logger: logger,
		repo:   repo,

		playerService: playerService,
	}
}

func (c *CardService) DrawCard(
	ctx context.Context,
	telegramChatID int64,
	username string,
) (*core.Player, error) {
	player, err := c.playerService.GetPlayer(ctx, telegramChatID, username)
	if err != nil {
		if errors.Is(err, core.ErrNoActiveGame) {
			return player, core.ErrNoActiveGame
		}

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

		if err := c.playerService.StopDrawing(ctx, telegramChatID, player); err != nil {
			c.logger.Error(err)

			return nil, core.ErrServerError
		}

		if player.Busted {
			return player, core.ErrBusted
		}
	}

	return player, nil
}

func (c *CardService) DrawCards(ctx context.Context, telegramChatID int64, amount int) (core.Cards, error) {
	cards := core.Cards{}
	for amount > 0 {
		card, err := c.repo.DrawCard(ctx, telegramChatID)
		if err != nil {
			return nil, err
		}

		cards = append(cards, card)

		amount--
	}

	return cards, nil
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

func (c *CardService) createNewPlayer(
	ctx context.Context,
	telegramChatID int64,
	username string,
) (*core.Player, error) {
	playerCards, err := c.DrawCards(ctx, telegramChatID, 2)
	if err != nil {
		c.logger.Error(err)

		return nil, core.ErrServerError
	}

	player := &core.Player{
		Username: username,
		Cards:    playerCards,
	}

	if playerCards.IsBlackjack() {
		player.Stop = true
	}

	if err := c.playerService.AddNewPlayer(ctx, telegramChatID, *player); err != nil {
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
	card, err := c.repo.DrawCard(ctx, telegramChatID)
	if err != nil {
		return card, err
	}

	return card, c.repo.AddCardToPlayer(ctx, telegramChatID, username, card)
}
