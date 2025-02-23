package inmemory

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
)

func TestMemStorage_GetByNameType(t *testing.T) {
	delta := common.LinkInt64(100)
	value := common.LinkFloat64(12.34)
	inMemStorage := &MemStorage{
		Metrics: map[string]common.Metrics{
			"metric1": {ID: "metric1", MType: common.Counter, Delta: delta},
			"metric2": {ID: "metric2", MType: common.Gauge, Value: value},
		},
	}

	tests := []struct {
		name    string
		mName   string
		mType   string
		want    common.Metrics
		wantErr error
	}{
		{
			name:    "Valid Counter Metric",
			mName:   "metric1",
			mType:   common.Counter,
			want:    common.Metrics{ID: "metric1", MType: common.Counter, Delta: delta},
			wantErr: nil,
		},
		{
			name:    "Valid Gauge Metric",
			mName:   "metric2",
			mType:   common.Gauge,
			want:    common.Metrics{ID: "metric2", MType: common.Gauge, Value: value},
			wantErr: nil,
		},
		{
			name:    "Non-existent Metric",
			mName:   "unknown",
			mType:   common.Counter,
			wantErr: storage.ErrValueNotFound,
		},
		{
			name:    "Mismatch Metric Type",
			mName:   "metric1",
			mType:   common.Gauge,
			wantErr: storage.ErrValueNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := inMemStorage.GetByNameType(context.Background(), tt.mName, tt.mType)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MemStorage.GetByNameType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && (got != tt.want) {
				t.Errorf("MemStorage.GetByNameType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_create(t *testing.T) {
	delta := common.LinkInt64(100)
	value := common.LinkFloat64(12.34)
	inMemStorage := NewMemStorage(nil, false)

	tests := []struct {
		name    string
		mName   string
		mType   string
		want    common.Metrics
		wantErr error
	}{
		{
			name:    "Valid Counter Metric",
			mName:   "metric1",
			mType:   common.Counter,
			want:    common.Metrics{ID: "metric1", MType: common.Counter, Delta: delta},
			wantErr: nil,
		},
		{
			name:    "Valid Gauge Metric",
			mName:   "metric2",
			mType:   common.Gauge,
			want:    common.Metrics{ID: "metric2", MType: common.Gauge, Value: value},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := inMemStorage.create(context.Background(), tt.want)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MemStorage.create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := inMemStorage.GetByNameType(context.Background(), tt.mName, tt.mType)
			if (err == nil) && (got != tt.want) {
				t.Errorf("MemStorage.create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_List(t *testing.T) {
	delta := common.LinkInt64(100)
	value := common.LinkFloat64(12.34)
	inMemStorage := &MemStorage{
		Metrics: map[string]common.Metrics{
			"metric1": {ID: "metric1", MType: common.Counter, Delta: delta},
			"metric2": {ID: "metric2", MType: common.Gauge, Value: value},
		},
	}

	tests := []struct {
		name    string
		want    []common.Metrics
		wantErr error
		count   int
	}{
		{
			name: "Get list metrics",
			want: []common.Metrics{
				{ID: "metric1", MType: common.Counter, Delta: delta},
				{ID: "metric2", MType: common.Gauge, Value: value},
			},
			wantErr: nil,
			count:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := inMemStorage.List(context.Background())
			sort.Slice(got, func(i, j int) bool {
				return got[i].ID < got[j].ID
			})

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MemStorage.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.count, len(got))
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMemStorage_GetByName(t *testing.T) {
	delta := common.LinkInt64(100)
	value := common.LinkFloat64(12.34)
	inMemStorage := &MemStorage{
		Metrics: map[string]common.Metrics{
			"metric1": {ID: "metric1", MType: common.Counter, Delta: delta},
			"metric2": {ID: "metric2", MType: common.Gauge, Value: value},
		},
	}

	tests := []struct {
		name    string
		mName   string
		mType   string
		want    common.Metrics
		wantErr error
	}{
		{
			name:    "valid_counter",
			mName:   "metric1",
			want:    common.Metrics{ID: "metric1", MType: common.Counter, Delta: delta},
			wantErr: nil,
		},
		{
			name:    "valid_gauge",
			mName:   "metric2",
			want:    common.Metrics{ID: "metric2", MType: common.Gauge, Value: value},
			wantErr: nil,
		},
		{
			name:    "non_existent_metric",
			mName:   "unknown",
			wantErr: storage.ErrValueNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := inMemStorage.GetByName(context.Background(), tt.mName)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MemStorage.GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && (got != tt.want) {
				t.Errorf("MemStorage.GetByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_Insert(t *testing.T) {
	delta := common.LinkInt64(100)
	value := common.LinkFloat64(12.34)
	inMemStorage := NewMemStorage(nil, false)

	tests := []struct {
		name string
		want common.Metrics
	}{
		{"valid_counter", common.Metrics{ID: "metric1", MType: common.Counter, Delta: delta}},
		{"valid_gauge", common.Metrics{ID: "metric2", MType: common.Gauge, Value: value}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := inMemStorage.Insert(context.Background(), tt.want)
			existsMetric, err := inMemStorage.GetByName(context.Background(), tt.want.ID)
			assert.Equal(t, got, existsMetric)
			assert.NoError(t, err)
		})
	}
}

func TestMemStorage_InsertBatch(t *testing.T) {
	delta := common.LinkInt64(100)
	value := common.LinkFloat64(12.34)
	updatedValue := common.LinkInt64(*delta + *delta)
	inMemStorage := NewMemStorage(nil, false)
	inMemStorage.Insert(context.Background(), common.Metrics{ID: "existsMetric", MType: common.Counter, Delta: delta})

	tests := []struct {
		name string
		body []common.Metrics
		want []common.Metrics
	}{
		{
			name: "insert_batch",
			body: []common.Metrics{
				{ID: "existsMetric", MType: common.Counter, Delta: delta},
				{ID: "metric1", MType: common.Counter, Delta: delta},
				{ID: "metric2", MType: common.Gauge, Value: value},
			},
			want: []common.Metrics{
				{ID: "existsMetric", MType: common.Counter, Delta: updatedValue},
				{ID: "metric1", MType: common.Counter, Delta: delta},
				{ID: "metric2", MType: common.Gauge, Value: value},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = inMemStorage.InsertBatch(context.Background(), tt.body)
			got, err := inMemStorage.List(context.Background())
			sort.Slice(got, func(i, j int) bool {
				return got[i].ID < got[j].ID
			})

			if err == nil {
				assert.Equal(t, tt.want, got)
				return
			}
			t.Errorf("MemStorage.InsertBatch() = %v, want %v", got, tt.want)
		})
	}
}

func TestMemStorage_UpsertByValue(t *testing.T) {
	incorrectDelta := "100"
	incorrectValue := "12.34"
	intDelta := common.LinkInt64(100)
	floatValue := common.LinkFloat64(12.34)
	updatedValue := common.LinkInt64(*intDelta + *intDelta)
	inMemStorage := NewMemStorage(nil, false)
	inMemStorage.Insert(context.Background(), common.Metrics{ID: "metric5", MType: common.Counter, Delta: intDelta})

	tests := []struct {
		name    string
		body    common.Metrics
		value   any
		want    common.Metrics
		wantErr error
	}{
		{
			"valid_counter",
			common.Metrics{ID: "metric1", MType: common.Counter},
			*intDelta,
			common.Metrics{ID: "metric1", MType: common.Counter, Delta: intDelta},
			nil,
		},
		{
			"valid_gauge",
			common.Metrics{ID: "metric2", MType: common.Gauge},
			*floatValue,
			common.Metrics{ID: "metric2", MType: common.Gauge, Value: floatValue},
			nil,
		},
		{
			"incorrect_counter",
			common.Metrics{ID: "metric3", MType: common.Counter},
			incorrectDelta,
			common.Metrics{ID: "metric3", MType: common.Counter, Delta: intDelta},
			storage.ErrMetricValue,
		},
		{
			"incorrect_gauge",
			common.Metrics{ID: "metric4", MType: common.Gauge},
			incorrectValue,
			common.Metrics{ID: "metric4", MType: common.Gauge, Value: floatValue},
			storage.ErrMetricValue,
		},
		{
			"update_count",
			common.Metrics{ID: "metric5", MType: common.Counter},
			*intDelta,
			common.Metrics{ID: "metric5", MType: common.Counter, Delta: updatedValue},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := inMemStorage.UpsertByValue(context.Background(), tt.body, tt.value)
			got, _ := inMemStorage.GetByName(context.Background(), tt.want.ID)

			assert.ErrorIs(t, err, tt.wantErr)
			if err == nil {
				assert.Equal(t, tt.want, got)
				return
			}
			if !errors.Is(err, storage.ErrMetricValue) {
				t.Errorf("MemStorage.UpsertByValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
