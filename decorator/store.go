// Package decorator demonstrates the Decorator pattern in Go.
//
// A decorator wraps an object that implements an interface and adds behaviour
// before or after delegating to the wrapped object. Because the decorator
// implements the same interface, it is completely transparent to callers —
// they never know (or care) how many layers are stacked.
//
// Stacking order matters: the outermost decorator runs first.
//
//	store := NewInMemoryStore()
//	store = NewLoggingStore(store, logger)   // logs every call
//	store = NewCachingStore(store)           // caches on top of logging
//	store = NewRetryStore(store, 3)          // retries on top of caching
package decorator

import (
	"errors"
	"fmt"
	"time"
)

// DataStore is the interface every layer — real or decorator — must satisfy.
type DataStore interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

// --- Concrete implementation ---

// InMemoryStore is the real, unwrapped store.
type InMemoryStore struct {
	data map[string]string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: make(map[string]string)}
}

func (s *InMemoryStore) Get(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("key %q not found", key)
	}
	return v, nil
}

func (s *InMemoryStore) Set(key, value string) error {
	s.data[key] = value
	return nil
}

func (s *InMemoryStore) Delete(key string) error {
	delete(s.data, key)
	return nil
}

// --- Decorators ---

// Logger is a minimal logging interface so the decorator isn't tied to fmt.
type Logger interface {
	Log(msg string)
}

// LoggingStore logs every operation then delegates to the wrapped store.
type LoggingStore struct {
	wrapped DataStore
	logger  Logger
}

func NewLoggingStore(ds DataStore, l Logger) *LoggingStore {
	return &LoggingStore{wrapped: ds, logger: l}
}

func (s *LoggingStore) Get(key string) (string, error) {
	s.logger.Log(fmt.Sprintf("Get(%q)", key))
	v, err := s.wrapped.Get(key)
	if err != nil {
		s.logger.Log(fmt.Sprintf("Get(%q) error: %v", key, err))
	}
	return v, err
}

func (s *LoggingStore) Set(key, value string) error {
	s.logger.Log(fmt.Sprintf("Set(%q, %q)", key, value))
	return s.wrapped.Set(key, value)
}

func (s *LoggingStore) Delete(key string) error {
	s.logger.Log(fmt.Sprintf("Delete(%q)", key))
	return s.wrapped.Delete(key)
}

// CachingStore memoises Get results; invalidates on Set/Delete.
type CachingStore struct {
	wrapped DataStore
	cache   map[string]string
}

func NewCachingStore(ds DataStore) *CachingStore {
	return &CachingStore{wrapped: ds, cache: make(map[string]string)}
}

func (s *CachingStore) Get(key string) (string, error) {
	if v, ok := s.cache[key]; ok {
		return v, nil
	}
	v, err := s.wrapped.Get(key)
	if err == nil {
		s.cache[key] = v
	}
	return v, err
}

func (s *CachingStore) Set(key, value string) error {
	delete(s.cache, key)
	return s.wrapped.Set(key, value)
}

func (s *CachingStore) Delete(key string) error {
	delete(s.cache, key)
	return s.wrapped.Delete(key)
}

// RetryStore retries failed Get/Set/Delete calls up to maxAttempts times.
type RetryStore struct {
	wrapped     DataStore
	maxAttempts int
	delay       time.Duration
}

func NewRetryStore(ds DataStore, maxAttempts int) *RetryStore {
	return &RetryStore{wrapped: ds, maxAttempts: maxAttempts, delay: 10 * time.Millisecond}
}

func (s *RetryStore) Get(key string) (string, error) {
	var err error
	for i := range s.maxAttempts {
		var v string
		v, err = s.wrapped.Get(key)
		if err == nil {
			return v, nil
		}
		if i < s.maxAttempts-1 {
			time.Sleep(s.delay)
		}
	}
	return "", fmt.Errorf("after %d attempts: %w", s.maxAttempts, err)
}

func (s *RetryStore) Set(key, value string) error {
	var err error
	for i := range s.maxAttempts {
		err = s.wrapped.Set(key, value)
		if err == nil {
			return nil
		}
		if i < s.maxAttempts-1 {
			time.Sleep(s.delay)
		}
	}
	return fmt.Errorf("after %d attempts: %w", s.maxAttempts, err)
}

func (s *RetryStore) Delete(key string) error {
	var err error
	for i := range s.maxAttempts {
		err = s.wrapped.Delete(key)
		if err == nil {
			return nil
		}
		if i < s.maxAttempts-1 {
			time.Sleep(s.delay)
		}
	}
	return fmt.Errorf("after %d attempts: %w", s.maxAttempts, err)
}

// ErrTransient is a sentinel used in tests to simulate a transient failure.
var ErrTransient = errors.New("transient error")
