package main

import (
	"fmt"
	"time"

	fo "design_patterns/functional_options"
)

func main() {
	// Minimal client — safe defaults apply automatically.
	minimal, _ := fo.NewHTTPClient()
	fmt.Println(minimal.Get("/health"))

	// Fully configured client — each option is named and order-independent.
	full, _ := fo.NewHTTPClient(
		fo.WithBaseURL("api.example.com"),
		fo.WithTimeout(5*time.Second),
		fo.WithMaxRetries(3),
		fo.WithUserAgent("myapp/2.0"),
		fo.WithTLS(),
	)
	fmt.Println(full.Get("/users"))

	// Validation catches bad config at construction time, not at request time.
	_, err := fo.NewHTTPClient(fo.WithMaxRetries(-1))
	fmt.Println("Invalid config error:", err)
}
