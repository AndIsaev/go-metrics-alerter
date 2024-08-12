package file

import (
	"encoding/json"
	"fmt"
	"github.com/AndIsaev/go-metrics-alerter/internal/common"
	"os"
)

type Producer struct {
	file *os.File
	// добавляем Writer в Producer
	encoder *json.Encoder
}

func NewProducer(path string) (*Producer, error) {
	fullPath := fmt.Sprintf("%s/%s", path, "metrics.txt")
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteMetrics(metrics *common.Metrics) error {
	return p.encoder.Encode(&metrics)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(path string) (*Consumer, error) {
	fullPath := fmt.Sprintf("%s/%s", path, "metrics.txt")
	file, err := os.OpenFile(fullPath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file: file,
		// создаём новый scanner
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadMetrics() (*common.Metrics, error) {
	// одиночное сканирование до следующей строки
	event := &common.Metrics{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}

// Close - close file
func (c *Consumer) Close() error {
	return c.file.Close()
}
