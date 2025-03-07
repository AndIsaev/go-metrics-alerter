// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/server/metrics.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	common "github.com/AndIsaev/go-metrics-alerter/internal/common"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CloseStorage mocks base method.
func (m *MockService) CloseStorage(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseStorage", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseStorage indicates an expected call of CloseStorage.
func (mr *MockServiceMockRecorder) CloseStorage(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseStorage", reflect.TypeOf((*MockService)(nil).CloseStorage), ctx)
}

// GetMetricByName mocks base method.
func (m *MockService) GetMetricByName(ctx context.Context, name string) (common.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetricByName", ctx, name)
	ret0, _ := ret[0].(common.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetricByName indicates an expected call of GetMetricByName.
func (mr *MockServiceMockRecorder) GetMetricByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetricByName", reflect.TypeOf((*MockService)(nil).GetMetricByName), ctx, name)
}

// GetMetricByNameType mocks base method.
func (m *MockService) GetMetricByNameType(ctx context.Context, name, mType string) (common.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetricByNameType", ctx, name, mType)
	ret0, _ := ret[0].(common.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetricByNameType indicates an expected call of GetMetricByNameType.
func (mr *MockServiceMockRecorder) GetMetricByNameType(ctx, name, mType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetricByNameType", reflect.TypeOf((*MockService)(nil).GetMetricByNameType), ctx, name, mType)
}

// InsertMetric mocks base method.
func (m *MockService) InsertMetric(ctx context.Context, metric common.Metrics) (common.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertMetric", ctx, metric)
	ret0, _ := ret[0].(common.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertMetric indicates an expected call of InsertMetric.
func (mr *MockServiceMockRecorder) InsertMetric(ctx, metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertMetric", reflect.TypeOf((*MockService)(nil).InsertMetric), ctx, metric)
}

// InsertMetrics mocks base method.
func (m *MockService) InsertMetrics(ctx context.Context, metrics []common.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertMetrics", ctx, metrics)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertMetrics indicates an expected call of InsertMetrics.
func (mr *MockServiceMockRecorder) InsertMetrics(ctx, metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertMetrics", reflect.TypeOf((*MockService)(nil).InsertMetrics), ctx, metrics)
}

// ListMetrics mocks base method.
func (m *MockService) ListMetrics(ctx context.Context) ([]common.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMetrics", ctx)
	ret0, _ := ret[0].([]common.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMetrics indicates an expected call of ListMetrics.
func (mr *MockServiceMockRecorder) ListMetrics(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMetrics", reflect.TypeOf((*MockService)(nil).ListMetrics), ctx)
}

// PingStorage mocks base method.
func (m *MockService) PingStorage(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingStorage", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// PingStorage indicates an expected call of PingStorage.
func (mr *MockServiceMockRecorder) PingStorage(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingStorage", reflect.TypeOf((*MockService)(nil).PingStorage), ctx)
}

// RunMigrationsStorage mocks base method.
func (m *MockService) RunMigrationsStorage(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunMigrationsStorage", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunMigrationsStorage indicates an expected call of RunMigrationsStorage.
func (mr *MockServiceMockRecorder) RunMigrationsStorage(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunMigrationsStorage", reflect.TypeOf((*MockService)(nil).RunMigrationsStorage), ctx)
}

// UpdateMetricByValue mocks base method.
func (m *MockService) UpdateMetricByValue(ctx context.Context, metric common.Metrics, value any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetricByValue", ctx, metric, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetricByValue indicates an expected call of UpdateMetricByValue.
func (mr *MockServiceMockRecorder) UpdateMetricByValue(ctx, metric, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricByValue", reflect.TypeOf((*MockService)(nil).UpdateMetricByValue), ctx, metric, value)
}
