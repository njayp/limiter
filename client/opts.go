package client

import (
	"net/http"
	"time"
)

type ClientOpts func(*ClientProps)

// WithLimit sets the maximum number of requests per second.
func WithLimit(limit int) ClientOpts {
	return func(props *ClientProps) {
		props.Limit = limit
	}
}

// WithInterval sets the time interval for the rate limit.
func WithInterval(interval time.Duration) ClientOpts {
	return func(props *ClientProps) {
		props.Interval = interval
	}
}

// WithRoundTripper sets the http.RoundTripper to use for the client.
func WithRoundTripper(roundTripper http.RoundTripper) ClientOpts {
	return func(props *ClientProps) {
		props.RoundTripper = roundTripper
	}
}
