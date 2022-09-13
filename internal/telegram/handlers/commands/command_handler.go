package commands

import (
	"context"
	"strings"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers"
)

type CommandHandler struct {
	ctx      context.Context
	services *services.Services

	messages config.Messages

	helper handlers.TelegramHandlerHelper
}

func NewCommandHandler(handlerParams handlers.HandlerParams) *CommandHandler {
	return &CommandHandler{
		ctx:      handlerParams.Ctx,
		services: handlerParams.Services,
		messages: handlerParams.Messages,

		helper: handlerParams.Helper,
	}
}

func (c *CommandHandler) escapeUnderscoreUsername(username string) string {
	return strings.ReplaceAll(username, "_", "\\_")
}
