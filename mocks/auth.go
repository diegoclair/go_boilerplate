// Code generated by MockGen. DO NOT EDIT.
// Source: infra/auth/instance.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	auth "github.com/diegoclair/go_boilerplate/infra/auth"
	gomock "github.com/golang/mock/gomock"
)

// MockAuthToken is a mock of AuthToken interface.
type MockAuthToken struct {
	ctrl     *gomock.Controller
	recorder *MockAuthTokenMockRecorder
}

// MockAuthTokenMockRecorder is the mock recorder for MockAuthToken.
type MockAuthTokenMockRecorder struct {
	mock *MockAuthToken
}

// NewMockAuthToken creates a new mock instance.
func NewMockAuthToken(ctrl *gomock.Controller) *MockAuthToken {
	mock := &MockAuthToken{ctrl: ctrl}
	mock.recorder = &MockAuthTokenMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthToken) EXPECT() *MockAuthTokenMockRecorder {
	return m.recorder
}

// CreateAccessToken mocks base method.
func (m *MockAuthToken) CreateAccessToken(ctx context.Context, input auth.TokenPayloadInput) (string, *auth.TokenPayload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccessToken", ctx, input)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*auth.TokenPayload)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateAccessToken indicates an expected call of CreateAccessToken.
func (mr *MockAuthTokenMockRecorder) CreateAccessToken(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccessToken", reflect.TypeOf((*MockAuthToken)(nil).CreateAccessToken), ctx, input)
}

// CreateRefreshToken mocks base method.
func (m *MockAuthToken) CreateRefreshToken(ctx context.Context, input auth.TokenPayloadInput) (string, *auth.TokenPayload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRefreshToken", ctx, input)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*auth.TokenPayload)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateRefreshToken indicates an expected call of CreateRefreshToken.
func (mr *MockAuthTokenMockRecorder) CreateRefreshToken(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRefreshToken", reflect.TypeOf((*MockAuthToken)(nil).CreateRefreshToken), ctx, input)
}

// VerifyToken mocks base method.
func (m *MockAuthToken) VerifyToken(ctx context.Context, token string) (*auth.TokenPayload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyToken", ctx, token)
	ret0, _ := ret[0].(*auth.TokenPayload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyToken indicates an expected call of VerifyToken.
func (mr *MockAuthTokenMockRecorder) VerifyToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", reflect.TypeOf((*MockAuthToken)(nil).VerifyToken), ctx, token)
}
