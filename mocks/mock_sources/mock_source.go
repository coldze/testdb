// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/coldze/testdb/logic/sources (interfaces: Source)

// Package mock_sources is a generated GoMock package.
package mock_sources

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"

	"github.com/coldze/testdb/logic/structs"
)

// MockSource is a mock of Source interface
type MockSource struct {
	ctrl     *gomock.Controller
	recorder *MockSourceMockRecorder
}

// MockSourceMockRecorder is the mock recorder for MockSource
type MockSourceMockRecorder struct {
	mock *MockSource
}

// NewMockSource creates a new mock instance
func NewMockSource(ctrl *gomock.Controller) *MockSource {
	mock := &MockSource{ctrl: ctrl}
	mock.recorder = &MockSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSource) EXPECT() *MockSourceMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockSource) Get(arg0 context.Context, arg1 structs.Request) (structs.Data, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(structs.Data)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockSourceMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSource)(nil).Get), arg0, arg1)
}
