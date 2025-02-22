package file

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

// Overwrite save metrics to disc
func (fm *Manager) Overwrite(ctx context.Context, metrics []common.Metrics) error {
	if ctx.Err() != nil {
		log.Println("context is done -> exit from Overwrite")
		return ctx.Err()
	}
	fullPath := fm.file.Name()
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Print(err)
		return errors.Unwrap(err)
	}
	defer file.Close()

	fm.producer = json.NewEncoder(file)
	err = fm.producer.Encode(metrics)
	if err != nil {
		log.Printf("error save row to disc")
		return errors.Unwrap(err)
	}

	return nil
}

// ReadFile read from file
func (fm *Manager) ReadFile() ([]common.Metrics, error) {
	var result []common.Metrics
	if err := fm.consumer.Decode(&result); err != nil {
		log.Println("warning read metrics from file")
		return nil, errors.Unwrap(err)
	}

	return result, nil
}

func (fm *Manager) Close() error {
	err := fm.file.Close()
	if err != nil {
		log.Printf("error close file")
		return errors.Unwrap(err)
	}
	return nil
}
