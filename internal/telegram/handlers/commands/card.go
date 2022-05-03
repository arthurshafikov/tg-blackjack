package commands

import (
	"errors"
	"fmt"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (c *CommandHandler) HandleDrawCard(message *tgbotapi.Message) error {
	var msgText string
	player, err := c.services.Cards.DrawCard(c.ctx, message.Chat.ID, message.From.UserName)
	if err != nil {
		if errors.Is(err, core.ErrCantDraw) {
			msgText += c.messages.PlayerCantDraw + "\n"
		} else if !errors.Is(err, core.ErrBusted) {
			return err
		}
	}

	msgText += fmt.Sprintf(c.messages.PlayerHand+"\n", c.escapeUnderscoreUsername(message.From.UserName))

	for _, card := range player.Cards {
		msgText += card.ToString() + " "
	}

	if player.Busted {
		msgText += "\n" + c.messages.PlayerHandBusted
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyToMessageID = message.MessageID

	if err := c.sendMessage(msg); err != nil {
		return err
	}

	if err := c.finishGameIfNeeded(message); err != nil {
		return err
	}

	return nil
}
