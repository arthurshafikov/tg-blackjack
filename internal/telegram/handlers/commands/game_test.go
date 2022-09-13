package commands

import (
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	game = &core.Game{
		Dealer: core.Cards{"♣5", "♦Q"},
		// Players: []core.Player{
		// 	{
		// 		Username: username,
		// 		Cards: core.Cards{},

		// 	},
		// },
	}
	gameStats = core.UsersStatistics{
		"someUser1": 1,
		"someUser2": -1,
	}
)

func TestHandleNewGame(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Games.EXPECT().NewGame(c.ctx, chatID).Return(game, nil),
		mockBag.Helper.EXPECT().NewMessage(chatID, "DealerHand\n❓ *♦Q* \n\nGameEnterHint").Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(nil),
	)

	err := c.HandleNewGame(receivedMsg)
	require.NoError(t, err)
}

func TestHandleNewGameActiveGameExists(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Games.EXPECT().NewGame(c.ctx, chatID).Return(nil, core.ErrActiveGame),
	)

	err := c.HandleNewGame(receivedMsg)
	require.ErrorContains(t, err, c.messages.ChatHasActiveGame)
}

func TestHandleNewGameSendMessageReturnsError(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Games.EXPECT().NewGame(c.ctx, chatID).Return(game, nil),
		mockBag.Helper.EXPECT().NewMessage(chatID, "DealerHand\n❓ *♦Q* \n\nGameEnterHint").Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(core.ErrServerError),
		mockBag.Logger.EXPECT().Error(gomock.Any()),
	)

	err := c.HandleNewGame(receivedMsg)
	require.ErrorIs(t, err, core.ErrServerError)
}

func TestFinishGameIfNeeded(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Games.EXPECT().CheckIfGameShouldBeFinished(c.ctx, chatID).Return(true, nil),
		mockBag.Games.EXPECT().FinishGame(c.ctx, chatID).Return(game, gameStats, nil),
		mockBag.Helper.EXPECT().
			NewMessage(chatID, "GameOver\n\nDealerHand\n*♣5* *♦Q* \n\n@someUser1 - Win\n@someUser2 - Lose\n\n GameStartHint").
			Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(nil),
	)

	err := c.finishGameIfNeeded(receivedMsg)
	require.NoError(t, err)
}

func TestFinishGameIfNeededCheckGameShouldBeFinishedReturnsError(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Games.EXPECT().CheckIfGameShouldBeFinished(c.ctx, chatID).Return(false, core.ErrServerError),
	)

	err := c.finishGameIfNeeded(receivedMsg)
	require.ErrorIs(t, err, core.ErrServerError)
}

func TestFinishGameIfNeededSendMessageReturnsError(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Games.EXPECT().CheckIfGameShouldBeFinished(c.ctx, chatID).Return(true, nil),
		mockBag.Games.EXPECT().FinishGame(c.ctx, chatID).Return(game, gameStats, nil),
		mockBag.Helper.EXPECT().
			NewMessage(chatID, "GameOver\n\nDealerHand\n*♣5* *♦Q* \n\n@someUser1 - Win\n@someUser2 - Lose\n\n GameStartHint").
			Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(core.ErrServerError),
		mockBag.Logger.EXPECT().Error(gomock.Any()),
	)

	err := c.finishGameIfNeeded(receivedMsg)
	require.ErrorIs(t, err, core.ErrServerError)
}
