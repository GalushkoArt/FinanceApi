// Code generated by MockGen. DO NOT EDIT.
// Source: ../service/symbol_service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	model "github.com/galushkoart/finance-api/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockSymbolService is a mock of SymbolService interface.
type MockSymbolService struct {
	ctrl     *gomock.Controller
	recorder *MockSymbolServiceMockRecorder
}

// MockSymbolServiceMockRecorder is the mock recorder for MockSymbolService.
type MockSymbolServiceMockRecorder struct {
	mock *MockSymbolService
}

// NewMockSymbolService creates a new mock instance.
func NewMockSymbolService(ctrl *gomock.Controller) *MockSymbolService {
	mock := &MockSymbolService{ctrl: ctrl}
	mock.recorder = &MockSymbolServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSymbolService) EXPECT() *MockSymbolServiceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockSymbolService) Add(ctx context.Context, symbol model.Symbol) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, symbol)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockSymbolServiceMockRecorder) Add(ctx, symbol interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockSymbolService)(nil).Add), ctx, symbol)
}

// Delete mocks base method.
func (m *MockSymbolService) Delete(ctx context.Context, symbolName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, symbolName)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSymbolServiceMockRecorder) Delete(ctx, symbolName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSymbolService)(nil).Delete), ctx, symbolName)
}

// GetAll mocks base method.
func (m *MockSymbolService) GetAll(ctx context.Context) ([]model.Symbol, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]model.Symbol)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockSymbolServiceMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockSymbolService)(nil).GetAll), ctx)
}

// GetBySymbol mocks base method.
func (m *MockSymbolService) GetBySymbol(ctx context.Context, name string) (model.Symbol, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBySymbol", ctx, name)
	ret0, _ := ret[0].(model.Symbol)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBySymbol indicates an expected call of GetBySymbol.
func (mr *MockSymbolServiceMockRecorder) GetBySymbol(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBySymbol", reflect.TypeOf((*MockSymbolService)(nil).GetBySymbol), ctx, name)
}

// Update mocks base method.
func (m *MockSymbolService) Update(ctx context.Context, symbol model.UpdateSymbol) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, symbol)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockSymbolServiceMockRecorder) Update(ctx, symbol interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSymbolService)(nil).Update), ctx, symbol)
}
