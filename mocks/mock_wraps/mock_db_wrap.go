// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/coldze/testdb/logic/sources/wraps (interfaces: DbWrap)

// Package mock_wraps is a generated GoMock package.
package mock_wraps

import (
	"reflect"

	"github.com/golang/mock/gomock"

	"github.com/coldze/testdb/logic/sources/wraps"
)

// MockDbWrap is a mock of DbWrap interface
type MockDbWrap struct {
	ctrl     *gomock.Controller
	recorder *MockDbWrapMockRecorder
}

// MockDbWrapMockRecorder is the mock recorder for MockDbWrap
type MockDbWrapMockRecorder struct {
	mock *MockDbWrap
}

// NewMockDbWrap creates a new mock instance
func NewMockDbWrap(ctrl *gomock.Controller) *MockDbWrap {
	mock := &MockDbWrap{ctrl: ctrl}
	mock.recorder = &MockDbWrapMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDbWrap) EXPECT() *MockDbWrapMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockDbWrap) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockDbWrapMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDbWrap)(nil).Close))
}

// Query mocks base method
func (m *MockDbWrap) Query(arg0 string, arg1 ...interface{}) (wraps.Scanner, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(wraps.Scanner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query
func (mr *MockDbWrapMockRecorder) Query(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockDbWrap)(nil).Query), varargs...)
}
