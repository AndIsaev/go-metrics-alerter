package file

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type FileManager struct {
	file     *os.File
	producer *json.Encoder
	consumer *json.Decoder
}

func NewFileManager(path string) (*FileManager, error) {
	log.Printf("init file manager")

	fullPath := fmt.Sprintf("%s/%s", path, "metrics.txt")
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	fm := &FileManager{
		file:     file,
		producer: json.NewEncoder(file),
		consumer: json.NewDecoder(file),
	}

	return fm, nil
}
