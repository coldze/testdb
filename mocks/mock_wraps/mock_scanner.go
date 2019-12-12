// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/coldze/testdb/logic/sources/wraps (interfaces: Scanner)

// Package mock_wraps is a generated GoMock package.
package mock_wraps

import (
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockScanner is a mock of Scanner interface
type MockScanner struct {
	ctrl     *gomock.Controller
	recorder *MockScannerMockRecorder
}

// MockScannerMockRecorder is the mock recorder for MockScanner
type MockScannerMockRecorder struct {
	mock *MockScanner
}

// NewMockScanner creates a new mock instance
func NewMockScanner(ctrl *gomock.Controller) *MockScanner {
	mock := &MockScanner{ctrl: ctrl}
	mock.recorder = &MockScannerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockScanner) EXPECT() *MockScannerMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockScanner) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockScannerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockScanner)(nil).Close))
}

// Err mocks base method
func (m *MockScanner) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err
func (mr *MockScannerMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockScanner)(nil).Err))
}

// Next mocks base method
func (m *MockScanner) Next() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Next indicates an expected call of Next
func (mr *MockScannerMockRecorder) Next() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockScanner)(nil).Next))
}

// Scan mocks base method
func (m *MockScanner) Scan(arg0 ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Scan", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Scan indicates an expected call of Scan
func (mr *MockScannerMockRecorder) Scan(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockScanner)(nil).Scan), arg0...)
}
