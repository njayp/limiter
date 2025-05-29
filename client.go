package limiter

import (
	"net/http"
)

// NewClient creates a new http.Client with the given options.
func NewClient(opts ...MiddlewareOpts) *http.Client {
	return &http.Client{
		Transport: NewMiddlewareRoundTripper(opts...),
	}
}

func InjectClient(client *http.Client, opts ...MiddlewareOpts) {
	// client.Transport should be overruled by opts
	opts = append([]MiddlewareOpts{WithRoundTripper(client.Transport)}, opts...)
	client.Transport = NewMiddlewareRoundTripper(opts...)
}
