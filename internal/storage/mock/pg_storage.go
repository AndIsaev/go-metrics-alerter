// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/pg_storage.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPgStorage is a mock of PgStorage interface.
type MockPgStorage struct {
	ctrl     *gomock.Controller
	recorder *MockPgStorageMockRecorder
}

// MockPgStorageMockRecorder is the mock recorder for MockPgStorage.
type MockPgStorageMockRecorder struct {
	mock *MockPgStorage
}

// NewMockPgStorage creates a new mock instance.
func NewMockPgStorage(ctrl *gomock.Controller) *MockPgStorage {
	mock := &MockPgStorage{ctrl: ctrl}
	mock.recorder = &MockPgStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPgStorage) EXPECT() *MockPgStorageMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockPgStorage) Close(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockPgStorageMockRecorder) Close(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPgStorage)(nil).Close), ctx)
}

// Ping mocks base method.
func (m *MockPgStorage) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockPgStorageMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockPgStorage)(nil).Ping), ctx)
}
