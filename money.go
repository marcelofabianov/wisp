package wisp

import (
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

type Money struct {
	amount   int64
	currency Currency
}

var ZeroMoney = Money{}

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

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Currency() Currency {
	return m.currency
}

func (m Money) IsZero() bool {
	return m == ZeroMoney
}

func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

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

func (m Money) Multiply(multiplier int64) Money {
	return Money{
		amount:   m.amount * multiplier,
		currency: m.currency,
	}
}

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

func (m Money) Float64() float64 {
	return float64(m.amount) / 100.0
}

func (m Money) String() string {
	return fmt.Sprintf("%s %.2f", m.currency, m.Float64())
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Amount   int64    `json:"amount"`
		Currency Currency `json:"currency"`
	}{
		Amount:   m.amount,
		Currency: m.currency,
	})
}

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
