package commands

import (
	"fmt"
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	player = &core.Player{
		Username: username,
		Cards:    core.Cards{"♣5", "♦Q"},
	}
	bustedPlayer = &core.Player{
		Username: username,
		Cards:    core.Cards{"♣5", "♦Q", "♥K"},
		Busted:   true,
		Stop:     true,
	}
	blackjackPlayer = &core.Player{
		Username: username,
		Cards:    core.Cards{"♣A", "♦Q"},
		Stop:     true,
	}
	responseConfig = tgbotapi.NewMessage(chatID, "")
)

func TestHandleDrawCard(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Cards.EXPECT().DrawCardFromDeckToPlayer(c.ctx, chatID, username).Return(player, nil),
		mockBag.Helper.EXPECT().NewMessage(chatID, "PlayerHand someUsername\n*♣5* *♦Q* ").Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(nil),
		mockBag.Games.EXPECT().CheckIfGameShouldBeFinished(c.ctx, chatID).Return(false, nil),
	)

	err := c.HandleDrawCard(receivedMsg)
	require.NoError(t, err)
}

func TestHandleDrawCardBustedPlayer(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Cards.EXPECT().DrawCardFromDeckToPlayer(c.ctx, chatID, username).Return(bustedPlayer, nil),
		mockBag.Helper.EXPECT().
			NewMessage(chatID, "PlayerHand someUsername\n*♣5* *♦Q* *♥K* \nPlayerHandBusted").
			Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(nil),
		mockBag.Games.EXPECT().CheckIfGameShouldBeFinished(c.ctx, chatID).Return(false, nil),
	)

	err := c.HandleDrawCard(receivedMsg)
	require.NoError(t, err)
}

func TestHandleDrawCardBlackjackPlayer(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Cards.EXPECT().DrawCardFromDeckToPlayer(c.ctx, chatID, username).Return(blackjackPlayer, nil),
		mockBag.Helper.EXPECT().NewMessage(chatID, "Blackjack\nPlayerHand someUsername\n*♣A* *♦Q* ").Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(nil),
		mockBag.Games.EXPECT().CheckIfGameShouldBeFinished(c.ctx, chatID).Return(false, nil),
	)

	err := c.HandleDrawCard(receivedMsg)
	require.NoError(t, err)
}

func TestHandleDrawCardDrawCardReturnsNoActiveGame(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Cards.EXPECT().DrawCardFromDeckToPlayer(c.ctx, chatID, username).Return(nil, core.ErrNoActiveGame),
	)

	err := c.HandleDrawCard(receivedMsg)
	require.ErrorContains(t, err, c.messages.ChatHasNoActiveGame)
}

func TestHandleDrawCardBustedPlayerCantDraw(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Cards.EXPECT().DrawCardFromDeckToPlayer(c.ctx, chatID, username).Return(bustedPlayer, core.ErrCantDraw),
		mockBag.Helper.EXPECT().
			NewMessage(chatID, "PlayerCantDraw\nPlayerHand someUsername\n*♣5* *♦Q* *♥K* \nPlayerHandBusted").
			Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(nil),
		mockBag.Games.EXPECT().CheckIfGameShouldBeFinished(c.ctx, chatID).Return(false, nil),
	)

	err := c.HandleDrawCard(receivedMsg)
	require.NoError(t, err)
}

func TestHandleDrawCardSendMessageReturnsError(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	gomock.InOrder(
		mockBag.Cards.EXPECT().DrawCardFromDeckToPlayer(c.ctx, chatID, username).Return(bustedPlayer, nil),
		mockBag.Helper.EXPECT().NewMessage(chatID, gomock.Any()).Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(core.ErrServerError),
		mockBag.Logger.EXPECT().Error(gomock.Any()),
	)

	err := c.HandleDrawCard(receivedMsg)
	require.ErrorIs(t, err, core.ErrServerError)
}

func TestHandleStopDrawing(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	player := *player
	player.Cards = nil
	gomock.InOrder(
		mockBag.Players.EXPECT().StopDrawing(c.ctx, chatID, &player).Return(nil),
		mockBag.Helper.EXPECT().NewMessage(chatID, "StoppedDrawing someUsername\n").Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(nil),
		mockBag.Games.EXPECT().CheckIfGameShouldBeFinished(c.ctx, chatID).Return(false, nil),
	)

	err := c.HandleStopDrawing(receivedMsg)
	require.NoError(t, err)
}

func TestHandleStopDrawingBustedPlayer(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	player := *player
	player.Cards = nil
	gomock.InOrder(
		mockBag.Players.EXPECT().StopDrawing(c.ctx, chatID, &player).Return(core.ErrBusted),
	)

	err := c.HandleStopDrawing(receivedMsg)
	require.ErrorContains(t, err, fmt.Sprintf(c.messages.PlayerAlreadyBusted, username))
}

func TestHandleStopDrawingCantDrawPlayer(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	player := *player
	player.Cards = nil
	gomock.InOrder(
		mockBag.Players.EXPECT().StopDrawing(c.ctx, chatID, &player).Return(core.ErrAlreadyStopped),
	)

	err := c.HandleStopDrawing(receivedMsg)
	require.ErrorContains(t, err, fmt.Sprintf(c.messages.PlayerAlreadyStopped, username))
}

func TestHandleStopDrawingNoActiveGame(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	player := *player
	player.Cards = nil
	gomock.InOrder(
		mockBag.Players.EXPECT().StopDrawing(c.ctx, chatID, &player).Return(core.ErrNoActiveGame),
	)

	err := c.HandleStopDrawing(receivedMsg)
	require.ErrorContains(t, err, c.messages.ChatHasNoActiveGame)
}

func TestHandleStopDrawingNotFound(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	player := *player
	player.Cards = nil
	gomock.InOrder(
		mockBag.Players.EXPECT().StopDrawing(c.ctx, chatID, &player).Return(core.ErrNotFound),
	)

	err := c.HandleStopDrawing(receivedMsg)
	require.ErrorContains(t, err, c.messages.GameEnterHint)
}

func TestHandleStopDrawingSendMessageReturnsError(t *testing.T) {
	handlerParams, mockBag := handlers.GetHandlerParamsWithMocks(t)
	c := NewCommandHandler(handlerParams)
	player := *player
	player.Cards = nil
	gomock.InOrder(
		mockBag.Players.EXPECT().StopDrawing(c.ctx, chatID, &player).Return(nil),
		mockBag.Helper.EXPECT().NewMessage(chatID, "StoppedDrawing someUsername\n").Return(responseConfig),
		mockBag.Helper.EXPECT().SendMessage(responseConfig).Return(core.ErrServerError),
		mockBag.Logger.EXPECT().Error(gomock.Any()),
	)

	err := c.HandleStopDrawing(receivedMsg)
	require.ErrorIs(t, err, core.ErrServerError)
}
