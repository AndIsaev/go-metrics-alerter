package utils

import (
	"errors"
	"log"
	"time"
)

type Object func() error

// Retry - декоратор, реализующий повторный вызов функции
func Retry(fn Object) Object {
	return func() error {
		sleep := time.Second * 1
		var err error

		err = fn()

		if err != nil {
			for attempt := 1; attempt < 4; attempt++ {
				if err = fn(); err != nil {
					log.Printf("attempt - %d; sleep - %v seconds\n", attempt, sleep)
					time.Sleep(sleep)
					sleep += time.Second * 2
				}
			}
		}
		return errors.Unwrap(err)
	}
}
