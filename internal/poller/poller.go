package poller

import (
	"context"
	"fmt"
	"time"
)

// Poller runs fn every interval until context is canceled
func Poller(ctx context.Context, interval time.Duration, fn func()) {
	// fmt.Println("Poller started, running every", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fn()

	for {
		select {
		case <-ticker.C:
			fn()
		case <-ctx.Done():
			fmt.Println("Poller stopped via context")
			return
		}
	}
}
