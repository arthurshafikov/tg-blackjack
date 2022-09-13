// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/telegram/handlers/helper.go

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	reflect "reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gomock "github.com/golang/mock/gomock"
)

// MockTelegramHandlerHelper is a mock of TelegramHandlerHelper interface.
type MockTelegramHandlerHelper struct {
	ctrl     *gomock.Controller
	recorder *MockTelegramHandlerHelperMockRecorder
}

// MockTelegramHandlerHelperMockRecorder is the mock recorder for MockTelegramHandlerHelper.
type MockTelegramHandlerHelperMockRecorder struct {
	mock *MockTelegramHandlerHelper
}

// NewMockTelegramHandlerHelper creates a new mock instance.
func NewMockTelegramHandlerHelper(ctrl *gomock.Controller) *MockTelegramHandlerHelper {
	mock := &MockTelegramHandlerHelper{ctrl: ctrl}
	mock.recorder = &MockTelegramHandlerHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTelegramHandlerHelper) EXPECT() *MockTelegramHandlerHelperMockRecorder {
	return m.recorder
}

// GetUpdatesChan mocks base method.
func (m *MockTelegramHandlerHelper) GetUpdatesChan(config tgbotapi.UpdateConfig) (tgbotapi.UpdatesChannel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUpdatesChan", config)
	ret0, _ := ret[0].(tgbotapi.UpdatesChannel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUpdatesChan indicates an expected call of GetUpdatesChan.
func (mr *MockTelegramHandlerHelperMockRecorder) GetUpdatesChan(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpdatesChan", reflect.TypeOf((*MockTelegramHandlerHelper)(nil).GetUpdatesChan), config)
}

// NewMessage mocks base method.
func (m *MockTelegramHandlerHelper) NewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewMessage", chatID, text)
	ret0, _ := ret[0].(tgbotapi.MessageConfig)
	return ret0
}

// NewMessage indicates an expected call of NewMessage.
func (mr *MockTelegramHandlerHelperMockRecorder) NewMessage(chatID, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewMessage", reflect.TypeOf((*MockTelegramHandlerHelper)(nil).NewMessage), chatID, text)
}

// NewUpdateChannel mocks base method.
func (m *MockTelegramHandlerHelper) NewUpdateChannel(offset int) tgbotapi.UpdateConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewUpdateChannel", offset)
	ret0, _ := ret[0].(tgbotapi.UpdateConfig)
	return ret0
}

// NewUpdateChannel indicates an expected call of NewUpdateChannel.
func (mr *MockTelegramHandlerHelperMockRecorder) NewUpdateChannel(offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewUpdateChannel", reflect.TypeOf((*MockTelegramHandlerHelper)(nil).NewUpdateChannel), offset)
}

// SendMessage mocks base method.
func (m *MockTelegramHandlerHelper) SendMessage(msg tgbotapi.MessageConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockTelegramHandlerHelperMockRecorder) SendMessage(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockTelegramHandlerHelper)(nil).SendMessage), msg)
}