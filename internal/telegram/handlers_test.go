package telegram

import (
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	mock_telegram "github.com/arthurshafikov/tg-blackjack/internal/telegram/mocks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func getBotWithHandler(t *testing.T) (*Bot, *mock_telegram.MockCommandHandler) {
	t.Helper()

	ctrl := gomock.NewController(t)
	commandHandlerMock := mock_telegram.NewMockCommandHandler(ctrl)
	bot := &Bot{
		commandHandler: commandHandlerMock,
	}

	return bot, commandHandlerMock
}

func TestHandleCommand(t *testing.T) {
	bot, commandHandlerMock := getBotWithHandler(t)
	command := &tgbotapi.Message{
		Text: "/Some@asds",
	}

	t.Run("HandleStart", func(t *testing.T) {
		command.Text = "/" + core.StartCommand
		command.Entities = &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(command.Text)}}
		gomock.InOrder(
			commandHandlerMock.EXPECT().HandleStart(command).Return(nil),
		)

		require.NoError(t, bot.handleCommand(command))
	})

	t.Run("HandleStats", func(t *testing.T) {
		command.Text = "/" + core.StatsCommand
		command.Entities = &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(command.Text)}}
		gomock.InOrder(
			commandHandlerMock.EXPECT().HandleStats(command).Return(nil),
		)

		require.NoError(t, bot.handleCommand(command))
	})

	t.Run("HandleNewGame", func(t *testing.T) {
		command.Text = "/" + core.NewGame
		command.Entities = &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(command.Text)}}
		gomock.InOrder(
			commandHandlerMock.EXPECT().HandleNewGame(command).Return(nil),
		)

		require.NoError(t, bot.handleCommand(command))
	})

	t.Run("HandleDrawCard", func(t *testing.T) {
		command.Text = "/" + core.DrawCard
		command.Entities = &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(command.Text)}}
		gomock.InOrder(
			commandHandlerMock.EXPECT().HandleDrawCard(command).Return(nil),
		)

		require.NoError(t, bot.handleCommand(command))
	})

	t.Run("HandleStopDrawing", func(t *testing.T) {
		command.Text = "/" + core.StopDrawing
		command.Entities = &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(command.Text)}}
		gomock.InOrder(
			commandHandlerMock.EXPECT().HandleStopDrawing(command).Return(nil),
		)

		require.NoError(t, bot.handleCommand(command))
	})
}
