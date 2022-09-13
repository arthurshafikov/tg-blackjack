package telegram

import (
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	mock_services "github.com/arthurshafikov/tg-blackjack/internal/services/mocks"
	mock_handlers "github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers/mocks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/mock/gomock"
)

func getTelegramBotWithTelegramHandlerHelper(
	t *testing.T,
) (*Bot, *mock_handlers.MockTelegramHandlerHelper, *mock_services.MockLogger) {
	t.Helper()

	ctrl := gomock.NewController(t)
	helperMock := mock_handlers.NewMockTelegramHandlerHelper(ctrl)
	loggerMock := mock_services.NewMockLogger(ctrl)
	bot := &Bot{
		helper: helperMock,
		logger: loggerMock,
	}

	return bot, helperMock, loggerMock
}

func TestHandleError(t *testing.T) {
	bot, helperMock, _ := getTelegramBotWithTelegramHandlerHelper(t)
	responseConfig := tgbotapi.NewMessage(chatID, core.ErrServerError.Error())
	responseConfig.ParseMode = MarkDownParseMode
	gomock.InOrder(
		helperMock.EXPECT().NewMessage(chatID, core.ErrServerError.Error()).Return(responseConfig),
		helperMock.EXPECT().SendMessage(responseConfig).Return(nil),
	)

	bot.handleError(chatID, core.ErrServerError)
}

func TestHandleErrorSendMessageReturnsError(t *testing.T) {
	bot, helperMock, loggerMock := getTelegramBotWithTelegramHandlerHelper(t)
	responseConfig := tgbotapi.NewMessage(chatID, core.ErrServerError.Error())
	responseConfig.ParseMode = MarkDownParseMode
	gomock.InOrder(
		helperMock.EXPECT().NewMessage(chatID, core.ErrServerError.Error()).Return(responseConfig),
		helperMock.EXPECT().SendMessage(responseConfig).Return(core.ErrServerError),
		loggerMock.EXPECT().Error(core.ErrServerError),
	)

	bot.handleError(chatID, core.ErrServerError)
}
