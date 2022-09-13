package telegram

import (
	"context"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	ctx      context.Context
	services *services.Services
	logger   services.Logger
	messages config.Messages

	commandHandler CommandHandler

	helper handlers.TelegramHandlerHelper
}

type CommandHandler interface {
	HandleStart(message *tgbotapi.Message) error
	HandleStats(message *tgbotapi.Message) error
	HandleNewGame(message *tgbotapi.Message) error
	HandleDrawCard(message *tgbotapi.Message) error
	HandleStopDrawing(message *tgbotapi.Message) error
}

type Deps struct {
	Ctx      context.Context
	Services *services.Services
	Logger   services.Logger
	Messages config.Messages

	CommandHandler CommandHandler

	Helper handlers.TelegramHandlerHelper
}

func NewBot(deps *Deps) *Bot {
	return &Bot{
		ctx:      deps.Ctx,
		services: deps.Services,
		logger:   deps.Logger,
		messages: deps.Messages,

		commandHandler: deps.CommandHandler,

		helper: deps.Helper,
	}
}

func (b *Bot) Start() error {
	u := b.helper.NewUpdateChannel(0)
	u.Timeout = 60

	updates, err := b.helper.GetUpdatesChan(u)
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
