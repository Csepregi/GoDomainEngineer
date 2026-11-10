package main

import "errors"
// Money is a monetary amount with a currency.
//
// Turn this into a proper value object:
//  1. Expose getters for `amount` and `currency` fields and make them unexported.
//  2. Reject invalid inputs in NewMoney (negative amount, empty currency).
//  3. Add an `Add` method that enforces the same-currency invariant.
type Money struct {
	amount   int64
	currency string
}

func (m *Money) Amount() int64 { return m.amount }
func (m *Money) Currency() string { return m.currency }

func NewMoney(amount int64, currency string) (Money, error) {
	// TODO: validate amount and currency.
	if amount < 0 {
		return Money{}, errors.New("amount must be positive")
	}
	if currency == "" {
		return Money{}, errors.New("currency must be non-empty")
	}
	return Money{amount: amount, currency: currency}, nil
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("currencies must match")
	}
	return Money{ amount: m.amount + other.amount, currency: m.currency}, nil
}