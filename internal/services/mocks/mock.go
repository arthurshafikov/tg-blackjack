// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/services/service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	reflect "reflect"

	core "github.com/arthurshafikov/tg-blackjack/internal/core"
	gomock "github.com/golang/mock/gomock"
)

// MockCards is a mock of Cards interface.
type MockCards struct {
	ctrl     *gomock.Controller
	recorder *MockCardsMockRecorder
}

// MockCardsMockRecorder is the mock recorder for MockCards.
type MockCardsMockRecorder struct {
	mock *MockCards
}

// NewMockCards creates a new mock instance.
func NewMockCards(ctrl *gomock.Controller) *MockCards {
	mock := &MockCards{ctrl: ctrl}
	mock.recorder = &MockCardsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCards) EXPECT() *MockCardsMockRecorder {
	return m.recorder
}

// DrawCardFromDeckToDealer mocks base method.
func (m *MockCards) DrawCardFromDeckToDealer(ctx context.Context, telegramChatID int64) (core.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DrawCardFromDeckToDealer", ctx, telegramChatID)
	ret0, _ := ret[0].(core.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DrawCardFromDeckToDealer indicates an expected call of DrawCardFromDeckToDealer.
func (mr *MockCardsMockRecorder) DrawCardFromDeckToDealer(ctx, telegramChatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DrawCardFromDeckToDealer", reflect.TypeOf((*MockCards)(nil).DrawCardFromDeckToDealer), ctx, telegramChatID)
}

// DrawCardFromDeckToPlayer mocks base method.
func (m *MockCards) DrawCardFromDeckToPlayer(ctx context.Context, telegramChatID int64, username string) (*core.Player, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DrawCardFromDeckToPlayer", ctx, telegramChatID, username)
	ret0, _ := ret[0].(*core.Player)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DrawCardFromDeckToPlayer indicates an expected call of DrawCardFromDeckToPlayer.
func (mr *MockCardsMockRecorder) DrawCardFromDeckToPlayer(ctx, telegramChatID, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DrawCardFromDeckToPlayer", reflect.TypeOf((*MockCards)(nil).DrawCardFromDeckToPlayer), ctx, telegramChatID, username)
}

// MockChats is a mock of Chats interface.
type MockChats struct {
	ctrl     *gomock.Controller
	recorder *MockChatsMockRecorder
}

// MockChatsMockRecorder is the mock recorder for MockChats.
type MockChatsMockRecorder struct {
	mock *MockChats
}

// NewMockChats creates a new mock instance.
func NewMockChats(ctrl *gomock.Controller) *MockChats {
	mock := &MockChats{ctrl: ctrl}
	mock.recorder = &MockChatsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChats) EXPECT() *MockChatsMockRecorder {
	return m.recorder
}

// CheckChatExists mocks base method.
func (m *MockChats) CheckChatExists(ctx context.Context, telegramChatID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckChatExists", ctx, telegramChatID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckChatExists indicates an expected call of CheckChatExists.
func (mr *MockChatsMockRecorder) CheckChatExists(ctx, telegramChatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckChatExists", reflect.TypeOf((*MockChats)(nil).CheckChatExists), ctx, telegramChatID)
}

// RegisterChat mocks base method.
func (m *MockChats) RegisterChat(ctx context.Context, telegramChatID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterChat", ctx, telegramChatID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterChat indicates an expected call of RegisterChat.
func (mr *MockChatsMockRecorder) RegisterChat(ctx, telegramChatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterChat", reflect.TypeOf((*MockChats)(nil).RegisterChat), ctx, telegramChatID)
}

// MockGames is a mock of Games interface.
type MockGames struct {
	ctrl     *gomock.Controller
	recorder *MockGamesMockRecorder
}

// MockGamesMockRecorder is the mock recorder for MockGames.
type MockGamesMockRecorder struct {
	mock *MockGames
}

// NewMockGames creates a new mock instance.
func NewMockGames(ctrl *gomock.Controller) *MockGames {
	mock := &MockGames{ctrl: ctrl}
	mock.recorder = &MockGamesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGames) EXPECT() *MockGamesMockRecorder {
	return m.recorder
}

// CheckIfGameShouldBeFinished mocks base method.
func (m *MockGames) CheckIfGameShouldBeFinished(ctx context.Context, telegramChatID int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIfGameShouldBeFinished", ctx, telegramChatID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIfGameShouldBeFinished indicates an expected call of CheckIfGameShouldBeFinished.
func (mr *MockGamesMockRecorder) CheckIfGameShouldBeFinished(ctx, telegramChatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIfGameShouldBeFinished", reflect.TypeOf((*MockGames)(nil).CheckIfGameShouldBeFinished), ctx, telegramChatID)
}

// FinishGame mocks base method.
func (m *MockGames) FinishGame(ctx context.Context, telegramChatID int64) (*core.Game, core.UsersStatistics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinishGame", ctx, telegramChatID)
	ret0, _ := ret[0].(*core.Game)
	ret1, _ := ret[1].(core.UsersStatistics)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FinishGame indicates an expected call of FinishGame.
func (mr *MockGamesMockRecorder) FinishGame(ctx, telegramChatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishGame", reflect.TypeOf((*MockGames)(nil).FinishGame), ctx, telegramChatID)
}

// NewGame mocks base method.
func (m *MockGames) NewGame(ctx context.Context, telegramChatID int64) (*core.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewGame", ctx, telegramChatID)
	ret0, _ := ret[0].(*core.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewGame indicates an expected call of NewGame.
func (mr *MockGamesMockRecorder) NewGame(ctx, telegramChatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewGame", reflect.TypeOf((*MockGames)(nil).NewGame), ctx, telegramChatID)
}

// MockPlayers is a mock of Players interface.
type MockPlayers struct {
	ctrl     *gomock.Controller
	recorder *MockPlayersMockRecorder
}

// MockPlayersMockRecorder is the mock recorder for MockPlayers.
type MockPlayersMockRecorder struct {
	mock *MockPlayers
}

// NewMockPlayers creates a new mock instance.
func NewMockPlayers(ctrl *gomock.Controller) *MockPlayers {
	mock := &MockPlayers{ctrl: ctrl}
	mock.recorder = &MockPlayersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlayers) EXPECT() *MockPlayersMockRecorder {
	return m.recorder
}

// AddNewPlayer mocks base method.
func (m *MockPlayers) AddNewPlayer(ctx context.Context, telegramChatID int64, player core.Player) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewPlayer", ctx, telegramChatID, player)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewPlayer indicates an expected call of AddNewPlayer.
func (mr *MockPlayersMockRecorder) AddNewPlayer(ctx, telegramChatID, player interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewPlayer", reflect.TypeOf((*MockPlayers)(nil).AddNewPlayer), ctx, telegramChatID, player)
}

// GetPlayer mocks base method.
func (m *MockPlayers) GetPlayer(ctx context.Context, telegramChatID int64, username string) (*core.Player, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlayer", ctx, telegramChatID, username)
	ret0, _ := ret[0].(*core.Player)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlayer indicates an expected call of GetPlayer.
func (mr *MockPlayersMockRecorder) GetPlayer(ctx, telegramChatID, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlayer", reflect.TypeOf((*MockPlayers)(nil).GetPlayer), ctx, telegramChatID, username)
}

// StopDrawing mocks base method.
func (m *MockPlayers) StopDrawing(ctx context.Context, telegramChatID int64, player *core.Player) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopDrawing", ctx, telegramChatID, player)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopDrawing indicates an expected call of StopDrawing.
func (mr *MockPlayersMockRecorder) StopDrawing(ctx, telegramChatID, player interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopDrawing", reflect.TypeOf((*MockPlayers)(nil).StopDrawing), ctx, telegramChatID, player)
}

// MockStatistics is a mock of Statistics interface.
type MockStatistics struct {
	ctrl     *gomock.Controller
	recorder *MockStatisticsMockRecorder
}

// MockStatisticsMockRecorder is the mock recorder for MockStatistics.
type MockStatisticsMockRecorder struct {
	mock *MockStatistics
}

// NewMockStatistics creates a new mock instance.
func NewMockStatistics(ctrl *gomock.Controller) *MockStatistics {
	mock := &MockStatistics{ctrl: ctrl}
	mock.recorder = &MockStatisticsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatistics) EXPECT() *MockStatisticsMockRecorder {
	return m.recorder
}

// GetStatistics mocks base method.
func (m *MockStatistics) GetStatistics(ctx context.Context, telegramChatID int64) (core.UsersStatistics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatistics", ctx, telegramChatID)
	ret0, _ := ret[0].(core.UsersStatistics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatistics indicates an expected call of GetStatistics.
func (mr *MockStatisticsMockRecorder) GetStatistics(ctx, telegramChatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatistics", reflect.TypeOf((*MockStatistics)(nil).GetStatistics), ctx, telegramChatID)
}

// IncrementStatistic mocks base method.
func (m *MockStatistics) IncrementStatistic(ctx context.Context, telegramChatID int64, gameResult core.UsersStatistics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementStatistic", ctx, telegramChatID, gameResult)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementStatistic indicates an expected call of IncrementStatistic.
func (mr *MockStatisticsMockRecorder) IncrementStatistic(ctx, telegramChatID, gameResult interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementStatistic", reflect.TypeOf((*MockStatistics)(nil).IncrementStatistic), ctx, telegramChatID, gameResult)
}

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *MockLogger) Error(err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", err)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerMockRecorder) Error(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), err)
}
