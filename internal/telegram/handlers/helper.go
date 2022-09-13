package handlers

import (
	"context"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramHandlerHelper interface {
	NewMessage(chatID int64, text string) tgbotapi.MessageConfig
	SendMessage(msg tgbotapi.MessageConfig) error
	NewUpdateChannel(offset int) tgbotapi.UpdateConfig
	GetUpdatesChan(config tgbotapi.UpdateConfig) (tgbotapi.UpdatesChannel, error)
}

type HandlerParams struct {
	Ctx      context.Context
	Services *services.Services
	Logger   services.Logger

	Messages config.Messages

	Helper TelegramHandlerHelper
}

type Helper struct {
	bot *tgbotapi.BotAPI
}

func NewHelper(bot *tgbotapi.BotAPI) *Helper {
	return &Helper{
		bot: bot,
	}
}

func (c *Helper) NewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, text)
}

func (c *Helper) SendMessage(msg tgbotapi.MessageConfig) error {
	msg.ParseMode = "markdown"
	_, err := c.bot.Send(msg)

	return err
}

func (c *Helper) NewUpdateChannel(offset int) tgbotapi.UpdateConfig {
	return tgbotapi.NewUpdate(offset)
}

func (c *Helper) GetUpdatesChan(config tgbotapi.UpdateConfig) (tgbotapi.UpdatesChannel, error) {
	return c.bot.GetUpdatesChan(config)
}
