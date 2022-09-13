package telegram

const MarkDownParseMode = "markdown"

func (b *Bot) handleError(chatID int64, err error) {
	msg := b.helper.NewMessage(chatID, err.Error())
	msg.ParseMode = MarkDownParseMode
	if err := b.helper.SendMessage(msg); err != nil {
		b.logger.Error(err)
	}
}
