package telegram

import (
	"context"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	ctx      context.Context
	bot      *tgbotapi.BotAPI
	services *services.Services

	commandHandler *commands.CommandHandler

	messages config.Messages
}

func NewBot(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	services *services.Services,
	messages config.Messages,
) *Bot {
	commandHandler := commands.NewCommandHandler(ctx, bot, services, messages)

	return &Bot{
		ctx:      ctx,
		bot:      bot,
		services: services,

		commandHandler: commandHandler,

		messages: messages,
	}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID

		if err := b.checkAuthorization(chatID); err != nil && update.Message.Command() != core.StartCommand {
			b.handleError(chatID, err)

			continue
		}

		// Handle commands
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(chatID, err)
			}

			continue
		}
	}

	return nil
}
