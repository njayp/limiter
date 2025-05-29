package client

import (
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client := NewHttpClient(WithCount(2))
	start := time.Now()

	for range 3 {
		client.Get("https://example.com")
	}

	// test should take more than 1 second
	duration := time.Since(start)
	expectedDuration := 1 * time.Second
	if duration < expectedDuration {
		t.Errorf("expected duration to be at least %v, got %v", expectedDuration, duration)
	}
}
