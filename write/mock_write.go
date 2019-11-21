package write

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockWrite is a mock of Write interface
type MockWrite struct {
	ctrl     *gomock.Controller
	recorder *MockWriteMockRecorder
}

// MockWriteMockRecorder is the mock recorder for MockWrite
type MockWriteMockRecorder struct {
	mock *MockWrite
}

// NewMockWrite creates a new mock instance
func NewMockWrite(ctrl *gomock.Controller) *MockWrite {
	mock := &MockWrite{ctrl: ctrl}
	mock.recorder = &MockWriteMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWrite) EXPECT() *MockWriteMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockWrite) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockWriteMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockWrite)(nil).Close))
}

// WriteToFile1 mocks base method
func (m *MockWrite) WriteToFile1(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteToFile1", arg0)
}

// WriteToFile1 indicates an expected call of WriteToFile1
func (mr *MockWriteMockRecorder) WriteToFile1(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteToFile1", reflect.TypeOf((*MockWrite)(nil).WriteToFile1), arg0)
}

// WriteToFile2 mocks base method
func (m *MockWrite) WriteToFile2(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteToFile2", arg0)
}

// WriteToFile2 indicates an expected call of WriteToFile2
func (mr *MockWriteMockRecorder) WriteToFile2(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteToFile2", reflect.TypeOf((*MockWrite)(nil).WriteToFile2), arg0)
}

// WriteToFile3 mocks base method
func (m *MockWrite) WriteToFile3(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteToFile3", arg0)
}

// WriteToFile3 indicates an expected call of WriteToFile3
func (mr *MockWriteMockRecorder) WriteToFile3(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteToFile3", reflect.TypeOf((*MockWrite)(nil).WriteToFile3), arg0)
}
