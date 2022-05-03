package mongodb

import (
	"context"
	"errors"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Card struct {
	collection *mongo.Collection
}

func NewCard(collection *mongo.Collection) *Card {
	return &Card{
		collection: collection,
	}
}

func (g *Game) AddCardToDealer(ctx context.Context, telegramChatID int64, card core.Card) error {
	filter := bson.M{"telegram_chat_id": telegramChatID}
	update := bson.M{"$push": bson.M{"active_game.dealer": card}}
	if err := g.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

func (g *Game) AddCardToPlayer(ctx context.Context, telegramChatID int64, username string, card core.Card) error {
	filter := bson.M{"$and": bson.A{
		bson.M{"telegram_chat_id": telegramChatID},
		bson.M{"active_game.players.username": username},
	}}
	update := bson.M{"$push": bson.M{"active_game.players.$.cards": card}}
	if err := g.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

func (g *Game) AddNewPlayer(ctx context.Context, telegramChatID int64, player core.Player) error {
	filter := bson.M{"telegram_chat_id": telegramChatID}
	update := bson.M{"$push": bson.M{"active_game.players": player}}
	if err := g.collection.FindOneAndUpdate(ctx, filter, update).Err(); err == nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

// todo test for concurrency
func (g *Game) DrawCard(ctx context.Context, telegramChatID int64) (core.Card, error) {
	var card core.Card

	filter := bson.M{"telegram_chat_id": telegramChatID}
	update := bson.M{"$pop": bson.M{"deck.cards": 1}}

	res := g.collection.FindOneAndUpdate(ctx, filter, update)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return card, core.ErrNotFound
		}

		return card, err
	}

	var chat core.Chat
	if err := res.Decode(&chat); err != nil {
		return card, err
	}

	if chat.Deck.IsEmpty() {
		var err error
		chat.Deck, err = g.setNewDeck(ctx, telegramChatID)
		if err != nil {
			return card, err
		}
	}

	return chat.Deck.DrawCard()
}

func (g *Game) DrawCards(ctx context.Context, telegramChatID int64, amount int) (core.Cards, error) {
	cards := core.Cards{}
	for i := 0; i < amount; i++ {
		card, err := g.DrawCard(ctx, telegramChatID)
		if err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (g *Game) StopDrawing(ctx context.Context, telegramChatID int64, player *core.Player) error {
	filter := bson.M{"$and": bson.A{
		bson.M{"telegram_chat_id": telegramChatID},
		bson.M{"active_game.players.username": player.Username},
	}}
	update := bson.M{"$set": bson.M{
		"active_game.players.$.stop":   player.Stop,
		"active_game.players.$.busted": player.Busted,
	}}
	if err := g.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

func (g *Game) GetPlayer(ctx context.Context, telegramChatID int64, username string) (*core.Player, error) {
	err := g.collection.FindOne(ctx, bson.M{"$and": bson.A{
		bson.M{"telegram_chat_id": telegramChatID},
		bson.M{"active_game": nil},
	}}).Err()
	if err == nil {
		return nil, core.ErrNoActiveGame
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	filter := bson.M{"$and": bson.A{
		bson.M{"telegram_chat_id": telegramChatID},
		bson.M{"active_game.players.username": username},
	}}
	res := g.collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, core.ErrNotFound
		}

		return nil, err
	}

	var chat core.Chat
	if err := res.Decode(&chat); err != nil {
		return nil, err
	}

	for _, player := range chat.ActiveGame.Players {
		if player.Username == username {
			return &player, nil
		}
	}

	return nil, core.ErrNotFound
}

func (g *Game) setNewDeck(ctx context.Context, telegramChatID int64) (core.Deck, error) {
	deck := *core.NewDeck()
	filter := bson.M{"telegram_chat_id": telegramChatID}
	update := bson.M{"$set": bson.M{"deck": deck}}

	res := g.collection.FindOneAndUpdate(ctx, filter, update)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return deck, core.ErrNotFound
		}

		return deck, err
	}

	return deck, nil
}
