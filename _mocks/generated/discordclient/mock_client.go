// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aplombomb/boombot/discord/ifaces (interfaces: DisgordClientAPI)

// Package mock_disgordclient is a generated GoMock package.
package mock_disgordclient

import (
	disgord "github.com/andersfylling/disgord"
	snowflake "github.com/andersfylling/snowflake/v4"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDisgordClientAPI is a mock of DisgordClientAPI interface
type MockDisgordClientAPI struct {
	ctrl     *gomock.Controller
	recorder *MockDisgordClientAPIMockRecorder
}

// MockDisgordClientAPIMockRecorder is the mock recorder for MockDisgordClientAPI
type MockDisgordClientAPIMockRecorder struct {
	mock *MockDisgordClientAPI
}

// NewMockDisgordClientAPI creates a new mock instance
func NewMockDisgordClientAPI(ctrl *gomock.Controller) *MockDisgordClientAPI {
	mock := &MockDisgordClientAPI{ctrl: ctrl}
	mock.recorder = &MockDisgordClientAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDisgordClientAPI) EXPECT() *MockDisgordClientAPIMockRecorder {
	return m.recorder
}

// Channel mocks base method
func (m *MockDisgordClientAPI) Channel(arg0 snowflake.Snowflake) disgord.ChannelQueryBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Channel", arg0)
	ret0, _ := ret[0].(disgord.ChannelQueryBuilder)
	return ret0
}

// Channel indicates an expected call of Channel
func (mr *MockDisgordClientAPIMockRecorder) Channel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Channel", reflect.TypeOf((*MockDisgordClientAPI)(nil).Channel), arg0)
}

// SendMsg mocks base method
func (m *MockDisgordClientAPI) SendMsg(arg0 snowflake.Snowflake, arg1 ...interface{}) (*disgord.Message, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendMsg", varargs...)
	ret0, _ := ret[0].(*disgord.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockDisgordClientAPIMockRecorder) SendMsg(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDisgordClientAPI)(nil).SendMsg), varargs...)
}

// VoiceConnectOptions mocks base method
func (m *MockDisgordClientAPI) VoiceConnectOptions(arg0, arg1 snowflake.Snowflake, arg2, arg3 bool) (disgord.VoiceConnection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VoiceConnectOptions", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(disgord.VoiceConnection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VoiceConnectOptions indicates an expected call of VoiceConnectOptions
func (mr *MockDisgordClientAPIMockRecorder) VoiceConnectOptions(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VoiceConnectOptions", reflect.TypeOf((*MockDisgordClientAPI)(nil).VoiceConnectOptions), arg0, arg1, arg2, arg3)
}
