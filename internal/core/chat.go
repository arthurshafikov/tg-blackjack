package core

const (
	TelegramChatIDField = "telegram_chat_id"
)

type Chat struct {
	TelegramChatID int64           `bson:"telegram_chat_id"`
	ActiveGame     *Game           `bson:"active_game"`
	Statistics     UsersStatistics `bson:"statistics"`
	Deck           *Deck           `bson:"deck"`
}
