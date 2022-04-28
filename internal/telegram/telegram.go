package telegram

import (
	"context"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers/commands"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers/queries"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	ctx      context.Context
	bot      *tgbotapi.BotAPI
	services *services.Services

	commandHandler *commands.CommandHandler
	queryHandler   *queries.QueryHandler

	messages config.Messages
}

func NewBot(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	services *services.Services,
	messages config.Messages,
) *Bot {
	commandHandler := commands.NewCommandHandler(ctx, bot, services, messages)
	queryHandler := queries.NewQueryHandler(ctx, bot, services, messages)

	return &Bot{
		ctx:      ctx,
		bot:      bot,
		services: services,

		commandHandler: commandHandler,
		queryHandler:   queryHandler,

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

	var chatID int64
	for update := range updates {
		if update.Message != nil { //nolint
			chatID = update.Message.Chat.ID
		} else if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
		} else {
			continue
		}

		if err := b.checkAuthorization(chatID); err != nil && update.Message.Command() != core.StartCommand {
			b.handleError(chatID, err)

			continue
		}

		if update.CallbackQuery != nil {
			if err := b.handleCallbackQuery(update.CallbackQuery); err != nil {
				b.handleError(chatID, err)
			}

			continue
		}

		// Handle commands
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(chatID, err)
			}

			continue
		}

		// Handle regular messages
		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(chatID, err)
		}
	}

	return nil
}
