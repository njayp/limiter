package rate

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	ctx := context.Background()
	limit := 2
	interval := time.Second
	l := NewLimiter(limit, interval)
	wg := &sync.WaitGroup{}
	start := time.Now()

	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Wait(ctx)
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	expectedDuration := 1 * time.Second
	if duration < expectedDuration {
		t.Errorf("expected duration to be at least %v, got %v", expectedDuration, duration)
	}
}
