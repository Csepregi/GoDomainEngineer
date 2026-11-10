package main

import (
	"errors"
	"fmt"
)

// Money is a value object representing a monetary amount.
// Fields are unexported, so instances can only be created through NewMoney,
// which guarantees the invariants. Once constructed, a Money never changes.
type Money struct {
	amount   int64
	currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount must not be negative")
	}
	if currency == "" {
		return Money{}, errors.New("currency is required")
	}

	return Money{amount: amount, currency: currency}, nil
}

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

// Add returns a new Money. The original values are never mutated.
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot add %s and %s", m.currency, other.currency)
	}

	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

func main() {
	a, err := NewMoney(100, "USD")
	if err != nil {
		panic(err)
	}

	b, err := NewMoney(50, "USD")
	if err != nil {
		panic(err)
	}

	sum, err := a.Add(b)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d %s\n", sum.Amount(), sum.Currency())
}
