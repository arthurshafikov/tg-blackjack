package commands

import (
	"errors"
	"fmt"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (c *CommandHandler) HandleNewGame(message *tgbotapi.Message) error {
	game, err := c.services.Games.NewGame(c.ctx, message.Chat.ID)
	if err != nil {
		if errors.Is(err, core.ErrActiveGame) {
			return fmt.Errorf(c.messages.ChatHasActiveGame)
		}

		return err
	}

	msgText := c.messages.DealerHand + "\n"
	for i, card := range game.Dealer {
		if i == 0 {
			msgText += "\xE2\x9D\x93 " // question mark
		} else {
			msgText += card.ToString() + " "
		}
	}

	msgText += "\n\n" + c.messages.GameEnterHint

	return c.sendMessage(tgbotapi.NewMessage(message.Chat.ID, msgText))
}
