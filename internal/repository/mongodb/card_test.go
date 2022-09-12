package mongodb

import (
	"context"
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var (
	telegramChatID int64     = 64
	card           core.Card = "♣4"
	username                 = "someUser"
	deck                     = &core.Deck{Cards: core.Cards{"♦4", "♦10", "♦K"}}
)

func getCardRepo(mt *mtest.T) (context.Context, *Card) {
	mt.Helper()

	return context.Background(), NewCard(mt.Coll)
}

func TestAddCardToDealer(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, cardRepo := getCardRepo(mt)
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{}},
			},
		)

		err := cardRepo.AddCardToDealer(ctx, telegramChatID, card)
		require.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		ctx, cardRepo := getCardRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{}),
		)

		err := cardRepo.AddCardToDealer(ctx, telegramChatID, card)
		require.ErrorIs(t, err, core.ErrNotFound)
	})
}

func TestAddCardToPlayer(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, cardRepo := getCardRepo(mt)
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{}},
			},
		)

		err := cardRepo.AddCardToPlayer(ctx, telegramChatID, username, card)
		require.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		ctx, cardRepo := getCardRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{}),
		)

		err := cardRepo.AddCardToPlayer(ctx, telegramChatID, username, card)
		require.ErrorIs(t, err, core.ErrNotFound)
	})
}

func TestDrawCardFromDeck(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, cardRepo := getCardRepo(mt)
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: "telegram_chat_id", Value: telegramChatID},
					{Key: "active_game", Value: nil},
					{Key: "deck", Value: deck},
					{Key: "statistics", Value: nil},
				}},
			},
		)

		resultCard, err := cardRepo.DrawCardFromDeck(ctx, telegramChatID)
		require.NoError(t, err)
		require.Equal(t, core.Card("♦K"), resultCard)
	})

	mt.Run("not found", func(mt *mtest.T) {
		ctx, cardRepo := getCardRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{}),
		)

		resultCard, err := cardRepo.DrawCardFromDeck(ctx, telegramChatID)
		require.ErrorIs(t, err, core.ErrNotFound)
		require.Equal(t, core.Card(""), resultCard)
	})
}

func TestSetNewDeck(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, cardRepo := getCardRepo(mt)
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{}},
			},
		)

		err := cardRepo.SetNewDeck(ctx, telegramChatID, core.NewDeck(1))
		require.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		ctx, cardRepo := getCardRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{}),
		)

		err := cardRepo.SetNewDeck(ctx, telegramChatID, core.NewDeck(1))
		require.ErrorIs(t, err, core.ErrNotFound)
	})
}
