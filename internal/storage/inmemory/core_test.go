package inmemory

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndIsaev/go-metrics-alerter/internal/storage/file"
)

func TestNewMemStorage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Не удалось создать временную директорию: %v", err)
	}
	defer os.RemoveAll(tempDir)

	newManager, err := file.NewManager(tempDir)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	tests := []struct {
		name     string
		syncSave bool
		fm       *file.Manager
	}{
		{name: "#1", syncSave: false, fm: nil},
		{name: "#2", syncSave: true, fm: newManager},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemStorage(tt.fm, tt.syncSave)
			assert.Equal(t, tt.fm, storage.fm)
			assert.Equal(t, tt.syncSave, storage.syncSave)
		})
	}
}
