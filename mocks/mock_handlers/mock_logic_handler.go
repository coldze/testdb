package mock_handlers

import (
	"net/http"
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockLogicHandler is a mock of LogicHandler interface
type MockLogicHandler struct {
	ctrl     *gomock.Controller
	recorder *MockLogicHandlerMockRecorder
}

// MockLogicHandlerMockRecorder is the mock recorder for MockLogicHandler
type MockLogicHandlerMockRecorder struct {
	mock *MockLogicHandler
}

func NewMockLogicHandler(ctrl *gomock.Controller) *MockLogicHandler {
	mock := &MockLogicHandler{ctrl: ctrl}
	mock.recorder = &MockLogicHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogicHandler) EXPECT() *MockLogicHandlerMockRecorder {
	return m.recorder
}

// Handle mocks base method
func (m *MockLogicHandler) Handle(arg0 *http.Request) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Handle
func (mr *MockLogicHandlerMockRecorder) Handle(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockLogicHandler)(nil).Handle), arg0)
}
