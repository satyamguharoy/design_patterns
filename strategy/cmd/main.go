package main

import (
	"fmt"

	s "design_patterns/strategy"
)

func main() {
	processor := s.NewPaymentProcessor(s.NewCreditCardStrategy("1234567890123456"))
	processor.Checkout(99.99)

	fmt.Println("--- Customer switched to PayPal ---")
	processor.SetStrategy(s.NewPayPalStrategy("alice@example.com"))
	processor.Checkout(49.00)

	fmt.Println("--- Bulk order via bank transfer ---")
	processor.SetStrategy(s.NewBankTransferStrategy("ACC-9876"))
	processor.Checkout(1500.00)
}
