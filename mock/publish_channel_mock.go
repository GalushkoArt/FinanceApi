// Code generated by MockGen. DO NOT EDIT.
// Source: audit_publisher.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	audit "github.com/GalushkoArt/GoAuditService/pkg/proto"
	gomock "github.com/golang/mock/gomock"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

// MockAuditPublisher is a mock of AuditPublisher interface.
type MockAuditPublisher struct {
	ctrl     *gomock.Controller
	recorder *MockAuditPublisherMockRecorder
}

// MockAuditPublisherMockRecorder is the mock recorder for MockAuditPublisher.
type MockAuditPublisherMockRecorder struct {
	mock *MockAuditPublisher
}

// NewMockAuditPublisher creates a new mock instance.
func NewMockAuditPublisher(ctrl *gomock.Controller) *MockAuditPublisher {
	mock := &MockAuditPublisher{ctrl: ctrl}
	mock.recorder = &MockAuditPublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuditPublisher) EXPECT() *MockAuditPublisherMockRecorder {
	return m.recorder
}

// Publish mocks base method.
func (m *MockAuditPublisher) Publish(ctx context.Context, request *audit.LogRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockAuditPublisherMockRecorder) Publish(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockAuditPublisher)(nil).Publish), ctx, request)
}

// MockPublishChannel is a mock of PublishChannel interface.
type MockPublishChannel struct {
	ctrl     *gomock.Controller
	recorder *MockPublishChannelMockRecorder
}

// MockPublishChannelMockRecorder is the mock recorder for MockPublishChannel.
type MockPublishChannelMockRecorder struct {
	mock *MockPublishChannel
}

// NewMockPublishChannel creates a new mock instance.
func NewMockPublishChannel(ctrl *gomock.Controller) *MockPublishChannel {
	mock := &MockPublishChannel{ctrl: ctrl}
	mock.recorder = &MockPublishChannelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPublishChannel) EXPECT() *MockPublishChannelMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockPublishChannel) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockPublishChannelMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPublishChannel)(nil).Close))
}

// PublishWithContext mocks base method.
func (m *MockPublishChannel) PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishWithContext", ctx, exchange, key, mandatory, immediate, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishWithContext indicates an expected call of PublishWithContext.
func (mr *MockPublishChannelMockRecorder) PublishWithContext(ctx, exchange, key, mandatory, immediate, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishWithContext", reflect.TypeOf((*MockPublishChannel)(nil).PublishWithContext), ctx, exchange, key, mandatory, immediate, msg)
}