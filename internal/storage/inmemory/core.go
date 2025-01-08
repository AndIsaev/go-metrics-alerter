package inmemory

import (
	"log"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage"
	"github.com/AndIsaev/go-metrics-alerter/internal/storage/file"
)

type MemStorage struct {
	Metrics  map[string]common.Metrics
	fm       *file.FileManager
	syncSave bool
}

func NewMemStorage(fm *file.FileManager, syncSave bool) *MemStorage {
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
