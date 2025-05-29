package client

import (
	"log/slog"
	"net/http"

	"github.com/njayp/limiter/client/rate"
)

// MiddlewareRoundTripper is a middleware for http.RoundTripper
// TODO different limiter per host
type MiddlewareRoundTripper struct {
	Original http.RoundTripper
	Limiter  *rate.Limiter
}

// NewMiddlewareRoundTripper creates a new MiddlewareRoundTripper with the given rate limit and interval.
// use the default http.RoundTripper
func NewMiddlewareRoundTripper(props *ClientProps) *MiddlewareRoundTripper {
	return &MiddlewareRoundTripper{
		Original: props.RoundTripper,
		Limiter:  rate.NewLimiter(props.Limit, props.Interval),
	}
}

func (c *MiddlewareRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// limit the number of requests per second and the number of concurrent requests
	err := c.Limiter.Wait(req.Context())
	if err != nil {
		return nil, err
	}

	slog.Debug("Limiter Client Starting RoundTrip",
		"method", req.Method,
		"url", req.URL.String(),
	)

	// Proceed with the actual request
	return c.Original.RoundTrip(req)
}
