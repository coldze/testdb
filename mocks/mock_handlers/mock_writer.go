// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/coldze/testdb/logic/handlers (interfaces: Writer)

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockWriter is a mock of Writer interface
type MockWriter struct {
	ctrl     *gomock.Controller
	recorder *MockWriterMockRecorder
}

// MockWriterMockRecorder is the mock recorder for MockWriter
type MockWriterMockRecorder struct {
	mock *MockWriter
}

// NewMockWriter creates a new mock instance
func NewMockWriter(ctrl *gomock.Controller) *MockWriter {
	mock := &MockWriter{ctrl: ctrl}
	mock.recorder = &MockWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWriter) EXPECT() *MockWriterMockRecorder {
	return m.recorder
}

// Data mocks base method
func (m *MockWriter) Data(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Data", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Data indicates an expected call of Data
func (mr *MockWriterMockRecorder) Data(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Data", reflect.TypeOf((*MockWriter)(nil).Data), arg0)
}

// Error mocks base method
func (m *MockWriter) Error(arg0 error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Error", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Error indicates an expected call of Error
func (mr *MockWriterMockRecorder) Error(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockWriter)(nil).Error), arg0)
}
