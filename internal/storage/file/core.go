package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

type Provider interface {
	Overwrite(ctx context.Context, metrics []common.Metrics) error
	ReadFile() ([]common.Metrics, error)
	Close() error
}

// Manager manage file on disk
type Manager struct {
	file     *os.File
	producer *json.Encoder
	consumer *json.Decoder
}

// NewManager init Manager
func NewManager(path string) (*Manager, error) {
	log.Printf("init file manager")

	fullPath := fmt.Sprintf("%s/%s", path, "metrics.txt")
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	fm := &Manager{
		file:     file,
		producer: json.NewEncoder(file),
		consumer: json.NewDecoder(file),
	}

	return fm, nil
}

// CreateDir create dir
func CreateDir(fileStoragePath string) error {
	if _, err := os.Stat(fileStoragePath); os.IsNotExist(err) {
		if err = os.Mkdir(fileStoragePath, 0755); err != nil {
			log.Printf("the directory %s not created\n", fileStoragePath)
			return err
		}
	}
	log.Printf("the directory %s is done\n", fileStoragePath)
	return nil
}
