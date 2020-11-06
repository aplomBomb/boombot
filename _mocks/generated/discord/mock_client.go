// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aplombomb/boombot/discord/ifaces (interfaces: DisgordClientAPI)

// Package mock_disgordclient is a generated GoMock package.
package mock_disgordclient

import (
	context "context"
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

// DeleteMessage mocks base method
func (m *MockDisgordClientAPI) DeleteMessage(arg0 context.Context, arg1, arg2 snowflake.Snowflake, arg3 ...disgord.Flag) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteMessage", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMessage indicates an expected call of DeleteMessage
func (mr *MockDisgordClientAPIMockRecorder) DeleteMessage(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessage", reflect.TypeOf((*MockDisgordClientAPI)(nil).DeleteMessage), varargs...)
}

// GetMessage mocks base method
func (m *MockDisgordClientAPI) GetMessage(arg0 context.Context, arg1, arg2 snowflake.Snowflake, arg3 ...disgord.Flag) (*disgord.Message, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMessage", varargs...)
	ret0, _ := ret[0].(*disgord.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessage indicates an expected call of GetMessage
func (mr *MockDisgordClientAPIMockRecorder) GetMessage(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessage", reflect.TypeOf((*MockDisgordClientAPI)(nil).GetMessage), varargs...)
}

// SendMsg mocks base method
func (m *MockDisgordClientAPI) SendMsg(arg0 context.Context, arg1 snowflake.Snowflake, arg2 ...interface{}) (*disgord.Message, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendMsg", varargs...)
	ret0, _ := ret[0].(*disgord.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockDisgordClientAPIMockRecorder) SendMsg(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDisgordClientAPI)(nil).SendMsg), varargs...)
}
