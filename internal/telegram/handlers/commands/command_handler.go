package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

func (c *CommandHandler) HandleStart(message *tgbotapi.Message) error {
	if err := c.services.Chats.RegisterChat(c.ctx, message.Chat.ID); err != nil {
		if errors.Is(err, core.ErrChatRegistered) {
			return fmt.Errorf(c.messages.ChatAlreadyRegistered)
		}

		return err
	}

	msg := c.helper.NewMessage(message.Chat.ID, c.messages.ChatCreatedSuccessfully)

	return c.helper.SendMessage(msg)
}

func (c *CommandHandler) escapeUnderscoreUsername(username string) string {
	return strings.ReplaceAll(username, "_", "\\_")
}
