package write

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
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

// Close mocks base method
func (m *MockWriter) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockWriterMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockWriter)(nil).Close))
}

// LogPostError mocks base method
func (m *MockWriter) LogPostError(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogPostError", arg0)
}

// LogPostError indicates an expected call of LogPostError
func (mr *MockWriterMockRecorder) LogPostError(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogPostError", reflect.TypeOf((*MockWriter)(nil).LogPostError), arg0)
}

// LogUnexpectedResponse mocks base method
func (m *MockWriter) LogUnexpectedResponse(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogUnexpectedResponse", arg0)
}

// LogUnexpectedResponse indicates an expected call of LogUnexpectedResponse
func (mr *MockWriterMockRecorder) LogUnexpectedResponse(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogUnexpectedResponse", reflect.TypeOf((*MockWriter)(nil).LogUnexpectedResponse), arg0)
}

// LogMissingCompanyName mocks base method
func (m *MockWriter) LogMissingCompanyName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogMissingCompanyName", arg0)
}

// LogMissingCompanyName indicates an expected call of LogMissingCompanyName
func (mr *MockWriterMockRecorder) LogMissingCompanyName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogMissingCompanyName", reflect.TypeOf((*MockWriter)(nil).LogMissingCompanyName), arg0)
}
