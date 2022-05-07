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

func getChatServiceDependencies(t *testing.T) (
	context.Context,
	*mock_services.MockLogger,
	*mock_repository.MockChats,
) {
	t.Helper()
	ctrl := gomock.NewController(t)

	return context.Background(),
		mock_services.NewMockLogger(ctrl),
		mock_repository.NewMockChats(ctrl)
}

func TestCheckChatExists(t *testing.T) {
	ctx, logger, repo := getChatServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().CheckChatExists(ctx, telegramChatID).Return(nil),
	)
	service := NewChatService(logger, repo)

	err := service.CheckChatExists(ctx, telegramChatID)
	require.NoError(t, err)
}

func TestCheckChatExistsNotFound(t *testing.T) {
	ctx, logger, repo := getChatServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().CheckChatExists(ctx, telegramChatID).Return(core.ErrNotFound),
	)
	service := NewChatService(logger, repo)

	err := service.CheckChatExists(ctx, telegramChatID)
	require.ErrorIs(t, err, core.ErrNotFound)
}

func TestCheckChatExistsServerError(t *testing.T) {
	ctx, logger, repo := getChatServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().CheckChatExists(ctx, telegramChatID).Return(core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewChatService(logger, repo)

	err := service.CheckChatExists(ctx, telegramChatID)
	require.ErrorIs(t, err, core.ErrServerError)
}

func TestRegisterChat(t *testing.T) {
	ctx, logger, repo := getChatServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().CheckChatExists(ctx, telegramChatID).Return(core.ErrNotFound),
		repo.EXPECT().RegisterChat(ctx, telegramChatID).Return(nil),
	)
	service := NewChatService(logger, repo)

	err := service.RegisterChat(ctx, telegramChatID)
	require.NoError(t, err)
}

func TestRegisterChatServerError(t *testing.T) {
	ctx, logger, repo := getChatServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().CheckChatExists(ctx, telegramChatID).Return(core.ErrNotFound),
		repo.EXPECT().RegisterChat(ctx, telegramChatID).Return(core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewChatService(logger, repo)

	err := service.RegisterChat(ctx, telegramChatID)
	require.ErrorIs(t, err, core.ErrServerError)
}

func TestRegisterChatExists(t *testing.T) {
	ctx, logger, repo := getChatServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().CheckChatExists(ctx, telegramChatID).Return(nil),
	)
	service := NewChatService(logger, repo)

	err := service.RegisterChat(ctx, telegramChatID)
	require.ErrorIs(t, err, core.ErrChatRegistered)
}
