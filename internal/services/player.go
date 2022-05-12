package services

import (
	"context"
	"errors"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/repository"
)

type PlayerService struct {
	logger Logger
	repo   repository.Players
}

func NewPlayerService(logger Logger, repo repository.Players) *PlayerService {
	return &PlayerService{
		logger: logger,
		repo:   repo,
	}
}

func (p *PlayerService) StopDrawing(
	ctx context.Context,
	telegramChatID int64,
	player *core.Player,
) error {
	playerStopped, err := p.checkIfPlayerIsStopped(ctx, telegramChatID, player.Username)
	if err != nil {
		if !errors.Is(err, core.ErrNotFound) && !errors.Is(err, core.ErrNoActiveGame) {
			p.logger.Error(err)

			return core.ErrServerError
		}

		return err
	}
	if playerStopped {
		return core.ErrAlreadyStopped
	}

	player.Stop = true
	if err := p.repo.SetPlayerStopAndBusted(ctx, telegramChatID, player); err != nil {
		if !errors.Is(err, core.ErrNotFound) {
			p.logger.Error(err)

			return core.ErrServerError
		}

		return core.ErrNotFound
	}

	return nil
}

func (p *PlayerService) GetPlayer(ctx context.Context, telegramChatID int64, username string) (*core.Player, error) {
	player, err := p.repo.GetPlayer(ctx, telegramChatID, username)
	if err != nil {
		if !errors.Is(err, core.ErrNotFound) && !errors.Is(err, core.ErrNoActiveGame) {
			p.logger.Error(err)

			return nil, core.ErrServerError
		}

		return nil, err
	}

	return player, nil
}

func (p *PlayerService) AddNewPlayer(ctx context.Context, telegramChatID int64, player core.Player) error {
	if err := p.repo.AddNewPlayer(ctx, telegramChatID, player); err != nil {
		p.logger.Error(err)

		return core.ErrServerError
	}

	return nil
}

func (p *PlayerService) checkIfPlayerIsStopped(
	ctx context.Context,
	telegramChatID int64,
	username string,
) (bool, error) {
	player, err := p.repo.GetPlayer(ctx, telegramChatID, username)
	if err != nil {
		return false, err
	}

	if player.Busted {
		return false, core.ErrBusted
	}

	return player.Stop, nil
}
