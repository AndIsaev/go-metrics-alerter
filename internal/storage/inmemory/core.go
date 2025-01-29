package inmemory

import (
	"log"
	"sync"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage/file"
)

// MemStorage storage in memory
type MemStorage struct {
	Metrics  map[string]common.Metrics
	fm       *file.Manager
	syncSave bool
	mu       sync.Mutex
}

// NewMemStorage init storage in memory
func NewMemStorage(fm *file.Manager, syncSave bool) *MemStorage {
	log.Printf("init in memory storage")
	return &MemStorage{
		Metrics:  make(map[string]common.Metrics),
		fm:       fm,
		syncSave: syncSave,
	}
}

func (m *MemStorage) System() storage.SystemRepository {
	return m
}

func (m *MemStorage) Metric() storage.MetricRepository {
	return m
}
