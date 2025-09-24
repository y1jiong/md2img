package util

import (
	"context"
	"time"
)

func Sleep(ctx context.Context, d time.Duration) {
	if d <= 0 {
		return
	}

	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
