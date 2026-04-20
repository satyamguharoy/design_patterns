package main

import (
	"fmt"

	d "design_patterns/decorator"
)

type consoleLogger struct{}

func (l *consoleLogger) Log(msg string) { fmt.Println("[LOG]", msg) }

func main() {
	// Build the stack from inside out.
	// Call order on Get("user:1"): Retry -> Cache -> Logging -> InMemory
	var store d.DataStore = d.NewInMemoryStore()
	store = d.NewLoggingStore(store, &consoleLogger{})
	store = d.NewCachingStore(store)
	store = d.NewRetryStore(store, 3)

	store.Set("user:1", "alice")

	fmt.Println("--- First Get (cache miss, hits store) ---")
	v, _ := store.Get("user:1")
	fmt.Println("Got:", v)

	fmt.Println("--- Second Get (cache hit, no log) ---")
	v, _ = store.Get("user:1")
	fmt.Println("Got:", v)

	fmt.Println("--- Delete then Get (miss) ---")
	store.Delete("user:1")
	_, err := store.Get("user:1")
	fmt.Println("Error:", err)
}
