// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/base.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	common "github.com/AndIsaev/go-metrics-alerter/internal/common"
	storage "github.com/AndIsaev/go-metrics-alerter/internal/storage"
	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Metric mocks base method.
func (m *MockStorage) Metric() storage.MetricRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Metric")
	ret0, _ := ret[0].(storage.MetricRepository)
	return ret0
}

// Metric indicates an expected call of Metric.
func (mr *MockStorageMockRecorder) Metric() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Metric", reflect.TypeOf((*MockStorage)(nil).Metric))
}

// System mocks base method.
func (m *MockStorage) System() storage.SystemRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "System")
	ret0, _ := ret[0].(storage.SystemRepository)
	return ret0
}

// System indicates an expected call of System.
func (mr *MockStorageMockRecorder) System() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "System", reflect.TypeOf((*MockStorage)(nil).System))
}

// MockSystemRepository is a mock of SystemRepository interface.
type MockSystemRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSystemRepositoryMockRecorder
}

// MockSystemRepositoryMockRecorder is the mock recorder for MockSystemRepository.
type MockSystemRepositoryMockRecorder struct {
	mock *MockSystemRepository
}

// NewMockSystemRepository creates a new mock instance.
func NewMockSystemRepository(ctrl *gomock.Controller) *MockSystemRepository {
	mock := &MockSystemRepository{ctrl: ctrl}
	mock.recorder = &MockSystemRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSystemRepository) EXPECT() *MockSystemRepositoryMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockSystemRepository) Close(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockSystemRepositoryMockRecorder) Close(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSystemRepository)(nil).Close), ctx)
}

// Ping mocks base method.
func (m *MockSystemRepository) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockSystemRepositoryMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockSystemRepository)(nil).Ping), ctx)
}

// RunMigrations mocks base method.
func (m *MockSystemRepository) RunMigrations(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunMigrations", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunMigrations indicates an expected call of RunMigrations.
func (mr *MockSystemRepositoryMockRecorder) RunMigrations(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunMigrations", reflect.TypeOf((*MockSystemRepository)(nil).RunMigrations), ctx)
}

// MockMetricRepository is a mock of MetricRepository interface.
type MockMetricRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMetricRepositoryMockRecorder
}

// MockMetricRepositoryMockRecorder is the mock recorder for MockMetricRepository.
type MockMetricRepositoryMockRecorder struct {
	mock *MockMetricRepository
}

// NewMockMetricRepository creates a new mock instance.
func NewMockMetricRepository(ctrl *gomock.Controller) *MockMetricRepository {
	mock := &MockMetricRepository{ctrl: ctrl}
	mock.recorder = &MockMetricRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricRepository) EXPECT() *MockMetricRepositoryMockRecorder {
	return m.recorder
}

// GetByName mocks base method.
func (m *MockMetricRepository) GetByName(ctx context.Context, name string) (common.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, name)
	ret0, _ := ret[0].(common.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockMetricRepositoryMockRecorder) GetByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockMetricRepository)(nil).GetByName), ctx, name)
}

// GetByNameType mocks base method.
func (m *MockMetricRepository) GetByNameType(ctx context.Context, name, mType string) (common.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNameType", ctx, name, mType)
	ret0, _ := ret[0].(common.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNameType indicates an expected call of GetByNameType.
func (mr *MockMetricRepositoryMockRecorder) GetByNameType(ctx, name, mType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNameType", reflect.TypeOf((*MockMetricRepository)(nil).GetByNameType), ctx, name, mType)
}

// Insert mocks base method.
func (m *MockMetricRepository) Insert(ctx context.Context, metric common.Metrics) (common.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, metric)
	ret0, _ := ret[0].(common.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockMetricRepositoryMockRecorder) Insert(ctx, metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockMetricRepository)(nil).Insert), ctx, metric)
}

// InsertBatch mocks base method.
func (m *MockMetricRepository) InsertBatch(ctx context.Context, metrics []common.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertBatch", ctx, metrics)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertBatch indicates an expected call of InsertBatch.
func (mr *MockMetricRepositoryMockRecorder) InsertBatch(ctx, metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertBatch", reflect.TypeOf((*MockMetricRepository)(nil).InsertBatch), ctx, metrics)
}

// List mocks base method.
func (m *MockMetricRepository) List(ctx context.Context) ([]common.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]common.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockMetricRepositoryMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockMetricRepository)(nil).List), ctx)
}

// UpsertByValue mocks base method.
func (m *MockMetricRepository) UpsertByValue(ctx context.Context, metric common.Metrics, value any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertByValue", ctx, metric, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertByValue indicates an expected call of UpsertByValue.
func (mr *MockMetricRepositoryMockRecorder) UpsertByValue(ctx, metric, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertByValue", reflect.TypeOf((*MockMetricRepository)(nil).UpsertByValue), ctx, metric, value)
}
