package concurrent

import (
	"context"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	ctx := context.Background()
	limit := 2
	sem := NewSemaphore(limit)

	t.Run("Wait and Release", func(t *testing.T) {
		token, err := sem.Wait(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		token.Release() // should not panic
	})

	t.Run("Fill Semaphore", func(t *testing.T) {
		tokens := make([]*Token, 0, limit)
		for range limit {
			token, err := sem.Wait(ctx)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			tokens = append(tokens, token)
		}

		go func() {
			// test should take at least 1 second to complete
			time.Sleep(time.Second)

			for _, token := range tokens {
				token.Release() // release all tokens
			}
		}()

		_, err := sem.Wait(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}
