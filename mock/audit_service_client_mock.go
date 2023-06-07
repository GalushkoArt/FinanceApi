// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/GalushkoArt/GoAuditService/pkg/proto (interfaces: AuditServiceClient)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	audit "github.com/GalushkoArt/GoAuditService/pkg/proto"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockAuditServiceClient is a mock of AuditServiceClient interface.
type MockAuditServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuditServiceClientMockRecorder
}

// MockAuditServiceClientMockRecorder is the mock recorder for MockAuditServiceClient.
type MockAuditServiceClientMockRecorder struct {
	mock *MockAuditServiceClient
}

// NewMockAuditServiceClient creates a new mock instance.
func NewMockAuditServiceClient(ctrl *gomock.Controller) *MockAuditServiceClient {
	mock := &MockAuditServiceClient{ctrl: ctrl}
	mock.recorder = &MockAuditServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuditServiceClient) EXPECT() *MockAuditServiceClientMockRecorder {
	return m.recorder
}

// Log mocks base method.
func (m *MockAuditServiceClient) Log(arg0 context.Context, arg1 *audit.LogRequest, arg2 ...grpc.CallOption) (*audit.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Log", varargs...)
	ret0, _ := ret[0].(*audit.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Log indicates an expected call of Log.
func (mr *MockAuditServiceClientMockRecorder) Log(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockAuditServiceClient)(nil).Log), varargs...)
}
