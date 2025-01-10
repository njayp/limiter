# limiter

## Runner

### Motivation
Runner is designed to facilitate working with rate-limited api's. It starts jobs at the specified limited/staggered rate, regardless of when previous jobs finish.

### Examples
As a basic example, let's send 1000 requests at 200 requests per minute.
```go
r := limiter.NewRunner(200, time.Minute, 100*time.Millisecond)

for range 1000 {
	go func() {
		r.Run(ctx, func(ctx context.Context) error {
			return do(ctx) // send request
		})
	}()
}
```
As a more useful example, let's wrap a client and limiter into a limited client. LimitedClient.Do can be used thread-safely to send however many requests are desired without worrying about exceeding the rate limit.

```go
type LimitedClient struct {
	runner *limiter.Runner
	client *Client
}

func NewLimitedClient(limit int, interval time.Duration) *LimitedClient {
	return &LimitedClient{
		runner: limiter.NewRunner(limit, interval, 100*time.Millisecond),
		client: NewClient(),
	}
}

func (lc *LimitedClient) Do(ctx context.Context, req *Request) (*Response, error) {
	var resp *Response
	err := lc.runner.Run(ctx, func(ctx context.Context) error {
		var err error
		resp, err = lc.client.Do(ctx, req)
		return err
	})

	return resp, err
}
```