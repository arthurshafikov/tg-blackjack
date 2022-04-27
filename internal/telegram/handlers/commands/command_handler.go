package commands

import (
	"context"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CommandHandler struct {
	ctx      context.Context
	bot      *tgbotapi.BotAPI
	services *services.Services

	messages config.Messages
}

func NewCommandHandler(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	services *services.Services,
	messages config.Messages,
) *CommandHandler {
	return &CommandHandler{
		ctx:      ctx,
		bot:      bot,
		services: services,
		messages: messages,
	}
}

func (c *CommandHandler) HandleStart(message *tgbotapi.Message) error {
	return nil
}
