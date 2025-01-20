package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

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
