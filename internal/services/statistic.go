package services

import (
	"context"
	"errors"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/repository"
)

type StatisticService struct {
	logger Logger
	repo   repository.Statistic
}

func NewStatisticService(logger Logger, repo repository.Statistic) *StatisticService {
	return &StatisticService{
		logger: logger,
		repo:   repo,
	}
}

func (s *StatisticService) GetStatistics(
	ctx context.Context,
	telegramChatID int64,
) (core.UsersStatistics, error) {
	statistics, err := s.repo.GetStatistics(ctx, telegramChatID)
	if err != nil {
		if !errors.Is(err, core.ErrNotFound) {
			s.logger.Error(err)

			return statistics, core.ErrServerError
		}

		return statistics, core.ErrNotFound
	}

	return statistics, nil
}

func (s *StatisticService) IncrementStatistic(
	ctx context.Context,
	telegramChatID int64,
	gameResult core.UsersStatistics,
) error {
	stats, err := s.GetStatistics(ctx, telegramChatID)
	if err != nil {
		return err
	}

	for username, result := range gameResult {
		stats[username] += result
	}

	return s.repo.SetStatistics(ctx, telegramChatID, stats)
}
