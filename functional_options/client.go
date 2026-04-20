// Package functional_options demonstrates the Functional Options pattern in Go.
//
// Problem: a struct with many optional fields leads to one of these bad outcomes:
//   - A constructor with a long parameter list (order-dependent, hard to read)
//   - Many constructor variants: NewClientWithTLS, NewClientWithTimeout, …
//   - A separate Config struct (verbose — callers must build it just to pass it)
//
// Solution: define a private options struct, expose it only through
// Option functions, and accept ...Option in the constructor.
//
//	client := NewHTTPClient(
//	    WithBaseURL("https://api.example.com"),
//	    WithTimeout(5 * time.Second),
//	    WithMaxRetries(3),
//	)
//
// Each option is self-documenting, order-independent, and easy to add without
// breaking existing callers.
package functional_options

import (
	"errors"
	"fmt"
	"time"
)

// options holds every configurable field. It is unexported — callers never
// build it directly; they use Option functions instead.
type options struct {
	baseURL    string
	timeout    time.Duration
	maxRetries int
	userAgent  string
	tlsEnabled bool
}

// defaults returns a safe baseline so callers that pass zero options still
// get a usable client.
func defaults() options {
	return options{
		baseURL:    "localhost",
		timeout:    30 * time.Second,
		maxRetries: 0,
		userAgent:  "go-http-client/1.0",
		tlsEnabled: false,
	}
}

// Option is a function that mutates the options struct.
// This is the type callers work with.
type Option func(*options)

// WithBaseURL sets the base URL for all requests.
func WithBaseURL(url string) Option {
	return func(o *options) {
		o.baseURL = url
	}
}

// WithTimeout sets the per-request timeout.
func WithTimeout(d time.Duration) Option {
	return func(o *options) {
		o.timeout = d
	}
}

// WithMaxRetries sets how many times a failed request is retried.
func WithMaxRetries(n int) Option {
	return func(o *options) {
		o.maxRetries = n
	}
}

// WithUserAgent overrides the User-Agent header.
func WithUserAgent(ua string) Option {
	return func(o *options) {
		o.userAgent = ua
	}
}

// WithTLS enables TLS for all requests.
func WithTLS() Option {
	return func(o *options) {
		o.tlsEnabled = true
	}
}

// HTTPClient is configured entirely through Option functions.
type HTTPClient struct {
	opts options
}

// NewHTTPClient applies defaults then overlays the provided options.
// Adding a new option never requires changing this signature.
func NewHTTPClient(opts ...Option) (*HTTPClient, error) {
	o := defaults()
	for _, opt := range opts {
		opt(&o)
	}
	if err := validate(o); err != nil {
		return nil, err
	}
	return &HTTPClient{opts: o}, nil
}

func validate(o options) error {
	if o.baseURL == "" {
		return errors.New("base URL must not be empty")
	}
	if o.timeout <= 0 {
		return errors.New("timeout must be positive")
	}
	if o.maxRetries < 0 {
		return errors.New("maxRetries must be non-negative")
	}
	return nil
}

// Get simulates an HTTP GET to demonstrate the configured client in action.
func (c *HTTPClient) Get(path string) string {
	scheme := "http"
	if c.opts.tlsEnabled {
		scheme = "https"
	}
	return fmt.Sprintf("[%s] GET %s://%s%s (timeout=%s retries=%d ua=%s)",
		c.opts.userAgent,
		scheme,
		c.opts.baseURL,
		path,
		c.opts.timeout,
		c.opts.maxRetries,
		c.opts.userAgent,
	)
}

// BaseURL, Timeout, MaxRetries expose config for testing without making
// the options struct public.
func (c *HTTPClient) BaseURL() string        { return c.opts.baseURL }
func (c *HTTPClient) Timeout() time.Duration { return c.opts.timeout }
func (c *HTTPClient) MaxRetries() int        { return c.opts.maxRetries }
func (c *HTTPClient) UserAgent() string      { return c.opts.userAgent }
func (c *HTTPClient) TLSEnabled() bool       { return c.opts.tlsEnabled }
