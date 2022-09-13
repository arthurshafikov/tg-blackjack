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

	if err := c.helper.SendMessage(msg); err != nil {
		c.logger.Error(fmt.Errorf("chatId: %v, func: CommandHandler.HandleStart, error: %w", message.Chat.ID, err))

		return core.ErrServerError
	}

	return nil
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

	if err := c.helper.SendMessage(msg); err != nil {
		c.logger.Error(fmt.Errorf("chatId: %v, func: CommandHandler.HandleStats, error: %w", message.Chat.ID, err))

		return core.ErrServerError
	}

	return nil
}
