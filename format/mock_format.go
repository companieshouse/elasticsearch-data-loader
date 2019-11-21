// Package format is a generated GoMock package.
package format

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockFormatter is a mock of Formatter interface
type MockFormatter struct {
	ctrl     *gomock.Controller
	recorder *MockFormatterMockRecorder
}

// MockFormatterMockRecorder is the mock recorder for MockFormatter
type MockFormatterMockRecorder struct {
	mock *MockFormatter
}

// NewMockFormatter creates a new mock instance
func NewMockFormatter(ctrl *gomock.Controller) *MockFormatter {
	mock := &MockFormatter{ctrl: ctrl}
	mock.recorder = &MockFormatterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFormatter) EXPECT() *MockFormatterMockRecorder {
	return m.recorder
}

// SplitCompanyNameEndings mocks base method
func (m *MockFormatter) SplitCompanyNameEndings(arg0 string) (string, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SplitCompanyNameEndings", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// SplitCompanyNameEndings indicates an expected call of SplitCompanyNameEndings
func (mr *MockFormatterMockRecorder) SplitCompanyNameEndings(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SplitCompanyNameEndings", reflect.TypeOf((*MockFormatter)(nil).SplitCompanyNameEndings), arg0)
}
