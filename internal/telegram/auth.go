package telegram

import (
	"errors"
	"fmt"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
)

func (b *Bot) checkAuthorization(chatID int64) error {
	if err := b.services.Chats.CheckChatExists(b.ctx, chatID); err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return fmt.Errorf(b.messages.ChatNotExists)
		}

		return core.ErrServerError
	}

	return nil
}
