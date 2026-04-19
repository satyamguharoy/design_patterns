package strategy_test

import (
	"testing"

	s "design_patterns/strategy"
)

func TestCreditCardPayment_Success(t *testing.T) {
	processor := s.NewPaymentProcessor(s.NewCreditCardStrategy("1234567890123456"))
	if err := processor.Checkout(99.99); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCreditCardPayment_InvalidCard(t *testing.T) {
	processor := s.NewPaymentProcessor(s.NewCreditCardStrategy("123"))
	if err := processor.Checkout(99.99); err == nil {
		t.Fatal("expected error for invalid card, got nil")
	}
}

func TestPayPalPayment_Success(t *testing.T) {
	processor := s.NewPaymentProcessor(s.NewPayPalStrategy("user@example.com"))
	if err := processor.Checkout(49.00); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPayPalPayment_MissingEmail(t *testing.T) {
	processor := s.NewPaymentProcessor(s.NewPayPalStrategy(""))
	if err := processor.Checkout(49.00); err == nil {
		t.Fatal("expected error for missing email, got nil")
	}
}

func TestBankTransferPayment_Success(t *testing.T) {
	processor := s.NewPaymentProcessor(s.NewBankTransferStrategy("ACC-9876"))
	if err := processor.Checkout(500.00); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestRuntimeStrategySwap shows the core value of the pattern: same processor,
// different algorithm, no structural change.
func TestRuntimeStrategySwap(t *testing.T) {
	processor := s.NewPaymentProcessor(s.NewCreditCardStrategy("1234567890123456"))
	if err := processor.Checkout(100.00); err != nil {
		t.Fatalf("credit card payment failed: %v", err)
	}

	processor.SetStrategy(s.NewPayPalStrategy("user@example.com"))
	if err := processor.Checkout(200.00); err != nil {
		t.Fatalf("paypal payment failed: %v", err)
	}
}
