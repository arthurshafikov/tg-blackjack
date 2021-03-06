package mongodb

import (
	"context"
	"errors"
	"fmt"

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

func (c *Card) AddCardToDealer(ctx context.Context, telegramChatID int64, card core.Card) error {
	filter := bson.M{core.TelegramChatIDField: telegramChatID}
	update := bson.M{"$push": bson.M{fmt.Sprintf("%s.%s", core.ActiveGameField, core.GameDealerField): card}}

	if err := c.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

func (c *Card) AddCardToPlayer(ctx context.Context, telegramChatID int64, username string, card core.Card) error {
	filter := bson.M{"$and": bson.A{
		bson.M{core.TelegramChatIDField: telegramChatID},
		bson.M{
			fmt.Sprintf(
				"%s.%s.%s",
				core.ActiveGameField,
				core.GamePlayersField,
				core.GamePlayerUsernameField,
			): username,
		},
	}}
	update := bson.M{"$push": bson.M{
		fmt.Sprintf(
			"%s.%s.$.%s",
			core.ActiveGameField,
			core.GamePlayersField,
			core.GamePlayerCardsField,
		): card,
	}}

	if err := c.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

func (c *Card) DrawCardFromDeck(ctx context.Context, telegramChatID int64) (core.Card, error) {
	var card core.Card

	filter := bson.M{core.TelegramChatIDField: telegramChatID}
	update := bson.M{"$pop": bson.M{fmt.Sprintf("%s.%s", core.DeckField, core.GamePlayerCardsField): 1}}

	res := c.collection.FindOneAndUpdate(ctx, filter, update)
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

	return chat.Deck.DrawCard()
}

func (c *Card) SetNewDeck(ctx context.Context, telegramChatID int64, deck *core.Deck) error {
	filter := bson.M{core.TelegramChatIDField: telegramChatID}
	update := bson.M{"$set": bson.M{core.DeckField: deck}}

	if err := c.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}
