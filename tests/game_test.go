package tests

import (
	"errors"

	"github.com/arthurshafikov/tg-blackjack/internal/core"
)

func (s *APITestSuite) TestGameOnePlayer() {
	err := s.services.Chats.RegisterChat(s.ctx, telegramChatID)
	r.NoError(err)

	_, err = s.services.Games.NewGame(s.ctx, telegramChatID)
	r.NoError(err)

	player := s.playerEntersTheGame(username)

	if !player.Stop {
		res, err := s.services.Games.CheckIfGameShouldBeFinished(s.ctx, telegramChatID)
		r.NoError(err)
		r.False(res)
		r.NoError(s.services.Players.StopDrawing(s.ctx, telegramChatID, player))
	}

	res, err := s.services.Games.CheckIfGameShouldBeFinished(s.ctx, telegramChatID)
	r.NoError(err)
	r.True(res)

	game, gameStats, err := s.services.Games.FinishGame(s.ctx, telegramChatID)
	r.NoError(err)
	r.Len(game.Players, 1)

	playerStat, ok := gameStats[username]
	r.True(ok)
	r.True(-1 <= playerStat && playerStat <= 2)

	globalStats, err := s.services.Statistics.GetStatistics(s.ctx, telegramChatID)
	s.NoError(err)
	r.Equal(playerStat, globalStats[username])
}

func (s *APITestSuite) TestGameThreePlayers() {
	err := s.services.Chats.RegisterChat(s.ctx, telegramChatID)
	r.NoError(err)

	_, err = s.services.Games.NewGame(s.ctx, telegramChatID)
	r.NoError(err)

	playersSlice := []*core.Player{
		s.playerEntersTheGame(playersUsernames[0]),
		s.playerEntersTheGame(playersUsernames[1]),
		s.playerEntersTheGame(playersUsernames[2]),
	}

	for i, player := range playersSlice {
		for !player.Stop {
			player, err = s.services.Cards.DrawCardFromDeckToPlayer(s.ctx, telegramChatID, player.Username)
			if errors.Is(err, core.ErrBusted) {
				break
			}
			r.NoError(err)
		}
		playersSlice[i] = player
	}

	for _, player := range playersSlice {
		r.True(player.Stop)
		if player.Cards.CountValue() > 21 {
			r.True(player.Busted)
		}
	}

	res, err := s.services.Games.CheckIfGameShouldBeFinished(s.ctx, telegramChatID)
	r.NoError(err)
	r.True(res)

	game, gameStats, err := s.services.Games.FinishGame(s.ctx, telegramChatID)
	r.NoError(err)
	r.Len(game.Players, 3)

	globalStats, err := s.services.Statistics.GetStatistics(s.ctx, telegramChatID)
	s.NoError(err)

	for _, player := range playersSlice {
		playerStat, ok := gameStats[player.Username]
		r.True(ok)
		r.True(-1 <= playerStat && playerStat <= 2)
		r.Equal(playerStat, globalStats[player.Username])
	}
}

func (s *APITestSuite) playerEntersTheGame(username string) *core.Player {
	player, err := s.services.Cards.DrawCardFromDeckToPlayer(s.ctx, telegramChatID, username)
	r.NoError(err)
	r.Len(player.Cards, 2)
	r.Equal(player.Username, username)
	r.True(player.Cards.IsBlackjack() || !player.Stop)
	r.False(player.Busted)

	return player
}
