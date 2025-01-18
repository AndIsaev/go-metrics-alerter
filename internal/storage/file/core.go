package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

// FileManager manage file on disk
type FileManager struct {
	file     *os.File
	producer *json.Encoder
	consumer *json.Decoder
}

// NewFileManager init FileManager
func NewFileManager(path string) (*FileManager, error) {
	log.Printf("init file manager")

	fullPath := fmt.Sprintf("%s/%s", path, "metrics.txt")
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	fm := &FileManager{
		file:     file,
		producer: json.NewEncoder(file),
		consumer: json.NewDecoder(file),
	}

	return fm, nil
}
