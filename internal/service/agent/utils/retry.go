package utils

import (
	"errors"
	"log"
	"time"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/metrics"
)

type Object func(url string, db *metrics.StorageMetrics) error

// Retry - декоратор, реализующий повторный вызов функции
func Retry(fn Object) Object {
	return func(url string, db *metrics.StorageMetrics) error {
		sleep := time.Second * 1
		var err error

		err = fn(url, db)

		if err != nil {
			for attempt := 3; attempt > 0; attempt-- {
				if err = fn(url, db); err != nil {
					log.Printf("attempt - %d; sleep - %v seconds\n", attempt, sleep)
					time.Sleep(sleep)
					sleep += time.Second * 2
				}
			}
		}
		return errors.Unwrap(err)
	}
}
