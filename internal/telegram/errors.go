package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, err.Error())
	msg.ParseMode = "markdown"
	b.bot.Send(msg)
}
