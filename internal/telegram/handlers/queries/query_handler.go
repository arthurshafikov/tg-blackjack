package queries

import (
	"context"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type QueryHandler struct {
	ctx      context.Context
	bot      *tgbotapi.BotAPI
	services *services.Services
	messages config.Messages
}

func NewQueryHandler(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	services *services.Services,
	messages config.Messages,
) *QueryHandler {
	return &QueryHandler{
		ctx:      ctx,
		bot:      bot,
		services: services,
		messages: messages,
	}
}

func (q *QueryHandler) sendMessage(msg tgbotapi.MessageConfig) error {
	msg.ParseMode = "markdown"
	_, err := q.bot.Send(msg)

	return err
}
