package telegram

import (
	"github.com/arthurshafikov/tg-blackjack/internal/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case core.StartCommand:
		return b.commandHandler.HandleStart(message)
	case core.StatsCommand:
		return b.commandHandler.HandleStats(message)
	}

	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	return nil
}

func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) error {
	return nil
}
