package repository

import (
	"context"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Cards interface {
	AddCardToDealer(ctx context.Context, telegramChatID int64, card core.Card) error
	AddCardToPlayer(ctx context.Context, telegramChatID int64, username string, card core.Card) error
	DrawCardFromDeck(ctx context.Context, telegramChatID int64) (core.Card, error)
	SetNewDeck(ctx context.Context, telegramChatID int64, deck *core.Deck) error
}

type Chats interface {
	CheckChatExists(ctx context.Context, telegramChatID int64) error
	RegisterChat(ctx context.Context, telegramChatID int64) error
}

type Games interface {
	GetActiveGame(ctx context.Context, telegramChatID int64) (*core.Game, error)
	NullActiveGame(ctx context.Context, telegramChatID int64) error
	SetActiveGame(ctx context.Context, telegramChatID int64, game core.Game) error
}

type Players interface {
	AddNewPlayer(ctx context.Context, telegramChatID int64, player core.Player) error
	GetPlayer(ctx context.Context, telegramChatID int64, username string) (*core.Player, error)
	SetPlayerStopAndBusted(ctx context.Context, telegramChatID int64, player *core.Player) error
}

type Statistic interface {
	GetStatistics(ctx context.Context, telegramChatID int64) (core.UsersStatistics, error)
	SetStatistics(ctx context.Context, telegramChatID int64, stats core.UsersStatistics) error
}

type Repository struct {
	Cards
	Chats
	Games
	Players
	Statistic
}

func NewRepository(db *mongo.Client) *Repository {
	chatsCollection := db.Database("homestead").Collection("chats")

	return &Repository{
		Cards:     mongodb.NewCard(chatsCollection),
		Chats:     mongodb.NewChat(chatsCollection),
		Games:     mongodb.NewGame(chatsCollection),
		Players:   mongodb.NewPlayer(chatsCollection),
		Statistic: mongodb.NewStatistic(chatsCollection),
	}
}
