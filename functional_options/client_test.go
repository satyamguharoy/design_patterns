package functional_options_test

import (
	"testing"
	"time"

	fo "design_patterns/functional_options"
)

func TestDefaults(t *testing.T) {
	c, err := fo.NewHTTPClient()
	if err != nil {
		t.Fatalf("unexpected error with defaults: %v", err)
	}
	if c.Timeout() != 30*time.Second {
		t.Errorf("expected default timeout 30s, got %v", c.Timeout())
	}
	if c.MaxRetries() != 0 {
		t.Errorf("expected default maxRetries 0, got %d", c.MaxRetries())
	}
	if c.TLSEnabled() {
		t.Error("expected TLS disabled by default")
	}
}

func TestWithBaseURL(t *testing.T) {
	c, _ := fo.NewHTTPClient(fo.WithBaseURL("api.example.com"))
	if c.BaseURL() != "api.example.com" {
		t.Errorf("expected api.example.com, got %s", c.BaseURL())
	}
}

func TestWithTimeout(t *testing.T) {
	c, _ := fo.NewHTTPClient(fo.WithTimeout(5 * time.Second))
	if c.Timeout() != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", c.Timeout())
	}
}

func TestWithMaxRetries(t *testing.T) {
	c, _ := fo.NewHTTPClient(fo.WithMaxRetries(3))
	if c.MaxRetries() != 3 {
		t.Errorf("expected 3 retries, got %d", c.MaxRetries())
	}
}

func TestWithUserAgent(t *testing.T) {
	c, _ := fo.NewHTTPClient(fo.WithUserAgent("myapp/2.0"))
	if c.UserAgent() != "myapp/2.0" {
		t.Errorf("expected myapp/2.0, got %s", c.UserAgent())
	}
}

func TestWithTLS(t *testing.T) {
	c, _ := fo.NewHTTPClient(fo.WithTLS())
	if !c.TLSEnabled() {
		t.Error("expected TLS enabled")
	}
}

// TestOptionsComposed verifies that multiple options apply independently
// and last-write-wins when the same option is provided twice.
func TestOptionsComposed(t *testing.T) {
	c, err := fo.NewHTTPClient(
		fo.WithBaseURL("api.example.com"),
		fo.WithTimeout(10*time.Second),
		fo.WithMaxRetries(5),
		fo.WithTLS(),
		fo.WithTimeout(3*time.Second), // overrides the earlier WithTimeout
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Timeout() != 3*time.Second {
		t.Errorf("expected last timeout 3s to win, got %v", c.Timeout())
	}
	if c.MaxRetries() != 5 {
		t.Errorf("expected 5 retries, got %d", c.MaxRetries())
	}
}

func TestValidation_NegativeRetries(t *testing.T) {
	_, err := fo.NewHTTPClient(fo.WithMaxRetries(-1))
	if err == nil {
		t.Fatal("expected error for negative retries, got nil")
	}
}

func TestValidation_ZeroTimeout(t *testing.T) {
	_, err := fo.NewHTTPClient(fo.WithTimeout(0))
	if err == nil {
		t.Fatal("expected error for zero timeout, got nil")
	}
}
