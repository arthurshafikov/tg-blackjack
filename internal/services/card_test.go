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
	new2Card       core.Card = "♥2"
	new5Card       core.Card = "♣5"
	newKCard       core.Card = "♣K"
	newACard       core.Card = "♥A"
	new9Card       core.Card = "♣9"
	cards                    = core.Cards{newKCard, new2Card}
	blackjack                = core.Cards{newACard, newKCard}
	sameCards                = core.Cards{new5Card, new5Card}
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

func TestDrawCards(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().DrawCardFromDeck(ctx, telegramChatID).Times(4).Return(new9Card, nil),
	)
	service := NewCardService(logger, repo, playerServiceMock, 1)

	cards, err := service.drawCards(ctx, telegramChatID, 4)

	require.NoError(t, err)
	require.Len(t, cards, 4)
}

func TestDrawCardFromDeckToDealer(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().DrawCardFromDeck(ctx, telegramChatID).Return(new9Card, nil),
		repo.EXPECT().AddCardToDealer(ctx, telegramChatID, new9Card).Return(nil),
	)
	service := NewCardService(logger, repo, playerServiceMock, 1)

	card, err := service.DrawCardFromDeckToDealer(ctx, telegramChatID)

	require.NoError(t, err)
	require.Equal(t, new9Card, card)
}

func TestDrawCardFromDeckToDealerServerError(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().DrawCardFromDeck(ctx, telegramChatID).Return(new9Card, core.ErrNotFound),
		logger.EXPECT().Error(core.ErrNotFound),
	)
	service := NewCardService(logger, repo, playerServiceMock, 1)

	card, err := service.DrawCardFromDeckToDealer(ctx, telegramChatID)

	require.ErrorIs(t, err, core.ErrServerError)
	require.Equal(t, new9Card, card)
}

func TestDrawCardDeckEmpty(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().DrawCardFromDeck(ctx, telegramChatID).Return(new9Card, core.ErrDeckEmpty),
		repo.EXPECT().SetNewDeck(ctx, telegramChatID, gomock.Any()).Return(nil),
	)
	service := NewCardService(logger, repo, playerServiceMock, 1)

	_, err := service.drawCard(ctx, telegramChatID)

	require.NoError(t, err)
}

func TestDrawCardNonEmptyError(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().DrawCardFromDeck(ctx, telegramChatID).Return(new9Card, core.ErrServerError),
	)
	service := NewCardService(logger, repo, playerServiceMock, 1)

	_, err := service.drawCard(ctx, telegramChatID)

	require.ErrorIs(t, err, core.ErrServerError)
}

func TestDrawCardFromDeckToPlayer(t *testing.T) {
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
		repo.EXPECT().DrawCardFromDeck(ctx, telegramChatID).Return(new5Card, nil),
		repo.EXPECT().AddCardToPlayer(ctx, telegramChatID, username, new5Card).Return(nil),
	)
	expected.Cards = append(expected.Cards, new5Card)
	service := NewCardService(logger, repo, playerServiceMock, 1)

	player, err := service.DrawCardFromDeckToPlayer(ctx, telegramChatID, username)

	require.NoError(t, err)
	require.Equal(t, expected, player)
}

func TestDrawCardFromDeckToPlayerPlayerNotExists(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	expected := core.Player{
		Username: username,
		Cards:    sameCards,
		Stop:     false,
		Busted:   false,
	}
	gomock.InOrder(
		playerServiceMock.EXPECT().GetPlayer(ctx, telegramChatID, username).Return(nil, core.ErrNotFound),
		repo.EXPECT().DrawCardFromDeck(ctx, telegramChatID).Times(2).Return(new5Card, nil),
		playerServiceMock.EXPECT().AddNewPlayer(ctx, telegramChatID, expected).Return(nil),
	)
	service := NewCardService(logger, repo, playerServiceMock, 1)

	player, err := service.DrawCardFromDeckToPlayer(ctx, telegramChatID, username)

	require.NoError(t, err)
	require.Equal(t, expected, *player)
}

func TestDrawCardFromDeckToPlayerGameNotExists(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	var expected *core.Player
	gomock.InOrder(
		playerServiceMock.EXPECT().GetPlayer(ctx, telegramChatID, username).Return(nil, core.ErrNoActiveGame),
	)
	service := NewCardService(logger, repo, playerServiceMock, 1)

	player, err := service.DrawCardFromDeckToPlayer(ctx, telegramChatID, username)

	require.ErrorIs(t, err, core.ErrNoActiveGame)
	require.Equal(t, expected, player)
}

func TestDrawCardFromDeckToPlayerServerError(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	var expected *core.Player
	gomock.InOrder(
		playerServiceMock.EXPECT().GetPlayer(ctx, telegramChatID, username).Return(nil, core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewCardService(logger, repo, playerServiceMock, 1)

	player, err := service.DrawCardFromDeckToPlayer(ctx, telegramChatID, username)

	require.ErrorIs(t, err, core.ErrServerError)
	require.Equal(t, expected, player)
}

func TestDrawCardFromDeckToPlayerBustedCase(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	expected := &core.Player{
		Username: username,
		Cards:    append(cards, core.Cards{new5Card, newKCard}...),
		Stop:     true,
		Busted:   true,
	}
	player := &core.Player{
		Username: username,
		Cards:    append(cards, new5Card),
	}
	gomock.InOrder(
		playerServiceMock.EXPECT().GetPlayer(ctx, telegramChatID, username).Return(player, nil),
		repo.EXPECT().DrawCardFromDeck(ctx, telegramChatID).Return(newKCard, nil),
		repo.EXPECT().AddCardToPlayer(ctx, telegramChatID, username, newKCard).Return(nil),
		playerServiceMock.EXPECT().StopDrawing(ctx, telegramChatID, player).Return(nil),
	)
	service := NewCardService(logger, repo, playerServiceMock, 1)

	player, err := service.DrawCardFromDeckToPlayer(ctx, telegramChatID, username)

	require.ErrorIs(t, err, core.ErrBusted)
	require.Equal(t, expected, player)
}

func TestDrawCardFromDeckToPlayer21ValueCase(t *testing.T) {
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
		repo.EXPECT().DrawCardFromDeck(ctx, telegramChatID).Return(new9Card, nil),
		repo.EXPECT().AddCardToPlayer(ctx, telegramChatID, username, new9Card).Return(nil),
		playerServiceMock.EXPECT().StopDrawing(ctx, telegramChatID, player).Return(nil),
	)
	service := NewCardService(logger, repo, playerServiceMock, 1)

	player, err := service.DrawCardFromDeckToPlayer(ctx, telegramChatID, username)

	require.NoError(t, err)
	require.Equal(t, expected, player)
	require.Equal(t, 21, player.Cards.CountValue())
}

func TestDrawCardFromDeckToPlayerBlackjackCase(t *testing.T) {
	ctx, logger, repo, playerServiceMock := getCardServiceDependencies(t)
	expected := &core.Player{
		Username: username,
		Cards:    blackjack,
		Stop:     true,
		Busted:   false,
	}
	gomock.InOrder(
		playerServiceMock.EXPECT().GetPlayer(ctx, telegramChatID, username).Return(nil, core.ErrNotFound),
		repo.EXPECT().DrawCardFromDeck(ctx, telegramChatID).Return(newACard, nil),
		repo.EXPECT().DrawCardFromDeck(ctx, telegramChatID).Return(newKCard, nil),
		playerServiceMock.EXPECT().AddNewPlayer(ctx, telegramChatID, *expected).Return(nil),
	)
	service := NewCardService(logger, repo, playerServiceMock, 1)

	player, err := service.DrawCardFromDeckToPlayer(ctx, telegramChatID, username)

	require.NoError(t, err)
	require.Equal(t, expected, player)
	require.Equal(t, 21, player.Cards.CountValue())
	require.True(t, player.Cards.IsBlackjack())
}
