package mock_handlers

import (
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockIDFactory is a mock of IDFactory
type MockIDFactory struct {
	ctrl     *gomock.Controller
	recorder *MockIDFactoryMockRecorder
}

// MockIDFactoryMockRecorder is the mock recorder for MockIDFactory
type MockIDFactoryMockRecorder struct {
	mock *MockIDFactory
}

func NewMockIDFactory(ctrl *gomock.Controller) *MockIDFactory {
	mock := &MockIDFactory{ctrl: ctrl}
	mock.recorder = &MockIDFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIDFactory) EXPECT() *MockIDFactoryMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockIDFactory) Create() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create")
	ret0, _ := ret[0].(string)
	return ret0
}

// Do indicates an expected call of Do
func (mr *MockIDFactoryMockRecorder) Create() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIDFactory)(nil).Create))
}
