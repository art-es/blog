// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	dto "github.com/art-es/blog/internal/auth/dto"
	gomock "github.com/golang/mock/gomock"
)

// Mockdatabus is a mock of databus interface.
type Mockdatabus struct {
	ctrl     *gomock.Controller
	recorder *MockdatabusMockRecorder
}

// MockdatabusMockRecorder is the mock recorder for Mockdatabus.
type MockdatabusMockRecorder struct {
	mock *Mockdatabus
}

// NewMockdatabus creates a new mock instance.
func NewMockdatabus(ctrl *gomock.Controller) *Mockdatabus {
	mock := &Mockdatabus{ctrl: ctrl}
	mock.recorder = &MockdatabusMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockdatabus) EXPECT() *MockdatabusMockRecorder {
	return m.recorder
}

// ProduceActivationEmail mocks base method.
func (m *Mockdatabus) ProduceActivationEmail(ctx context.Context, msg *dto.UserActivationEmailMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProduceActivationEmail", ctx, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProduceActivationEmail indicates an expected call of ProduceActivationEmail.
func (mr *MockdatabusMockRecorder) ProduceActivationEmail(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProduceActivationEmail", reflect.TypeOf((*Mockdatabus)(nil).ProduceActivationEmail), ctx, msg)
}
