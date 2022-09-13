// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/telegram/telegram.go

// Package mock_telegram is a generated GoMock package.
package mock_telegram

import (
	reflect "reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gomock "github.com/golang/mock/gomock"
)

// MockCommandHandler is a mock of CommandHandler interface.
type MockCommandHandler struct {
	ctrl     *gomock.Controller
	recorder *MockCommandHandlerMockRecorder
}

// MockCommandHandlerMockRecorder is the mock recorder for MockCommandHandler.
type MockCommandHandlerMockRecorder struct {
	mock *MockCommandHandler
}

// NewMockCommandHandler creates a new mock instance.
func NewMockCommandHandler(ctrl *gomock.Controller) *MockCommandHandler {
	mock := &MockCommandHandler{ctrl: ctrl}
	mock.recorder = &MockCommandHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandHandler) EXPECT() *MockCommandHandlerMockRecorder {
	return m.recorder
}

// HandleDrawCard mocks base method.
func (m *MockCommandHandler) HandleDrawCard(message *tgbotapi.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleDrawCard", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleDrawCard indicates an expected call of HandleDrawCard.
func (mr *MockCommandHandlerMockRecorder) HandleDrawCard(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleDrawCard", reflect.TypeOf((*MockCommandHandler)(nil).HandleDrawCard), message)
}

// HandleNewGame mocks base method.
func (m *MockCommandHandler) HandleNewGame(message *tgbotapi.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleNewGame", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleNewGame indicates an expected call of HandleNewGame.
func (mr *MockCommandHandlerMockRecorder) HandleNewGame(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleNewGame", reflect.TypeOf((*MockCommandHandler)(nil).HandleNewGame), message)
}

// HandleStart mocks base method.
func (m *MockCommandHandler) HandleStart(message *tgbotapi.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleStart", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleStart indicates an expected call of HandleStart.
func (mr *MockCommandHandlerMockRecorder) HandleStart(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleStart", reflect.TypeOf((*MockCommandHandler)(nil).HandleStart), message)
}

// HandleStats mocks base method.
func (m *MockCommandHandler) HandleStats(message *tgbotapi.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleStats", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleStats indicates an expected call of HandleStats.
func (mr *MockCommandHandlerMockRecorder) HandleStats(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleStats", reflect.TypeOf((*MockCommandHandler)(nil).HandleStats), message)
}

// HandleStopDrawing mocks base method.
func (m *MockCommandHandler) HandleStopDrawing(message *tgbotapi.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleStopDrawing", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleStopDrawing indicates an expected call of HandleStopDrawing.
func (mr *MockCommandHandlerMockRecorder) HandleStopDrawing(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleStopDrawing", reflect.TypeOf((*MockCommandHandler)(nil).HandleStopDrawing), message)
}