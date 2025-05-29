package client

import (
	"net/http"
	"time"
)

type ClientProps struct {
	// Count is the maximum number of requests per interval. The default is 10.
	Count int
	// Interval is the time interval for the rate limit. The default is 1 second.
	Interval time.Duration
	// RoundTripper is the http.RoundTripper to use for the client. The default is http.DefaultTransport.
	RoundTripper http.RoundTripper
}

// NewClientProps creates a new ClientProps with default values.
func NewClientProps(opts ...ClientOpts) *ClientProps {
	props := &ClientProps{
		Count:        10,                    // default limit
		Interval:     time.Second,           // default interval
		RoundTripper: http.DefaultTransport, // default RoundTripper
	}

	for _, opt := range opts {
		opt(props)
	}

	return props
}

// NewHttpClient creates a new http.Client with the given options.
func NewHttpClient(opts ...ClientOpts) *http.Client {
	props := NewClientProps(opts...)

	return &http.Client{
		Transport: NewMiddlewareRoundTripper(props),
	}
}
