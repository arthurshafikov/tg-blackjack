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

var statistics = core.UsersStatistics{
	"someUser1": 11,
	"someUser2": -34,
	"someUser3": 2,
}

func getStatisticRepo(mt *mtest.T) (context.Context, *Statistic) {
	mt.Helper()

	return context.Background(), NewStatistic(mt.Coll)
}

func TestGetStatistics(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, statisticRepo := getStatisticRepo(mt)
		statistics := core.UsersStatistics{
			"someUser1": 11,
			"someUser2": -34,
			"someUser3": 2,
		}
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: primitive.NewObjectID()},
				{Key: "telegram_chat_id", Value: telegramChatID},
				{Key: "active_game", Value: nil},
				{Key: "deck", Value: deck},
				{Key: "statistics", Value: statistics},
			}),
		)

		result, err := statisticRepo.GetStatistics(ctx, telegramChatID)
		require.NoError(t, err)
		require.Equal(t, statistics, result)
	})

	mt.Run("not found", func(mt *mtest.T) {
		ctx, statisticRepo := getStatisticRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch),
		)

		result, err := statisticRepo.GetStatistics(ctx, telegramChatID)
		require.ErrorIs(t, err, core.ErrNotFound)
		require.Nil(t, result)
	})
}

func TestSetStatistics(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, statisticRepo := getStatisticRepo(mt)
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 1},
				{Key: "value", Value: bson.D{}},
			},
		)

		err := statisticRepo.SetStatistics(ctx, telegramChatID, statistics)
		require.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		ctx, statisticRepo := getStatisticRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch),
		)

		err := statisticRepo.SetStatistics(ctx, telegramChatID, statistics)
		require.ErrorIs(t, err, core.ErrNotFound)
	})
}
