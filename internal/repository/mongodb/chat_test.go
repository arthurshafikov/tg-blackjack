package mongodb

import (
	"context"
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func getChatRepo(mt *mtest.T) (context.Context, *Chat) {
	mt.Helper()

	return context.Background(), NewChat(mt.Coll)
}

func TestCheckChatExists(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, chatRepo := getChatRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{}),
		)

		err := chatRepo.CheckChatExists(ctx, telegramChatID)
		require.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		ctx, chatRepo := getChatRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch),
		)

		err := chatRepo.CheckChatExists(ctx, telegramChatID)
		require.ErrorIs(t, err, core.ErrNotFound)
	})
}

func TestRegisterChat(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		ctx, chatRepo := getChatRepo(mt)
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch, bson.D{}),
		)

		err := chatRepo.RegisterChat(ctx, telegramChatID)
		require.NoError(t, err)
	})
}
