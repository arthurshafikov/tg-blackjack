package commands

import (
	"errors"
	"fmt"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (c *CommandHandler) HandleStart(message *tgbotapi.Message) error {
	if err := c.services.Chats.RegisterChat(c.ctx, message.Chat.ID); err != nil {
		if errors.Is(err, core.ErrChatRegistered) {
			return fmt.Errorf(c.messages.ChatAlreadyRegistered)
		}

		return err
	}

	msg := c.helper.NewMessage(message.Chat.ID, c.messages.ChatCreatedSuccessfully)

	return c.helper.SendMessage(msg)
}

func (c *CommandHandler) HandleStats(message *tgbotapi.Message) error {
	stats, err := c.services.Statistics.GetStatistics(c.ctx, message.Chat.ID)
	if err != nil {
		return err
	}

	msgText := c.messages.TopPlayers + "\n\n"

	sortedUserStats := stats.SortByValue()
	for pos, userStats := range sortedUserStats {
		msgText += fmt.Sprintf(
			"*%v.* @%s — %v points\n",
			pos+1,
			c.escapeUnderscoreUsername(userStats.Username),
			userStats.Points,
		)
	}
	msg := c.helper.NewMessage(message.Chat.ID, msgText)

	return c.helper.SendMessage(msg)
}
