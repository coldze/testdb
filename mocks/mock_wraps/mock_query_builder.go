package mock_wraps

import (
	"reflect"

	"github.com/golang/mock/gomock"

	"github.com/coldze/testdb/logic/sources/wraps"
	"github.com/coldze/testdb/logic/structs"
)

// MockQueryBuilder is a mock of IDFactory
type MockQueryBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockQueryBuilderMockRecorder
}

// MockQueryBuilderMockRecorder is the mock recorder for MockQueryBuilder
type MockQueryBuilderMockRecorder struct {
	mock *MockQueryBuilder
}

func NewMockQueryBuilder(ctrl *gomock.Controller) *MockQueryBuilder {
	mock := &MockQueryBuilder{ctrl: ctrl}
	mock.recorder = &MockQueryBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockQueryBuilder) EXPECT() *MockQueryBuilderMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockQueryBuilder) Build(arg0 structs.Request) (wraps.Query, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Build", arg0)
	ret0, _ := ret[0].(wraps.Query)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do
func (mr *MockQueryBuilderMockRecorder) Execute(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockQueryBuilder)(nil).Build), arg0)
}

