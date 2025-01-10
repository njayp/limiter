package limiter

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunner(t *testing.T) { // Active jobs test
	ctx := context.Background()
	limit := 2
	interval := time.Second
	stagger := 0 * time.Millisecond
	r := NewRunner(limit, interval, stagger)
	count := atomic.Int32{}
	const expectedCount = 4

	fn := func(ctx context.Context) error {
		count.Add(1)
		time.Sleep(500 * time.Millisecond)
		return nil
	}

	wg := &sync.WaitGroup{}
	start := time.Now()
	for range expectedCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Run(ctx, fn)
		}()
	}

	time.Sleep(100 * time.Millisecond)
	activeJobs := r.ActiveJobs()
	if activeJobs != 2 {
		t.Errorf("expected 2 active jobs, got %d", activeJobs)
	}

	wg.Wait()
	duration := time.Since(start)

	load := count.Load()
	if load != expectedCount {
		t.Errorf("expected %d executions, got %d", expectedCount, load)
	}

	// last two jobs will start at t=1s and finish at t=1.5s
	expectedDuration := 1*time.Second + 500*time.Millisecond
	if duration < expectedDuration {
		t.Errorf("expected duration to be at least %v, got %v", expectedDuration, duration)
	}
}

func TestRunnerStagger(t *testing.T) {
	ctx := context.Background()
	limit := 10
	interval := time.Millisecond
	stagger := time.Second
	r := NewRunner(limit, interval, stagger)
	a := atomic.Int32{}

	fn := func(ctx context.Context) error {
		a.Add(1)
		return nil
	}

	for range 10 {
		go func() {
			r.Run(ctx, fn)
		}()
	}

	time.Sleep(time.Second + 100*time.Millisecond)
	got := a.Load()
	expectedCount := int32(2)
	if got != expectedCount {
		t.Errorf("expected %d executions, got %d", expectedCount, got)
	}
}
