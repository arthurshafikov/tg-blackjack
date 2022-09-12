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
	player     = core.Player{Username: username, Cards: core.Cards{}}
	activeGame = core.Game{
		Players: []core.Player{
			{
				Username: "otherUser1",
			},
			{
				Username: username,
				Stop:     true,
			},
			{
				Username: "otherUser2",
				Stop:     true,
			},
		},
	}
)

func getPlayerRepo(mt *mtest.T) (context.Context, *Player) {
	mt.Helper()

	return context.Background(), NewPlayer(mt.Coll)
}

func TestAddNewPlayer(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, playerRepo := getPlayerRepo(mt)
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{}},
			},
		)

		err := playerRepo.AddNewPlayer(ctx, telegramChatID, player)
		require.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		ctx, playerRepo := getPlayerRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch),
		)

		err := playerRepo.AddNewPlayer(ctx, telegramChatID, player)
		require.ErrorIs(t, err, core.ErrNotFound)
	})
}

func TestGetPlayer(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, playerRepo := getPlayerRepo(mt)

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch),
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: primitive.NewObjectID()},
				{Key: "telegram_chat_id", Value: telegramChatID},
				{Key: "active_game", Value: activeGame},
				{Key: "deck", Value: deck},
				{Key: "statistics", Value: nil},
			}),
		)

		result, err := playerRepo.GetPlayer(ctx, telegramChatID, username)
		require.NoError(t, err)
		require.Equal(t, &activeGame.Players[1], result)
	})

	mt.Run("active game exists", func(mt *mtest.T) {
		ctx, playerRepo := getPlayerRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{}),
		)

		_, err := playerRepo.GetPlayer(ctx, telegramChatID, username)
		require.ErrorIs(t, err, core.ErrNoActiveGame)
	})

	mt.Run("player not found", func(mt *mtest.T) {
		ctx, playerRepo := getPlayerRepo(mt)

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch),
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: primitive.NewObjectID()},
				{Key: "telegram_chat_id", Value: telegramChatID},
				{Key: "active_game", Value: activeGame},
				{Key: "deck", Value: deck},
				{Key: "statistics", Value: nil},
			}),
		)

		_, err := playerRepo.GetPlayer(ctx, telegramChatID, "nonExistingUser")
		require.ErrorIs(t, err, core.ErrNotFound)
	})
}

func TestSetPlayerStopAndBusted(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, playerRepo := getPlayerRepo(mt)

		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{}},
			},
		)

		err := playerRepo.SetPlayerStopAndBusted(ctx, telegramChatID, &player)
		require.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		ctx, playerRepo := getPlayerRepo(mt)

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{}),
		)

		err := playerRepo.SetPlayerStopAndBusted(ctx, telegramChatID, &player)
		require.ErrorIs(t, err, core.ErrNotFound)
	})
}
