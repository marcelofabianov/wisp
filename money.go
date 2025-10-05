package wisp

import (
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

// Money represents a monetary amount with a specific currency.
// It stores the amount in the smallest currency unit (e.g., cents for USD, centavos for BRL)
// to avoid floating-point precision issues in financial calculations.
//
// The Money type ensures:
//   - Precise decimal arithmetic without floating-point errors
//   - Currency safety - operations only work with same currencies
//   - Immutability - operations return new instances
//   - Thread safety through immutable design
//
// Examples:
//   money, err := NewMoney(1050, BRL)  // R$ 10.50 (1050 centavos)
//   money, err := NewMoney(2500, USD)  // $25.00 (2500 cents)
//   fmt.Println(money.String())        // "BRL 10.50"
//   fmt.Println(money.Amount())        // 1050
//   fmt.Println(money.Currency())      // BRL
//
// All arithmetic operations (Add, Subtract, Multiply) return new Money instances
// and validate that currencies match when required.
type Money struct {
	amount   int64    // Amount in smallest currency unit (cents, centavos, etc.)
	currency Currency // The currency of this monetary amount
}

// ZeroMoney represents the zero value for Money type.
// It has zero amount and zero currency, making it invalid for most operations.
var ZeroMoney = Money{}

// NewMoney creates a new Money value with the specified amount and currency.
// The amount should be provided in the smallest currency unit (cents, centavos, etc.).
//
// Returns an error if the currency is invalid or zero.
//
// Examples:
//   money, err := NewMoney(1050, BRL)  // R$ 10.50
//   money, err := NewMoney(2500, USD)  // $25.00
//   money, err := NewMoney(0, BRL)     // R$ 0.00 (valid)
//   money, err := NewMoney(1000, Currency{}) // Error: invalid currency
func NewMoney(amountInCents int64, currency Currency) (Money, error) {
	if currency.IsZero() || !currency.IsValid() {
		return ZeroMoney, fault.New(
			"a valid currency is required to create money",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_currency", currency.String()),
		)
	}

	return Money{
		amount:   amountInCents,
		currency: currency,
	}, nil
}

// Amount returns the monetary amount in the smallest currency unit (e.g., cents).
func (m Money) Amount() int64 {
	return m.amount
}

// Currency returns the currency of the monetary amount.
func (m Money) Currency() Currency {
	return m.currency
}

// IsZero returns true if the Money is the zero value (ZeroMoney).
func (m Money) IsZero() bool {
	return m == ZeroMoney
}

// Equals checks if two Money instances are equal by comparing both amount and currency.
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// GreaterThan checks if the Money is greater than another.
// Returns an error if the currencies are different.
func (m Money) GreaterThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fault.New(
			"cannot compare money of different currencies",
			fault.WithCode(fault.DomainViolation),
			fault.WithContext("currency_a", m.currency),
			fault.WithContext("currency_b", other.currency),
		)
	}
	return m.amount > other.amount, nil
}

// LessThan checks if the Money is less than another.
// Returns an error if the currencies are different.
func (m Money) LessThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fault.New(
			"cannot compare money of different currencies",
			fault.WithCode(fault.DomainViolation),
			fault.WithContext("currency_a", m.currency),
			fault.WithContext("currency_b", other.currency),
		)
	}
	return m.amount < other.amount, nil
}

// Add returns a new Money instance with the sum of two amounts.
// Returns an error if the currencies are different.
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return ZeroMoney, fault.New(
			"cannot add money of different currencies",
			fault.WithCode(fault.DomainViolation),
			fault.WithContext("currency_a", m.currency),
			fault.WithContext("currency_b", other.currency),
		)
	}
	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract returns a new Money instance with the difference of two amounts.
// Returns an error if the currencies are different.
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return ZeroMoney, fault.New(
			"cannot subtract money of different currencies",
			fault.WithCode(fault.DomainViolation),
			fault.WithContext("currency_a", m.currency),
			fault.WithContext("currency_b", other.currency),
		)
	}
	return Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

// Multiply returns a new Money instance with the amount multiplied by a factor.
// This operation is useful for calculations like quantity * price.
func (m Money) Multiply(multiplier int64) Money {
	return Money{
		amount:   m.amount * multiplier,
		currency: m.currency,
	}
}

// Split divides the Money into n parts, distributing any remainder.
// This is useful for scenarios like splitting a bill among several people.
// The remainder is distributed one by one to the first parts.
// Returns an error if n is not a positive number.
func (m Money) Split(n int) ([]Money, error) {
	if n <= 0 {
		return nil, fault.New(
			"split count must be positive",
			fault.WithCode(fault.Invalid),
			fault.WithContext("split_count", n),
		)
	}

	parts := make([]Money, n)
	base := m.amount / int64(n)
	remainder := m.amount % int64(n)

	for i := 0; i < n; i++ {
		parts[i] = Money{amount: base, currency: m.currency}
		if remainder > 0 {
			parts[i].amount++
			remainder--
		}
	}

	return parts, nil
}

// IsNegative returns true if the monetary amount is negative.
func (m Money) IsNegative() bool {
	return m.amount < 0
}

// Float64 returns the monetary amount as a float64, converting from cents.
// Note: Use with caution, as floating-point arithmetic can lead to precision issues.
// This is primarily for display or interoperability, not for financial calculations.
func (m Money) Float64() float64 {
	return float64(m.amount) / 100.0
}

// String returns a formatted string representation of the money, like "BRL 10.50".
func (m Money) String() string {
	return fmt.Sprintf("%s %.2f", m.currency, m.Float64())
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes Money into a JSON object with "amount" and "currency" fields.
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Amount   int64    `json:"amount"`
		Currency Currency `json:"currency"`
	}{
		Amount:   m.amount,
		Currency: m.currency,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object into a Money instance, validating the currency.
func (m *Money) UnmarshalJSON(data []byte) error {
	dto := &struct {
		Amount   int64    `json:"amount"`
		Currency Currency `json:"currency"`
	}{}

	if err := json.Unmarshal(data, dto); err != nil {
		return fault.Wrap(err,
			"invalid JSON format for money",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_json", string(data)),
		)
	}

	if dto.Currency.IsZero() || !dto.Currency.IsValid() {
		return fault.New(
			"invalid or missing currency in JSON for money",
			fault.WithCode(fault.Invalid),
			fault.WithContext("received_currency", dto.Currency),
		)
	}

	m.amount = dto.Amount
	m.currency = dto.Currency

	return nil
}
