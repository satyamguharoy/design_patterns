package decorator_test

import (
	"fmt"
	"testing"

	d "design_patterns/decorator"
)

// --- test helpers ---

type testLogger struct{ logs []string }

func (l *testLogger) Log(msg string) { l.logs = append(l.logs, msg) }

// flakyStore fails the first N calls to Get, then delegates.
type flakyStore struct {
	inner    d.DataStore
	failures int
	calls    int
}

func (f *flakyStore) Get(key string) (string, error) {
	f.calls++
	if f.calls <= f.failures {
		return "", fmt.Errorf("transient error (call %d)", f.calls)
	}
	return f.inner.Get(key)
}
func (f *flakyStore) Set(key, value string) error  { return f.inner.Set(key, value) }
func (f *flakyStore) Delete(key string) error      { return f.inner.Delete(key) }

// --- LoggingStore ---

func TestLoggingStore_LogsGetAndSet(t *testing.T) {
	logger := &testLogger{}
	store := d.NewLoggingStore(d.NewInMemoryStore(), logger)

	store.Set("k", "v")
	store.Get("k")

	if len(logger.logs) < 2 {
		t.Errorf("expected at least 2 log entries, got %d: %v", len(logger.logs), logger.logs)
	}
}

func TestLoggingStore_LogsGetError(t *testing.T) {
	logger := &testLogger{}
	store := d.NewLoggingStore(d.NewInMemoryStore(), logger)

	store.Get("missing")

	logged := false
	for _, l := range logger.logs {
		if len(l) > 0 {
			logged = true
		}
	}
	if !logged {
		t.Error("expected error to be logged")
	}
}

// --- CachingStore ---

func TestCachingStore_ReturnsCachedValue(t *testing.T) {
	inner := d.NewInMemoryStore()
	inner.Set("x", "original")

	cache := d.NewCachingStore(inner)
	cache.Get("x") // populate cache

	// mutate inner directly — cache should still return old value
	inner.Set("x", "updated")

	v, _ := cache.Get("x")
	if v != "original" {
		t.Errorf("expected cached value 'original', got %q", v)
	}
}

func TestCachingStore_InvalidatesOnSet(t *testing.T) {
	inner := d.NewInMemoryStore()
	cache := d.NewCachingStore(inner)

	cache.Set("x", "first")
	cache.Get("x") // populate cache

	cache.Set("x", "second") // should bust cache
	v, _ := cache.Get("x")
	if v != "second" {
		t.Errorf("expected 'second' after cache invalidation, got %q", v)
	}
}

func TestCachingStore_InvalidatesOnDelete(t *testing.T) {
	cache := d.NewCachingStore(d.NewInMemoryStore())
	cache.Set("x", "value")
	cache.Get("x")
	cache.Delete("x")

	_, err := cache.Get("x")
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

// --- RetryStore ---

func TestRetryStore_RetriesOnTransientFailure(t *testing.T) {
	inner := d.NewInMemoryStore()
	inner.Set("key", "value")

	flaky := &flakyStore{inner: inner, failures: 2}
	store := d.NewRetryStore(flaky, 3)

	v, err := store.Get("key")
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if v != "value" {
		t.Errorf("expected 'value', got %q", v)
	}
}

func TestRetryStore_FailsAfterMaxAttempts(t *testing.T) {
	inner := d.NewInMemoryStore()
	inner.Set("key", "value")

	flaky := &flakyStore{inner: inner, failures: 5}
	store := d.NewRetryStore(flaky, 3)

	_, err := store.Get("key")
	if err == nil {
		t.Fatal("expected error after exhausting retries, got nil")
	}
}

// --- Stacking ---

func TestStackedDecorators(t *testing.T) {
	logger := &testLogger{}

	// Stack: RetryStore -> CachingStore -> LoggingStore -> InMemoryStore
	var store d.DataStore = d.NewInMemoryStore()
	store = d.NewLoggingStore(store, logger)
	store = d.NewCachingStore(store)
	store = d.NewRetryStore(store, 2)

	store.Set("hello", "world")
	v, err := store.Get("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "world" {
		t.Errorf("expected 'world', got %q", v)
	}
	if len(logger.logs) == 0 {
		t.Error("expected log entries from stacked decorators")
	}
}
