package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/core"
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
	if err := c.services.Chats.RegisterChat(c.ctx, message.Chat.ID); err != nil {
		if errors.Is(err, core.ErrChatRegistered) {
			return fmt.Errorf(c.messages.ChatAlreadyRegistered)
		}

		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, c.messages.ChatCreatedSuccessfully)

	return c.sendMessage(msg)
}

func (c *CommandHandler) sendMessage(msg tgbotapi.MessageConfig) error {
	msg.ParseMode = "markdown"
	_, err := c.bot.Send(msg)

	return err
}

func (c *CommandHandler) escapeUnderscoreUsername(username string) string {
	return strings.Replace(username, "_", "\\_", -1)
}
