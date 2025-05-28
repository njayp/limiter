package limiter

import (
	"net/http"
	"time"
)

// MiddlewareRoundTripper is a middleware for http.RoundTripper
type MiddlewareRoundTripper struct {
	Original http.RoundTripper
	Limiter  *Limiter
}

func NewMiddlewareRoundTripper(limit int, interval time.Duration, concurrentLimit int) *MiddlewareRoundTripper {
	return &MiddlewareRoundTripper{
		Original: http.DefaultTransport,
		Limiter:  NewLimiter(limit, interval, concurrentLimit),
	}
}

func (c *MiddlewareRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// limit the number of requests per second and the number of concurrent requests
	token, err := c.Limiter.Wait(req.Context())
	if err != nil {
		return nil, err
	}
	defer token.Release() // Ensure the token is released after the request is done

	// Proceed with the actual request
	return c.Original.RoundTrip(req)
}

func HttpClient(limit int, interval time.Duration, concurrentLimit int) *http.Client {
	return &http.Client{
		Transport: NewMiddlewareRoundTripper(limit, interval, concurrentLimit),
	}
}
