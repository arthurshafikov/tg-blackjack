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
	game = core.Game{
		Dealer: core.Cards{new5Card, new5Card},
		Players: []core.Player{
			{
				Username: "some_stopped_user",
				Cards:    core.Cards{new9Card, new5Card, new5Card},
				Stop:     true,
				Busted:   false,
			},
			{
				Username: "some_busted_user",
				Cards:    core.Cards{new9Card, new9Card, new9Card},
				Stop:     true,
				Busted:   true,
			},
			{
				Username: "some_blackjack_user",
				Cards:    blackjack,
				Stop:     true,
				Busted:   false,
			},
			{
				Username: "some_21_user",
				Cards:    core.Cards{"♣10", "♥2", new9Card},
				Stop:     true,
				Busted:   false,
			},
		},
	}
	emptyGame = core.Game{
		Dealer:  core.Cards{},
		Players: []core.Player{},
	}
)

func getGameServiceDependencies(t *testing.T) (
	context.Context,
	*mock_services.MockLogger,
	*mock_repository.MockGames,
	*mock_services.MockStatistics,
	*mock_services.MockCards,
) {
	t.Helper()
	ctrl := gomock.NewController(t)

	return context.Background(),
		mock_services.NewMockLogger(ctrl),
		mock_repository.NewMockGames(ctrl),
		mock_services.NewMockStatistics(ctrl),
		mock_services.NewMockCards(ctrl)
}

func TestNewGame(t *testing.T) {
	ctx, logger, repo, stats, cards := getGameServiceDependencies(t)
	expected := emptyGame
	expected.Dealer = core.Cards{new9Card, new9Card}
	gomock.InOrder(
		repo.EXPECT().SetActiveGame(ctx, telegramChatID, emptyGame).Return(nil),
		cards.EXPECT().DrawCardFromDeckToDealer(ctx, telegramChatID).Times(2).Return(new9Card, nil),
	)
	service := NewGameService(logger, repo, stats, cards)

	result, err := service.NewGame(ctx, telegramChatID)
	require.NoError(t, err)
	require.Equal(t, expected, *result)
}

func TestNewGameActiveGameExists(t *testing.T) {
	ctx, logger, repo, stats, cards := getGameServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().SetActiveGame(ctx, telegramChatID, emptyGame).Return(core.ErrActiveGame),
	)
	service := NewGameService(logger, repo, stats, cards)

	result, err := service.NewGame(ctx, telegramChatID)
	require.ErrorIs(t, err, core.ErrActiveGame)
	require.Nil(t, result)
}

func TestNewGameServerError(t *testing.T) {
	ctx, logger, repo, stats, cards := getGameServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().SetActiveGame(ctx, telegramChatID, emptyGame).Return(core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewGameService(logger, repo, stats, cards)

	result, err := service.NewGame(ctx, telegramChatID)
	require.ErrorIs(t, err, core.ErrServerError)
	require.Nil(t, result)
}

func TestNewGameDrawToDealerCausedServerError(t *testing.T) {
	ctx, logger, repo, stats, cards := getGameServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().SetActiveGame(ctx, telegramChatID, emptyGame).Return(nil),
		cards.EXPECT().DrawCardFromDeckToDealer(ctx, telegramChatID).Return(new5Card, core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewGameService(logger, repo, stats, cards)

	result, err := service.NewGame(ctx, telegramChatID)
	require.ErrorIs(t, err, core.ErrServerError)
	require.Nil(t, result)
}

func TestCheckIfGameShouldBeFinished(t *testing.T) {
	ctx, logger, repo, stats, cards := getGameServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetActiveGame(ctx, telegramChatID).Return(&game, nil),
	)
	service := NewGameService(logger, repo, stats, cards)

	result, err := service.CheckIfGameShouldBeFinished(ctx, telegramChatID)
	require.NoError(t, err)
	require.True(t, result)
}

func TestCheckIfGameShouldBeFinishedActivePlayers(t *testing.T) {
	ctx, logger, repo, stats, cards := getGameServiceDependencies(t)
	game := game
	game.Players = append(game.Players, core.Player{
		Stop: false,
	})
	gomock.InOrder(
		repo.EXPECT().GetActiveGame(ctx, telegramChatID).Return(&game, nil),
	)
	service := NewGameService(logger, repo, stats, cards)

	result, err := service.CheckIfGameShouldBeFinished(ctx, telegramChatID)
	require.NoError(t, err)
	require.False(t, result)
}

func TestCheckIfGameShouldBeFinishedServerError(t *testing.T) {
	ctx, logger, repo, stats, cards := getGameServiceDependencies(t)
	gomock.InOrder(
		repo.EXPECT().GetActiveGame(ctx, telegramChatID).Return(&emptyGame, core.ErrServerError),
		logger.EXPECT().Error(core.ErrServerError),
	)
	service := NewGameService(logger, repo, stats, cards)

	result, err := service.CheckIfGameShouldBeFinished(ctx, telegramChatID)
	require.ErrorIs(t, err, core.ErrServerError)
	require.True(t, result)
}

func TestFinishGame(t *testing.T) {
	ctx, logger, repo, stats, cards := getGameServiceDependencies(t)
	expectedGame := game
	expectedGame.Dealer = append(expectedGame.Dealer, new9Card)
	expectedStats := core.UsersStatistics{
		"some_stopped_user":   0,
		"some_busted_user":    -1,
		"some_blackjack_user": 2,
		"some_21_user":        1,
	}
	gomock.InOrder(
		repo.EXPECT().GetActiveGame(ctx, telegramChatID).Return(&game, nil),
		cards.EXPECT().DrawCardFromDeckToDealer(ctx, telegramChatID).Return(new9Card, nil),
		repo.EXPECT().NullActiveGame(ctx, telegramChatID).Return(nil),
		stats.EXPECT().IncrementStatistic(ctx, telegramChatID, expectedStats).Return(nil),
	)
	service := NewGameService(logger, repo, stats, cards)

	resultGame, resultStats, err := service.FinishGame(ctx, telegramChatID)
	require.NoError(t, err)
	require.Equal(t, expectedGame, *resultGame)
	require.Equal(t, expectedStats, resultStats)
}

func TestFinishGameDealerHasBlackJack(t *testing.T) {
	ctx, logger, repo, stats, cards := getGameServiceDependencies(t)
	game := game
	game.Dealer = blackjack
	expectedGame := game
	expectedStats := core.UsersStatistics{
		"some_stopped_user":   -1,
		"some_busted_user":    -1,
		"some_blackjack_user": 0,
		"some_21_user":        -1,
	}
	gomock.InOrder(
		repo.EXPECT().GetActiveGame(ctx, telegramChatID).Return(&game, nil),
		repo.EXPECT().NullActiveGame(ctx, telegramChatID).Return(nil),
		stats.EXPECT().IncrementStatistic(ctx, telegramChatID, expectedStats).Return(nil),
	)
	service := NewGameService(logger, repo, stats, cards)

	resultGame, resultStats, err := service.FinishGame(ctx, telegramChatID)
	require.NoError(t, err)
	require.Equal(t, expectedGame, *resultGame)
	require.Equal(t, expectedStats, resultStats)
}

func TestFinishGameDealerBusted(t *testing.T) {
	ctx, logger, repo, stats, cards := getGameServiceDependencies(t)
	game := game
	game.Dealer = core.Cards{new9Card, new5Card}
	expectedGame := game
	expectedGame.Dealer = append(expectedGame.Dealer, newKCard)
	expectedStats := core.UsersStatistics{
		"some_stopped_user":   1,
		"some_busted_user":    -1,
		"some_blackjack_user": 2,
		"some_21_user":        1,
	}
	gomock.InOrder(
		repo.EXPECT().GetActiveGame(ctx, telegramChatID).Return(&game, nil),
		cards.EXPECT().DrawCardFromDeckToDealer(ctx, telegramChatID).Return(newKCard, nil),
		repo.EXPECT().NullActiveGame(ctx, telegramChatID).Return(nil),
		stats.EXPECT().IncrementStatistic(ctx, telegramChatID, expectedStats).Return(nil),
	)
	service := NewGameService(logger, repo, stats, cards)

	resultGame, resultStats, err := service.FinishGame(ctx, telegramChatID)
	require.NoError(t, err)
	require.Equal(t, expectedGame, *resultGame)
	require.Equal(t, expectedStats, resultStats)
}
