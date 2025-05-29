# limiter

## Limiter

### Motivation
Limiter is designed to facilitate working with rate-limited APIs. It impl a token bucket that will allow a maximum number of function calls within a given interval.


### Install
```shell
go get github.com/njayp/limiter
```

### Client Examples
Prebuilt middleware is now available. The following client sends 2 requests/second.

```go
client := limiter.NewClient(WithCount(2))

for range 3 {
	client.Get("https://example.com")
}
```

### Rate Limiter Example
As a basic example, let's send 1000 requests at 200 requests per minute.
```go
l := limiter.NewLimiter(200, time.Minute, 100*time.Millisecond)

for range 1000 {
	go func() {
		// wait for turn
		err := l.Wait(ctx)
		if err != nil {
			// handle error
			return
		}

		// send request
		err = do(ctx)
		if err != nil {
			// handle error
		}
	}()
}
```

As a more useful example, let's wrap a client and limiter into a limited client. LimitedClient.Do can be used thread-safely to send however many requests are desired without worrying about exceeding the rate limit.
```go
type LimitedClient struct {
	limiter *limiter.Limiter
	client  *Client
}

func NewLimitedClient(limit int, interval time.Duration) *LimitedClient {
	return &LimitedClient{
		limiter: limiter.NewLimiter(limit, interval, 100*time.Millisecond),
		client:  NewClient(),
	}
}

func (lc *LimitedClient) Do(ctx context.Context, req *Request) (*Response, error) {
	err := lc.limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return lc.client.Do(ctx, req)
}
```