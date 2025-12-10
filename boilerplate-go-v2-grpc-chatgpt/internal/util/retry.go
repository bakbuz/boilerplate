package util

import (
	"context"
	"time"
)

func Retry(ctx context.Context, attempts int, base time.Duration, fn func() error) error {
	var err error
	d := base
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		select {
		case <-time.After(d):
			d *= 2
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}
