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
	cardService      Cards
}

func NewGameService(
	logger Logger,
	repo repository.Games,
	statisticService Statistics,
	cardService Cards,
) *GameService {
	return &GameService{
		logger:           logger,
		repo:             repo,
		statisticService: statisticService,
		cardService:      cardService,
	}
}

func (g *GameService) NewGame(ctx context.Context, telegramChatID int64) (*core.Game, error) {
	game := core.Game{
		Dealer:  core.Cards{},
		Players: []core.Player{},
	}

	if err := g.repo.SetActiveGame(ctx, telegramChatID, game); err != nil {
		if !errors.Is(err, core.ErrActiveGame) {
			g.logger.Error(err)

			return nil, core.ErrServerError
		}

		return nil, core.ErrActiveGame
	}

	for i := 0; i < 2; i++ {
		card, err := g.cardService.DrawCardFromDeckToDealer(ctx, telegramChatID)
		if err != nil {
			g.logger.Error(err)

			return nil, core.ErrServerError
		}

		game.Dealer = append(game.Dealer, card)
	}

	return &game, nil
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

	for _, player := range game.Players {
		if !player.Stop {
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
	dealerValue := game.Dealer.CountValue()

	gameResult := core.UsersStatistics{}

	for _, player := range game.Players {
		result := 0
		playerValue := player.Cards.CountValue()

		// todo implement blackjack logic
		if playerValue < dealerValue {
			result = -1
		} else if playerValue > dealerValue {
			result = +1
		}

		gameResult[player.Username] = result
	}

	if err := g.statisticService.IncrementStatistic(ctx, telegramChatID, gameResult); err != nil {
		g.logger.Error(err)

		return nil, core.ErrServerError
	}

	return gameResult, nil
}
