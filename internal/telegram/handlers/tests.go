package handlers

import (
	"context"
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	mock_services "github.com/arthurshafikov/tg-blackjack/internal/services/mocks"
	mock_handlers "github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers/mocks"
	"github.com/golang/mock/gomock"
)

var messages = config.Messages{
	ChatAlreadyRegistered:   "ChatAlreadyRegistered",
	ChatNotExists:           "ChatNotExists",
	ChatCreatedSuccessfully: "ChatCreatedSuccessfully",
	ChatHasActiveGame:       "ChatHasActiveGame",
	ChatHasNoActiveGame:     "ChatHasNoActiveGame",
	Blackjack:               "Blackjack",
	GameOver:                "GameOver",
	Win:                     "Win",
	Lose:                    "Lose",
	Push:                    "Push",
	BlackjackResult:         "BlackjackResult",
	PlayerCantDraw:          "PlayerCantDraw",
	PlayerHand:              "PlayerHand %s",
	PlayerHandBusted:        "PlayerHandBusted",
	PlayerAlreadyStopped:    "PlayerAlreadyStopped %s",
	PlayerAlreadyBusted:     "PlayerAlreadyBusted %s",
	StoppedDrawing:          "StoppedDrawing %s",
	DealerHand:              "DealerHand",
	DealerBlackjack:         "DealerBlackjack",
	GameStartHint:           "GameStartHint",
	GameEnterHint:           "GameEnterHint",
	TopPlayers:              "TopPlayers",
}

type MockBag struct {
	Cards      *mock_services.MockCards
	Chats      *mock_services.MockChats
	Games      *mock_services.MockGames
	Players    *mock_services.MockPlayers
	Statistics *mock_services.MockStatistics
	Logger     *mock_services.MockLogger
	Helper     *mock_handlers.MockTelegramHandlerHelper
}

func GetHandlerParamsWithMocks(t *testing.T) (HandlerParams, *MockBag) {
	t.Helper()
	ctrl := gomock.NewController(t)
	cardsMock := mock_services.NewMockCards(ctrl)
	chatsMock := mock_services.NewMockChats(ctrl)
	gamesMock := mock_services.NewMockGames(ctrl)
	playersMock := mock_services.NewMockPlayers(ctrl)
	statisticsMock := mock_services.NewMockStatistics(ctrl)
	helperMock := mock_handlers.NewMockTelegramHandlerHelper(ctrl)
	loggerMock := mock_services.NewMockLogger(ctrl)
	services := &services.Services{
		Cards:      cardsMock,
		Chats:      chatsMock,
		Games:      gamesMock,
		Players:    playersMock,
		Statistics: statisticsMock,
	}

	return HandlerParams{
			Ctx:      context.Background(),
			Services: services,
			Logger:   loggerMock,
			Messages: messages,
			Helper:   helperMock,
		}, &MockBag{
			Cards:      cardsMock,
			Chats:      chatsMock,
			Games:      gamesMock,
			Players:    playersMock,
			Statistics: statisticsMock,
			Logger:     loggerMock,
			Helper:     helperMock,
		}
}
