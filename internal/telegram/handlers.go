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
	case core.NewGameCommand:
		return b.commandHandler.HandleNewGame(message)
	case core.DrawCardCommand:
		return b.commandHandler.HandleDrawCard(message)
	case core.StopDrawingCommand:
		return b.commandHandler.HandleStopDrawing(message)
	}

	return nil
}
