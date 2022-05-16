package commands

import (
	"errors"
	"fmt"
	"log"

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

func (c *CommandHandler) finishGameIfNeeded(message *tgbotapi.Message) error {
	gameShouldBeFinished, err := c.services.Games.CheckIfGameShouldBeFinished(c.ctx, message.Chat.ID)
	if err != nil {
		return err
	}

	if gameShouldBeFinished {
		game, gameStats, err := c.services.Games.FinishGame(c.ctx, message.Chat.ID)
		if err != nil {
			return err
		}

		msgText := c.messages.GameOver + "\n\n"
		if game.Dealer.IsBlackjack() {
			msgText += c.messages.DealerBlackjack
		} else {
			msgText += c.messages.DealerHand
		}
		msgText += "\n"

		for _, card := range game.Dealer {
			msgText += card.ToString() + " "
		}

		for username, result := range gameStats {
			var resultText string
			switch result {
			case -1:
				resultText = c.messages.Lose
			case 0:
				resultText = c.messages.Push
			case 1:
				resultText = c.messages.Win
			case 2:
				resultText = c.messages.BlackjackResult
			default:
				log.Println("wrong value for result")
			}

			msgText += fmt.Sprintf("\n\n@%s - %s", c.escapeUnderscoreUsername(username), resultText)
		}

		msgText += fmt.Sprintf("\n\n %s", c.messages.GameStartHint)

		return c.sendMessage(tgbotapi.NewMessage(message.Chat.ID, msgText))
	}

	return nil
}
