package telegram

func (b *Bot) handleError(chatID int64, err error) {
	msg := b.helper.NewMessage(chatID, err.Error())
	msg.ParseMode = "markdown"
	if err := b.helper.SendMessage(msg); err != nil {
		b.logger.Error(err)
	}
}
