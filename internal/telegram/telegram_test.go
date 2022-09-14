package telegram

import (
	"context"
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	mock_services "github.com/arthurshafikov/tg-blackjack/internal/services/mocks"
	mock_handlers "github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers/mocks"
	mock_telegram "github.com/arthurshafikov/tg-blackjack/internal/telegram/mocks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

var (
	updateConfig = tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 60,
	}
	updateChannel = make(chan tgbotapi.Update)
	msg           = &tgbotapi.Message{
		Text: "some msg",
		Chat: &tgbotapi.Chat{
			ID: chatID,
		},
	}
)

type MockBag struct {
	Cards      *mock_services.MockCards
	Chats      *mock_services.MockChats
	Games      *mock_services.MockGames
	Players    *mock_services.MockPlayers
	Statistics *mock_services.MockStatistics
	Logger     *mock_services.MockLogger

	Helper *mock_handlers.MockTelegramHandlerHelper

	CommandHandler *mock_telegram.MockCommandHandler
}

func getBotWithMockBag(t *testing.T) (*Bot, *MockBag) {
	t.Helper()
	ctrl := gomock.NewController(t)
	cardsMock := mock_services.NewMockCards(ctrl)
	chatsMock := mock_services.NewMockChats(ctrl)
	gamesMock := mock_services.NewMockGames(ctrl)
	playersMock := mock_services.NewMockPlayers(ctrl)
	statisticsMock := mock_services.NewMockStatistics(ctrl)
	helperMock := mock_handlers.NewMockTelegramHandlerHelper(ctrl)
	loggerMock := mock_services.NewMockLogger(ctrl)
	services := &services.Services{
		Cards:      cardsMock,
		Chats:      chatsMock,
		Games:      gamesMock,
		Players:    playersMock,
		Statistics: statisticsMock,
	}
	commandHandlerMock := mock_telegram.NewMockCommandHandler(ctrl)

	return &Bot{
			ctx:      context.Background(),
			helper:   helperMock,
			logger:   loggerMock,
			services: services,

			commandHandler: commandHandlerMock,
		}, &MockBag{
			Cards:      cardsMock,
			Chats:      chatsMock,
			Games:      gamesMock,
			Players:    playersMock,
			Statistics: statisticsMock,
			Logger:     loggerMock,

			Helper: helperMock,

			CommandHandler: commandHandlerMock,
		}
}

func TestStart(t *testing.T) {
	g := errgroup.Group{}
	bot, mockBag := getBotWithMockBag(t)

	gomock.InOrder(
		mockBag.Helper.EXPECT().NewUpdateChannel(0).Return(updateConfig),
		mockBag.Helper.EXPECT().GetUpdatesChan(updateConfig).Return(updateChannel, nil),
	)
	g.Go(func() error {
		defer close(updateChannel)

		t.Run("start command", func(t *testing.T) {
			command := *msg
			command.Text = "/" + core.StartCommand
			command.Entities = &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(command.Text)}}
			gomock.InOrder(
				mockBag.Chats.EXPECT().CheckChatExists(bot.ctx, chatID).Return(core.ErrServerError),
				mockBag.CommandHandler.EXPECT().HandleStart(&command).Return(nil),
			)
			updateChannel <- tgbotapi.Update{Message: &command}
		})

		t.Run("stats command", func(t *testing.T) {
			command := *msg
			command.Text = "/" + core.StatsCommand
			command.Entities = &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(command.Text)}}
			gomock.InOrder(
				mockBag.Chats.EXPECT().CheckChatExists(bot.ctx, chatID).Return(nil),
				mockBag.CommandHandler.EXPECT().HandleStats(&command).Return(nil),
			)
			updateChannel <- tgbotapi.Update{Message: &command}
		})

		t.Run("new game command", func(t *testing.T) {
			command := *msg
			command.Text = "/" + core.NewGameCommand
			command.Entities = &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(command.Text)}}
			gomock.InOrder(
				mockBag.Chats.EXPECT().CheckChatExists(bot.ctx, chatID).Return(nil),
				mockBag.CommandHandler.EXPECT().HandleNewGame(&command).Return(nil),
			)
			updateChannel <- tgbotapi.Update{Message: &command}
		})

		t.Run("draw card command", func(t *testing.T) {
			command := *msg
			command.Text = "/" + core.DrawCardCommand
			command.Entities = &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(command.Text)}}
			gomock.InOrder(
				mockBag.Chats.EXPECT().CheckChatExists(bot.ctx, chatID).Return(nil),
				mockBag.CommandHandler.EXPECT().HandleDrawCard(&command).Return(nil),
			)
			updateChannel <- tgbotapi.Update{Message: &command}
		})

		t.Run("stop drawing command", func(t *testing.T) {
			command := *msg
			command.Text = "/" + core.StopDrawingCommand
			command.Entities = &[]tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(command.Text)}}
			gomock.InOrder(
				mockBag.Chats.EXPECT().CheckChatExists(bot.ctx, chatID).Return(nil),
				mockBag.CommandHandler.EXPECT().HandleStopDrawing(&command).Return(nil),
			)
			updateChannel <- tgbotapi.Update{Message: &command}
		})

		return nil
	})

	err := bot.Start()
	require.NoError(t, err)

	require.NoError(t, g.Wait())
}
