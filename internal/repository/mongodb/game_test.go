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

var expectedGame = &core.Game{Dealer: core.Cards{"♠10", "♠5"}}

func getGameRepo(mt *mtest.T) (context.Context, *Game) {
	mt.Helper()

	return context.Background(), NewGame(mt.Coll)
}

func TestGetActiveGame(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, gameRepo := getGameRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: primitive.NewObjectID()},
				{Key: "telegram_chat_id", Value: telegramChatID},
				{Key: "active_game", Value: expectedGame},
				{Key: "deck", Value: nil},
				{Key: "statistics", Value: nil},
			}),
		)

		game, err := gameRepo.GetActiveGame(ctx, telegramChatID)
		require.NoError(t, err)
		require.Equal(t, expectedGame, game)
	})

	mt.Run("not found", func(mt *mtest.T) {
		ctx, gameRepo := getGameRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch),
		)

		game, err := gameRepo.GetActiveGame(ctx, telegramChatID)
		require.ErrorIs(t, err, core.ErrNotFound)
		require.Nil(t, game)
	})
}

func TestNullActiveGame(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, gameRepo := getGameRepo(mt)
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{}},
			},
		)

		err := gameRepo.NullActiveGame(ctx, telegramChatID)
		require.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		ctx, gameRepo := getGameRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch),
		)

		err := gameRepo.NullActiveGame(ctx, telegramChatID)
		require.ErrorIs(t, err, core.ErrNotFound)
	})
}

func TestSetActiveGame(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, gameRepo := getGameRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{}),
			bson.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{}},
			},
		)

		err := gameRepo.SetActiveGame(ctx, telegramChatID, *expectedGame)
		require.NoError(t, err)
	})

	mt.Run("active game exists", func(mt *mtest.T) {
		ctx, gameRepo := getGameRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: primitive.NewObjectID()},
				{Key: "telegram_chat_id", Value: telegramChatID},
				{Key: "active_game", Value: expectedGame},
				{Key: "deck", Value: nil},
				{Key: "statistics", Value: nil},
			}),
		)

		err := gameRepo.SetActiveGame(ctx, telegramChatID, *expectedGame)
		require.ErrorIs(t, err, core.ErrActiveGame)
	})

	mt.Run("not found", func(mt *mtest.T) {
		ctx, gameRepo := getGameRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch),
		)

		err := gameRepo.SetActiveGame(ctx, telegramChatID, *expectedGame)
		require.ErrorIs(t, err, core.ErrNotFound)
	})
}
