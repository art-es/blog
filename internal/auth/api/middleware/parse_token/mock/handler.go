// Code generated by MockGen. DO NOT EDIT.
// Source: endpoint.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	dto "github.com/art-es/blog/internal/auth/dto"
	gomock "github.com/golang/mock/gomock"
)

// Mockusecase is a mock of usecase interface.
type Mockusecase struct {
	ctrl     *gomock.Controller
	recorder *MockusecaseMockRecorder
}

// MockusecaseMockRecorder is the mock recorder for Mockusecase.
type MockusecaseMockRecorder struct {
	mock *Mockusecase
}

// NewMockusecase creates a new mock instance.
func NewMockusecase(ctrl *gomock.Controller) *Mockusecase {
	mock := &Mockusecase{ctrl: ctrl}
	mock.recorder = &MockusecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockusecase) EXPECT() *MockusecaseMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *Mockusecase) Use(ctx context.Context, in *dto.AccessTokenParseIn) (*dto.ParseTokenOut, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Use", ctx, in)
	ret0, _ := ret[0].(*dto.ParseTokenOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockusecaseMockRecorder) Do(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Use", reflect.TypeOf((*Mockusecase)(nil).Use), ctx, in)
}
