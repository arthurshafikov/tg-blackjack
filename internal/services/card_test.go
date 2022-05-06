package services

import (
	"context"
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
	mock_repository "github.com/arthurshafikov/tg-blackjack/internal/repository/mocks"
	mock_services "github.com/arthurshafikov/tg-blackjack/internal/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	telegramChatID int64     = 1231234
	username                 = "some_user"
	cards                    = core.Cards{"♣10", "♥2"}
	blackjack                = core.Cards{"♣10", "♥A"}
	newCard        core.Card = "♣5"
	newKCard       core.Card = "♣K"
	new9Card       core.Card = "♣9"
)

func getCardServiceDependencies(t *testing.T) (
	context.Context,
	*mock_services.MockLogger,
	*mock_repository.MockCards,
	*mock_services.MockPlayers,
) {
	t.Helper()
	ctrl := gomock.NewController(t)

	return context.Background(),
		mock_services.NewMockLogger(ctrl),
		mock_repository.NewMockCards(ctrl),
		mock_services.NewMockPlayers(ctrl)
}

func TestDrawCard(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	expected := &core.Player{
		Username: username,
		Cards:    cards,
		Stop:     false,
		Busted:   false,
	}
	gomock.InOrder(
		playerServiceMock.EXPECT().GetPlayer(ctx, telegramChatID, username).Return(&core.Player{
			Username: username,
			Cards:    cards,
		}, nil),
		repo.EXPECT().DrawCard(ctx, telegramChatID).Return(newCard, nil),
		repo.EXPECT().AddCardToPlayer(ctx, telegramChatID, username, newCard).Return(nil),
	)
	expected.Cards = append(expected.Cards, newCard)
	service := NewCardService(logger, repo, playerServiceMock)

	player, err := service.DrawCard(ctx, telegramChatID, username)

	require.NoError(t, err)
	require.Equal(t, expected, player)
}

func TestDrawCardPlayerNotExists(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	expected := core.Player{
		Username: username,
		Cards:    cards,
		Stop:     false,
		Busted:   false,
	}
	gomock.InOrder(
		playerServiceMock.EXPECT().GetPlayer(ctx, telegramChatID, username).Return(nil, core.ErrNotFound),
		repo.EXPECT().DrawCards(ctx, telegramChatID, 2).Return(cards, nil),
		playerServiceMock.EXPECT().AddNewPlayer(ctx, telegramChatID, expected).Return(nil),
	)
	service := NewCardService(logger, repo, playerServiceMock)

	player, err := service.DrawCard(ctx, telegramChatID, username)

	require.NoError(t, err)
	require.Equal(t, expected, *player)
}

func TestDrawCardGameNotExists(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	var expected *core.Player
	gomock.InOrder(
		playerServiceMock.EXPECT().GetPlayer(ctx, telegramChatID, username).Return(nil, core.ErrNoActiveGame),
	)
	service := NewCardService(logger, repo, playerServiceMock)

	player, err := service.DrawCard(ctx, telegramChatID, username)

	require.ErrorIs(t, err, core.ErrNoActiveGame)
	require.Equal(t, expected, player)
}

func TestDrawCardServerError(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	var expected *core.Player
	gomock.InOrder(
		playerServiceMock.EXPECT().GetPlayer(ctx, telegramChatID, username).Return(nil, core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewCardService(logger, repo, playerServiceMock)

	player, err := service.DrawCard(ctx, telegramChatID, username)

	require.ErrorIs(t, err, core.ErrServerError)
	require.Equal(t, expected, player)
}

func TestDrawCardBustedCase(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	expected := &core.Player{
		Username: username,
		Cards:    append(cards, core.Cards{newCard, newKCard}...),
		Stop:     true,
		Busted:   true,
	}
	player := &core.Player{
		Username: username,
		Cards:    append(cards, newCard),
	}
	gomock.InOrder(
		playerServiceMock.EXPECT().GetPlayer(ctx, telegramChatID, username).Return(player, nil),
		repo.EXPECT().DrawCard(ctx, telegramChatID).Return(newKCard, nil),
		repo.EXPECT().AddCardToPlayer(ctx, telegramChatID, username, newKCard).Return(nil),
		playerServiceMock.EXPECT().StopDrawing(ctx, telegramChatID, player).Return(nil),
	)
	service := NewCardService(logger, repo, playerServiceMock)

	player, err := service.DrawCard(ctx, telegramChatID, username)

	require.ErrorIs(t, err, core.ErrBusted)
	require.Equal(t, expected, player)
}

func TestDrawCard21ValueCase(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	expected := &core.Player{
		Username: username,
		Cards:    append(cards, new9Card),
		Stop:     true,
		Busted:   false,
	}
	player := &core.Player{
		Username: username,
		Cards:    cards,
	}
	gomock.InOrder(
		playerServiceMock.EXPECT().GetPlayer(ctx, telegramChatID, username).Return(player, nil),
		repo.EXPECT().DrawCard(ctx, telegramChatID).Return(new9Card, nil),
		repo.EXPECT().AddCardToPlayer(ctx, telegramChatID, username, new9Card).Return(nil),
		playerServiceMock.EXPECT().StopDrawing(ctx, telegramChatID, player).Return(nil),
	)
	service := NewCardService(logger, repo, playerServiceMock)

	player, err := service.DrawCard(ctx, telegramChatID, username)

	require.NoError(t, err)
	require.Equal(t, expected, player)
	require.Equal(t, 21, player.Cards.CountValue())
}

func TestDrawCardBlackjackCase(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	expected := &core.Player{
		Username: username,
		Cards:    blackjack,
		Stop:     true,
		Busted:   false,
	}
	gomock.InOrder(
		playerServiceMock.EXPECT().GetPlayer(ctx, telegramChatID, username).Return(nil, core.ErrNotFound),
		repo.EXPECT().DrawCards(ctx, telegramChatID, 2).Return(blackjack, nil),
		playerServiceMock.EXPECT().AddNewPlayer(ctx, telegramChatID, *expected).Return(nil),
	)
	service := NewCardService(logger, repo, playerServiceMock)

	player, err := service.DrawCard(ctx, telegramChatID, username)

	require.NoError(t, err)
	require.Equal(t, expected, player)
	require.Equal(t, 21, player.Cards.CountValue())
	require.True(t, player.Cards.IsBlackjack())
}

func TestDrawCardFromDeckToDealer(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().DrawCard(ctx, telegramChatID).Return(new9Card, nil),
		repo.EXPECT().AddCardToDealer(ctx, telegramChatID, new9Card).Return(nil),
	)
	service := NewCardService(logger, repo, playerServiceMock)

	card, err := service.DrawCardFromDeckToDealer(ctx, telegramChatID)

	require.NoError(t, err)
	require.Equal(t, new9Card, card)
}

func TestDrawCardFromDeckToDealerServerError(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().DrawCard(ctx, telegramChatID).Return(new9Card, core.ErrNotFound),
		logger.EXPECT().Error(core.ErrNotFound),
	)
	service := NewCardService(logger, repo, playerServiceMock)

	card, err := service.DrawCardFromDeckToDealer(ctx, telegramChatID)

	require.ErrorIs(t, err, core.ErrServerError)
	require.Equal(t, new9Card, card)
}
