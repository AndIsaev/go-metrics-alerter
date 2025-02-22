package utils

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

type Object func(ctx context.Context, metrics []common.Metrics) error

// Retry - декоратор, реализующий повторный вызов функции
func Retry(fn Object) Object {
	return func(ctx context.Context, metrics []common.Metrics) error {
		sleep := time.Second * 1
		var err error

		err = fn(ctx, metrics)

		if err != nil {
			for attempt := 1; attempt < 4; attempt++ {
				if err = fn(ctx, metrics); err != nil {
					log.Printf("attempt - %d; sleep - %v seconds\n", attempt, sleep)
					time.Sleep(sleep)
					sleep += time.Second * 2
				}
			}
		}
		return errors.Unwrap(err)
	}
}
