package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Player struct {
	collection *mongo.Collection
}

func NewPlayer(collection *mongo.Collection) *Player {
	return &Player{
		collection: collection,
	}
}

func (p *Player) AddNewPlayer(ctx context.Context, telegramChatID int64, player core.Player) error {
	filter := bson.M{core.TelegramChatIDField: telegramChatID}
	update := bson.M{"$push": bson.M{fmt.Sprintf("%s.%s", core.ActiveGameField, core.GamePlayersField): player}}

	if err := p.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}

func (p *Player) GetPlayer(ctx context.Context, telegramChatID int64, username string) (*core.Player, error) {
	err := p.collection.FindOne(ctx, bson.M{"$and": bson.A{
		bson.M{core.TelegramChatIDField: telegramChatID},
		bson.M{core.ActiveGameField: nil},
	}}).Err()
	if err == nil {
		return nil, core.ErrNoActiveGame
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

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
	res := p.collection.FindOne(ctx, filter)
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

func (p *Player) SetPlayerStopAndBusted(ctx context.Context, telegramChatID int64, player *core.Player) error {
	filter := bson.M{"$and": bson.A{
		bson.M{core.TelegramChatIDField: telegramChatID},
		bson.M{
			fmt.Sprintf(
				"%s.%s.%s",
				core.ActiveGameField,
				core.GamePlayersField,
				core.GamePlayerUsernameField,
			): player.Username,
		},
	}}
	update := bson.M{"$set": bson.M{
		fmt.Sprintf(
			"%s.%s.$.%s",
			core.ActiveGameField,
			core.GamePlayersField,
			core.GamePlayerStopField,
		): player.Stop,
		fmt.Sprintf(
			"%s.%s.$.%s",
			core.ActiveGameField,
			core.GamePlayersField,
			core.GamePlayerBustedField,
		): player.Busted,
	}}
	if err := p.collection.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}

		return err
	}

	return nil
}
