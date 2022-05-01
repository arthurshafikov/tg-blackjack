package services

import (
	"context"
	"errors"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/repository"
)

type GameService struct {
	logger Logger

	repo repository.Games

	statisticService Statistics
}

func NewGameService(logger Logger, repo repository.Games, statisticService Statistics) *GameService {
	return &GameService{
		logger:           logger,
		repo:             repo,
		statisticService: statisticService,
	}
}

func (g *GameService) NewGame(ctx context.Context, telegramChatID int64) error {
	deck := core.NewDeck()
	dealerHand, err := deck.DrawCards(2)
	if err != nil {
		g.logger.Error(err)

		return core.ErrServerError
	}

	game := core.Game{
		Deck:         *deck,
		DealerHand:   dealerHand,
		PlayersHands: []core.PlayerHand{},
	}

	if err := g.repo.SetActiveGame(ctx, telegramChatID, game); err != nil {
		if !errors.Is(err, core.ErrActiveGame) {
			g.logger.Error(err)

			return core.ErrServerError
		}

		return core.ErrActiveGame
	}

	return nil
}

func (g *GameService) CheckIfGameShouldBeFinished(ctx context.Context, telegramChatID int64) (bool, error) {
	result := true
	game, err := g.repo.GetActiveGame(ctx, telegramChatID)
	if err != nil {
		if !errors.Is(err, core.ErrNotFound) {
			g.logger.Error(err)

			return result, core.ErrServerError
		}

		return result, core.ErrNotFound
	}

	for _, playerHand := range game.PlayersHands {
		if !playerHand.Stop {
			result = false
		}
	}

	return result, nil
}

func (g *GameService) FinishGame(ctx context.Context, telegramChatID int64) (core.UsersStatistics, error) {
	game, err := g.repo.FinishActiveGame(ctx, telegramChatID)
	if err != nil {
		if !errors.Is(err, core.ErrNotFound) {
			g.logger.Error(err)

			return nil, core.ErrServerError
		}

		return nil, core.ErrNotFound
	}

	// todo dealer should draw cards, implement...
	dealerValue := game.DealerHand.CountValue()

	gameResult := core.UsersStatistics{}

	for _, playerHand := range game.PlayersHands {
		result := 0
		playerValue := playerHand.Cards.CountValue()

		// todo implement blackjack logic
		if playerValue < dealerValue {
			result = -1
		} else if playerValue > dealerValue {
			result = +1
		}

		gameResult[playerHand.Username] = result
	}

	if err := g.statisticService.IncrementStatistic(ctx, telegramChatID, gameResult); err != nil {
		g.logger.Error(err)

		return nil, core.ErrServerError
	}

	return gameResult, nil
}
