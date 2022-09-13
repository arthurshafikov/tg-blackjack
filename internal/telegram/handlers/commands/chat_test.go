package commands

import (
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	username          = "someUsername"
	chatID      int64 = 28
	receivedMsg       = &tgbotapi.Message{
		Text: "Some message",
		Chat: &tgbotapi.Chat{
			ID: chatID,
		},
		From: &tgbotapi.User{
			UserName: username,
		},
	}
	stats = core.UsersStatistics{
		"someUser1": 22,
		"someUser2": -4,
	}
)

func TestHandleStart(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	responseConfig := tgbotapi.NewMessage(chatID, c.messages.ChatCreatedSuccessfully)
	gomock.InOrder(
		mockBag.Chats.EXPECT().RegisterChat(c.ctx, chatID).Return(nil),
		mockBag.Helper.EXPECT().NewMessage(chatID, c.messages.ChatCreatedSuccessfully).Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(nil),
	)

	err := c.HandleStart(receivedMsg)
	require.NoError(t, err)
}

func TestHandleStartChatRegistered(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Chats.EXPECT().RegisterChat(c.ctx, chatID).Return(core.ErrChatRegistered),
	)

	err := c.HandleStart(receivedMsg)
	require.ErrorContains(t, err, c.messages.ChatAlreadyRegistered)
}

func TestHandleStartSendMessageReturnsError(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	responseConfig := tgbotapi.NewMessage(chatID, c.messages.ChatCreatedSuccessfully)
	gomock.InOrder(
		mockBag.Chats.EXPECT().RegisterChat(c.ctx, chatID).Return(nil),
		mockBag.Helper.EXPECT().NewMessage(chatID, c.messages.ChatCreatedSuccessfully).Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(core.ErrServerError),
		mockBag.Logger.EXPECT().Error(gomock.Any()),
	)

	err := c.HandleStart(receivedMsg)
	require.ErrorIs(t, err, core.ErrServerError)
}

func TestHandleStats(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	responseConfig := tgbotapi.NewMessage(chatID, c.messages.ChatCreatedSuccessfully)
	gomock.InOrder(
		mockBag.Statistics.EXPECT().GetStatistics(c.ctx, chatID).Return(stats, nil),
		mockBag.Helper.EXPECT().NewMessage(
			chatID,
			"TopPlayers\n\n*1.* @someUser1 — 22 points\n*2.* @someUser2 — -4 points\n",
		).Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(nil),
	)

	err := c.HandleStats(receivedMsg)
	require.NoError(t, err)
}

func TestHandleStatsGetStatisticsReturnsError(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Statistics.EXPECT().GetStatistics(c.ctx, chatID).Return(nil, core.ErrServerError),
	)

	err := c.HandleStats(receivedMsg)
	require.ErrorIs(t, err, core.ErrServerError)
}

func TestHandleStatsSendMessageReturnsError(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	responseConfig := tgbotapi.NewMessage(chatID, c.messages.ChatCreatedSuccessfully)
	gomock.InOrder(
		mockBag.Statistics.EXPECT().GetStatistics(c.ctx, chatID).Return(stats, nil),
		mockBag.Helper.EXPECT().NewMessage(
			chatID,
			"TopPlayers\n\n*1.* @someUser1 — 22 points\n*2.* @someUser2 — -4 points\n",
		).Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(core.ErrServerError),
		mockBag.Logger.EXPECT().Error(gomock.Any()),
	)

	err := c.HandleStats(receivedMsg)
	require.ErrorIs(t, err, core.ErrServerError)
}
