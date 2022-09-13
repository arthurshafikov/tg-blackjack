package telegram

import (
	"context"
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/core"
	mock_repository "github.com/arthurshafikov/tg-blackjack/internal/repository/mocks"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var chatID int64 = 54

func getUserRepo(t *testing.T) (context.Context, *mock_repository.MockChats) {
	t.Helper()
	ctrl := gomock.NewController(t)

	return context.Background(), mock_repository.NewMockChats(ctrl)
}

func TestCheckAuthorization(t *testing.T) {
	ctx, chatRepo := getUserRepo(t)
	services := &services.Services{
		Chats: chatRepo,
	}
	bot := &Bot{
		ctx:      ctx,
		services: services,
	}
	gomock.InOrder(
		chatRepo.EXPECT().CheckChatExists(ctx, chatID).Return(nil),
	)

	err := bot.checkAuthorization(chatID)
	require.NoError(t, err)
}

func TestCheckAuthorizationReturnsNotFound(t *testing.T) {
	ctx, chatRepo := getUserRepo(t)
	services := &services.Services{
		Chats: chatRepo,
	}
	messages := config.Messages{
		ChatNotExists: "ChatNotExists",
	}
	bot := &Bot{
		ctx:      ctx,
		services: services,
		messages: messages,
	}
	gomock.InOrder(
		chatRepo.EXPECT().CheckChatExists(ctx, chatID).Return(core.ErrNotFound),
	)

	err := bot.checkAuthorization(chatID)
	require.ErrorContains(t, err, "ChatNotExists")
}

func TestCheckAuthorizationServerError(t *testing.T) {
	ctx, chatRepo := getUserRepo(t)
	services := &services.Services{
		Chats: chatRepo,
	}
	bot := &Bot{
		ctx:      ctx,
		services: services,
	}
	gomock.InOrder(
		chatRepo.EXPECT().CheckChatExists(ctx, chatID).Return(core.ErrServerError),
	)

	err := bot.checkAuthorization(chatID)
	require.ErrorIs(t, err, core.ErrServerError)
}
