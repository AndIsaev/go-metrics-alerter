package utils

import (
	"context"
	"errors"
	"testing"

	"github.com/AndIsaev/go-metrics-alerter/internal/common"
)

func TestRetry(t *testing.T) {
	ctx := context.Background()
	var metrics []common.Metrics

	callCount := 0

	fn := func(ctx context.Context, metrics []common.Metrics) error {
		callCount++
		if callCount < 4 {
			return errors.New("temporary error")
		}
		return nil
	}

	decoratedFn := Retry(fn)

	err := decoratedFn(ctx, metrics)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if callCount != 4 {
		t.Errorf("expected 4 attempts, got %d", callCount)
	}
}
