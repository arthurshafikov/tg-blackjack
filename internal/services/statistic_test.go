package services

import (
	"context"
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	mock_repository "github.com/arthurshafikov/tg-blackjack/internal/repository/mocks"
	mock_services "github.com/arthurshafikov/tg-blackjack/internal/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	statistics = core.UsersStatistics{
		"some_user":  10,
		"other_user": -5,
	}
	gameResult = core.UsersStatistics{
		"some_user":  -1,
		"other_user": 2,
	}
	setStatistics = core.UsersStatistics{
		"some_user":  9,
		"other_user": -3,
	}
)

func getStatisticServiceDependencies(t *testing.T) (
	context.Context,
	*mock_services.MockLogger,
	*mock_repository.MockStatistic,
) {
	t.Helper()
	ctrl := gomock.NewController(t)

	return context.Background(),
		mock_services.NewMockLogger(ctrl),
		mock_repository.NewMockStatistic(ctrl)
}

func TestGetStatistics(t *testing.T) {
	ctx, logger, repo := getStatisticServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetStatistics(ctx, telegramChatID).Return(statistics, nil),
	)
	service := NewStatisticService(logger, repo)

	result, err := service.GetStatistics(ctx, telegramChatID)
	require.NoError(t, err)
	require.Equal(t, statistics, result)
}

func TestGetStatisticsServerError(t *testing.T) {
	ctx, logger, repo := getStatisticServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetStatistics(ctx, telegramChatID).Return(statistics, core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewStatisticService(logger, repo)

	_, err := service.GetStatistics(ctx, telegramChatID)
	require.ErrorIs(t, err, core.ErrServerError)
}

func TestGetStatisticsNotFound(t *testing.T) {
	ctx, logger, repo := getStatisticServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetStatistics(ctx, telegramChatID).Return(statistics, core.ErrNotFound),
	)
	service := NewStatisticService(logger, repo)

	_, err := service.GetStatistics(ctx, telegramChatID)
	require.ErrorIs(t, err, core.ErrNotFound)
}

func TestIncrementStatistic(t *testing.T) {
	ctx, logger, repo := getStatisticServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetStatistics(ctx, telegramChatID).Return(statistics, nil),
		repo.EXPECT().SetStatistics(ctx, telegramChatID, setStatistics).Return(nil),
	)
	service := NewStatisticService(logger, repo)

	err := service.IncrementStatistic(ctx, telegramChatID, gameResult)
	require.NoError(t, err)
}

func TestIncrementStatisticServerError(t *testing.T) {
	ctx, logger, repo := getStatisticServiceDependencies(t)
	statistics := core.UsersStatistics{
		"some_user":  10,
		"other_user": -5,
	}
	gomock.InOrder(
		repo.EXPECT().GetStatistics(ctx, telegramChatID).Return(statistics, nil),
		repo.EXPECT().SetStatistics(ctx, telegramChatID, setStatistics).Return(core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewStatisticService(logger, repo)

	err := service.IncrementStatistic(ctx, telegramChatID, gameResult)
	require.ErrorIs(t, err, core.ErrServerError)
}

func TestIncrementStatisticNotFound(t *testing.T) {
	ctx, logger, repo := getStatisticServiceDependencies(t)
	statistics := core.UsersStatistics{
		"some_user":  10,
		"other_user": -5,
	}
	gomock.InOrder(
		repo.EXPECT().GetStatistics(ctx, telegramChatID).Return(statistics, nil),
		repo.EXPECT().SetStatistics(ctx, telegramChatID, setStatistics).Return(core.ErrNotFound),
	)
	service := NewStatisticService(logger, repo)

	err := service.IncrementStatistic(ctx, telegramChatID, gameResult)
	require.ErrorIs(t, err, core.ErrNotFound)
}
