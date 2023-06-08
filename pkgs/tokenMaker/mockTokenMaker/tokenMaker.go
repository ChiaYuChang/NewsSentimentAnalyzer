// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker (interfaces: TokenMaker,Payload)

// Package mock_tokenMaker is a generated GoMock package.
package mock_tokenMaker

import (
	reflect "reflect"

	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	gomock "github.com/golang/mock/gomock"
)

// MockTokenMaker is a mock of TokenMaker interface.
type MockTokenMaker struct {
	ctrl     *gomock.Controller
	recorder *MockTokenMakerMockRecorder
}

// MockTokenMakerMockRecorder is the mock recorder for MockTokenMaker.
type MockTokenMakerMockRecorder struct {
	mock *MockTokenMaker
}

// NewMockTokenMaker creates a new mock instance.
func NewMockTokenMaker(ctrl *gomock.Controller) *MockTokenMaker {
	mock := &MockTokenMaker{ctrl: ctrl}
	mock.recorder = &MockTokenMakerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenMaker) EXPECT() *MockTokenMakerMockRecorder {
	return m.recorder
}

// MakeToken mocks base method.
func (m *MockTokenMaker) MakeToken(arg0 string, arg1 int32, arg2 tokenmaker.Role) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeToken", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MakeToken indicates an expected call of MakeToken.
func (mr *MockTokenMakerMockRecorder) MakeToken(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeToken", reflect.TypeOf((*MockTokenMaker)(nil).MakeToken), arg0, arg1, arg2)
}

// String mocks base method.
func (m *MockTokenMaker) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockTokenMakerMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockTokenMaker)(nil).String))
}

// ValidateToken mocks base method.
func (m *MockTokenMaker) ValidateToken(arg0 string) (tokenmaker.Payload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", arg0)
	ret0, _ := ret[0].(tokenmaker.Payload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockTokenMakerMockRecorder) ValidateToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockTokenMaker)(nil).ValidateToken), arg0)
}

// MockPayload is a mock of Payload interface.
type MockPayload struct {
	ctrl     *gomock.Controller
	recorder *MockPayloadMockRecorder
}

// MockPayloadMockRecorder is the mock recorder for MockPayload.
type MockPayloadMockRecorder struct {
	mock *MockPayload
}

// NewMockPayload creates a new mock instance.
func NewMockPayload(ctrl *gomock.Controller) *MockPayload {
	mock := &MockPayload{ctrl: ctrl}
	mock.recorder = &MockPayloadMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPayload) EXPECT() *MockPayloadMockRecorder {
	return m.recorder
}

// GetRole mocks base method.
func (m *MockPayload) GetRole() tokenmaker.Role {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole")
	ret0, _ := ret[0].(tokenmaker.Role)
	return ret0
}

// GetRole indicates an expected call of GetRole.
func (mr *MockPayloadMockRecorder) GetRole() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockPayload)(nil).GetRole))
}

// GetUserID mocks base method.
func (m *MockPayload) GetUserID() int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserID")
	ret0, _ := ret[0].(int32)
	return ret0
}

// GetUserID indicates an expected call of GetUserID.
func (mr *MockPayloadMockRecorder) GetUserID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserID", reflect.TypeOf((*MockPayload)(nil).GetUserID))
}

// GetUserInfo mocks base method.
func (m *MockPayload) GetUserInfo() tokenmaker.UserInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfo")
	ret0, _ := ret[0].(tokenmaker.UserInfo)
	return ret0
}

// GetUserInfo indicates an expected call of GetUserInfo.
func (mr *MockPayloadMockRecorder) GetUserInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfo", reflect.TypeOf((*MockPayload)(nil).GetUserInfo))
}

// GetUsername mocks base method.
func (m *MockPayload) GetUsername() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsername")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetUsername indicates an expected call of GetUsername.
func (mr *MockPayloadMockRecorder) GetUsername() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsername", reflect.TypeOf((*MockPayload)(nil).GetUsername))
}

// String mocks base method.
func (m *MockPayload) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockPayloadMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockPayload)(nil).String))
}
