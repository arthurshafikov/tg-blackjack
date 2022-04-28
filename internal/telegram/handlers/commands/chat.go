package commands

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (c *CommandHandler) HandleStats(message *tgbotapi.Message) error {
	stats, err := c.services.Chats.GetStatistics(c.ctx, message.Chat.ID)
	if err != nil {
		return err
	}

	msgText := c.messages.TopPlayers + "\n\n"

	sortedUserStats := stats.SortByValue()
	for pos, userStats := range sortedUserStats {
		msgText += fmt.Sprintf("*%v.* %s — %v points\n", pos+1, userStats.Username, userStats.Points)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)

	return c.sendMessage(msg)
}
