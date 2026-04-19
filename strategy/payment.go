// Package strategy demonstrates the Strategy pattern in Go.
//
// The Strategy pattern defines a family of algorithms, encapsulates each one
// behind an interface, and makes them interchangeable at runtime.
//
// Key difference from Dependency Injection: DI is about supplying dependencies;
// Strategy is specifically about swapping behaviour (algorithms) at runtime.
// In practice they often look the same in Go — both use interfaces — but the
// intent differs.
package strategy

import (
	"errors"
	"fmt"
)

// PaymentStrategy is the strategy interface. Every concrete payment method
// must implement this.
type PaymentStrategy interface {
	Pay(amount float64) error
	Name() string
}

// PaymentProcessor is the context. It holds a strategy and delegates to it.
// The processor itself never changes; only the strategy does.
type PaymentProcessor struct {
	strategy PaymentStrategy
}

func NewPaymentProcessor(s PaymentStrategy) *PaymentProcessor {
	return &PaymentProcessor{strategy: s}
}

// SetStrategy allows swapping the algorithm at runtime — the defining trait of
// the Strategy pattern.
func (p *PaymentProcessor) SetStrategy(s PaymentStrategy) {
	p.strategy = s
}

func (p *PaymentProcessor) Checkout(amount float64) error {
	fmt.Printf("Processing $%.2f via %s\n", amount, p.strategy.Name())
	return p.strategy.Pay(amount)
}

// --- Concrete Strategies ---

// CreditCardStrategy charges a credit card.
type CreditCardStrategy struct {
	cardNumber string
}

func NewCreditCardStrategy(cardNumber string) *CreditCardStrategy {
	return &CreditCardStrategy{cardNumber: cardNumber}
}

func (c *CreditCardStrategy) Name() string { return "Credit Card" }

func (c *CreditCardStrategy) Pay(amount float64) error {
	if len(c.cardNumber) < 16 {
		return errors.New("invalid card number")
	}
	fmt.Printf("  Charged $%.2f to card ending in %s\n", amount, c.cardNumber[len(c.cardNumber)-4:])
	return nil
}

// PayPalStrategy charges via PayPal.
type PayPalStrategy struct {
	email string
}

func NewPayPalStrategy(email string) *PayPalStrategy {
	return &PayPalStrategy{email: email}
}

func (p *PayPalStrategy) Name() string { return "PayPal" }

func (p *PayPalStrategy) Pay(amount float64) error {
	if p.email == "" {
		return errors.New("paypal email required")
	}
	fmt.Printf("  Charged $%.2f to PayPal account %s\n", amount, p.email)
	return nil
}

// BankTransferStrategy initiates a bank transfer.
type BankTransferStrategy struct {
	accountNumber string
}

func NewBankTransferStrategy(accountNumber string) *BankTransferStrategy {
	return &BankTransferStrategy{accountNumber: accountNumber}
}

func (b *BankTransferStrategy) Name() string { return "Bank Transfer" }

func (b *BankTransferStrategy) Pay(amount float64) error {
	if b.accountNumber == "" {
		return errors.New("account number required")
	}
	fmt.Printf("  Transferred $%.2f to account %s\n", amount, b.accountNumber)
	return nil
}
