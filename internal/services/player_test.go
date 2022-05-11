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

var player = &core.Player{
	Username: username,
}

func getPlayerServiceDependencies(t *testing.T) (
	context.Context,
	*mock_services.MockLogger,
	*mock_repository.MockPlayers,
) {
	t.Helper()
	ctrl := gomock.NewController(t)

	return context.Background(),
		mock_services.NewMockLogger(ctrl),
		mock_repository.NewMockPlayers(ctrl)
}

func TestStopDrawing(t *testing.T) {
	ctx, logger, repo := getPlayerServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetPlayer(ctx, telegramChatID, player.Username).Return(player, nil),
		repo.EXPECT().SetPlayerStopAndBusted(ctx, telegramChatID, player).Return(nil),
	)
	service := NewPlayerService(logger, repo)

	err := service.StopDrawing(ctx, telegramChatID, player)
	require.NoError(t, err)
}

func TestStopDrawingPlayerAlreadyStopped(t *testing.T) {
	ctx, logger, repo := getPlayerServiceDependencies(t)
	player := *player
	player.Stop = true
	gomock.InOrder(
		repo.EXPECT().GetPlayer(ctx, telegramChatID, player.Username).Return(&player, nil),
	)
	service := NewPlayerService(logger, repo)

	err := service.StopDrawing(ctx, telegramChatID, &player)
	require.ErrorIs(t, err, core.ErrAlreadyStopped)
}

func TestStopDrawingGetPlayerServerError(t *testing.T) {
	ctx, logger, repo := getPlayerServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetPlayer(ctx, telegramChatID, player.Username).Return(player, core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewPlayerService(logger, repo)

	err := service.StopDrawing(ctx, telegramChatID, player)
	require.ErrorIs(t, err, core.ErrServerError)
}

func TestStopDrawingServerError(t *testing.T) {
	ctx, logger, repo := getPlayerServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetPlayer(ctx, telegramChatID, player.Username).Return(player, nil),
		repo.EXPECT().SetPlayerStopAndBusted(ctx, telegramChatID, player).Return(core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewPlayerService(logger, repo)

	err := service.StopDrawing(ctx, telegramChatID, player)
	require.ErrorIs(t, err, core.ErrServerError)
}

func TestStopDrawingNotFound(t *testing.T) {
	ctx, logger, repo := getPlayerServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetPlayer(ctx, telegramChatID, player.Username).Return(player, nil),
		repo.EXPECT().SetPlayerStopAndBusted(ctx, telegramChatID, player).Return(core.ErrNotFound),
	)
	service := NewPlayerService(logger, repo)

	err := service.StopDrawing(ctx, telegramChatID, player)
	require.ErrorIs(t, err, core.ErrNotFound)
}

func TestGetPlayer(t *testing.T) {
	ctx, logger, repo := getPlayerServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetPlayer(ctx, telegramChatID, player.Username).Return(player, nil),
	)
	service := NewPlayerService(logger, repo)

	result, err := service.GetPlayer(ctx, telegramChatID, player.Username)
	require.NoError(t, err)
	require.Equal(t, player, result)
}

func TestGetPlayerServerError(t *testing.T) {
	ctx, logger, repo := getPlayerServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetPlayer(ctx, telegramChatID, player.Username).Return(nil, core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewPlayerService(logger, repo)

	result, err := service.GetPlayer(ctx, telegramChatID, player.Username)
	require.ErrorIs(t, err, core.ErrServerError)
	require.Nil(t, result)
}

func TestGetPlayerNotFound(t *testing.T) {
	ctx, logger, repo := getPlayerServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetPlayer(ctx, telegramChatID, player.Username).Return(nil, core.ErrNotFound),
	)
	service := NewPlayerService(logger, repo)

	result, err := service.GetPlayer(ctx, telegramChatID, player.Username)
	require.ErrorIs(t, err, core.ErrNotFound)
	require.Nil(t, result)
}

func TestAddPlayer(t *testing.T) {
	ctx, logger, repo := getPlayerServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().AddNewPlayer(ctx, telegramChatID, *player).Return(nil),
	)
	service := NewPlayerService(logger, repo)

	err := service.AddNewPlayer(ctx, telegramChatID, *player)
	require.NoError(t, err)
}

func TestAddPlayerServerError(t *testing.T) {
	ctx, logger, repo := getPlayerServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().AddNewPlayer(ctx, telegramChatID, *player).Return(core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewPlayerService(logger, repo)

	err := service.AddNewPlayer(ctx, telegramChatID, *player)
	require.ErrorIs(t, err, core.ErrServerError)
}
