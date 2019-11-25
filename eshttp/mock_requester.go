package eshttp

import (
	"net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRequester is a mock of Requester interface
type MockRequester struct {
	ctrl     *gomock.Controller
	recorder *MockRequesterMockRecorder
}

// MockRequesterMockRecorder is the mock recorder for MockRequester
type MockRequesterMockRecorder struct {
	mock *MockRequester
}

// NewMockRequester creates a new mock instance
func NewMockRequester(ctrl *gomock.Controller) *MockRequester {
	mock := &MockRequester{ctrl: ctrl}
	mock.recorder = &MockRequesterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRequester) EXPECT() *MockRequesterMockRecorder {
	return m.recorder
}

// PostBulkToElasticSearch mocks base method
func (m *MockRequester) PostBulkToElasticSearch(arg0 []byte, arg1 string, arg2 string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostBulkToElasticSearch", arg0, arg1, arg2)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostBulkToElasticSearch indicates an expected call of PostBulkToElasticSearch
func (mr *MockRequesterMockRecorder) PostBulkToElasticSearch(arg0 interface{}, arg1 interface{}, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostBulkToElasticSearch", reflect.TypeOf((*MockRequester)(nil).PostBulkToElasticSearch), arg0, arg1, arg2)
}
