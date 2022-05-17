package commands

import (
	"errors"
	"fmt"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (c *CommandHandler) HandleDrawCard(message *tgbotapi.Message) error {
	var msgText string
	player, err := c.services.Cards.DrawCardFromDeckToPlayer(c.ctx, message.Chat.ID, message.From.UserName)
	if err != nil {
		switch err { //nolint
		case core.ErrCantDraw:
			msgText += c.messages.PlayerCantDraw + "\n"
		case core.ErrNoActiveGame:
			return fmt.Errorf(c.messages.ChatHasNoActiveGame)
		case core.ErrBusted:
			break
		default:
			return err
		}
	}

	if player.Cards.IsBlackjack() {
		msgText += c.messages.Blackjack + "\n"
	}

	msgText += fmt.Sprintf(c.messages.PlayerHand+"\n", c.escapeUnderscoreUsername(message.From.UserName))

	for _, card := range player.Cards {
		msgText += card.ToString() + " "
	}

	if player.Busted {
		msgText += "\n" + c.messages.PlayerHandBusted
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)

	if err := c.sendMessage(msg); err != nil {
		return err
	}

	if err := c.finishGameIfNeeded(message); err != nil {
		return err
	}

	return nil
}

func (c *CommandHandler) HandleStopDrawing(message *tgbotapi.Message) error {
	var msgText string

	player := core.Player{
		Username: message.From.UserName,
	}
	if err := c.services.Players.StopDrawing(c.ctx, message.Chat.ID, &player); err != nil {
		if errors.Is(err, core.ErrAlreadyStopped) {
			return fmt.Errorf(c.messages.PlayerAlreadyStopped, c.escapeUnderscoreUsername(player.Username))
		}
		if errors.Is(err, core.ErrBusted) {
			return fmt.Errorf(c.messages.PlayerAlreadyBusted, c.escapeUnderscoreUsername(player.Username))
		}
		if errors.Is(err, core.ErrNoActiveGame) {
			return fmt.Errorf(c.messages.ChatHasNoActiveGame)
		}
		if errors.Is(err, core.ErrNotFound) {
			return fmt.Errorf(c.messages.GameEnterHint)
		}

		return err
	}

	msgText += fmt.Sprintf(c.messages.StoppedDrawing+"\n", c.escapeUnderscoreUsername(message.From.UserName))

	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)

	if err := c.sendMessage(msg); err != nil {
		return err
	}

	if err := c.finishGameIfNeeded(message); err != nil {
		return err
	}

	return nil
}
