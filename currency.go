package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/marcelofabianov/fault"
)

// Currency represents a standardized currency code (e.g., BRL, USD, EUR).
// It ensures that only valid and recognized currency codes are used in the system,
// preventing errors in financial calculations and data exchange.
//
// The type is based on a predefined list of allowed currencies.
//
// Examples:
//   - BRL (Brazilian Real)
//   - USD (United States Dollar)
//   - EUR (Euro)
type Currency string

// Predefined and supported currency codes.
const (
	BRL Currency = "BRL" // Brazilian Real
	USD Currency = "USD" // United States Dollar
	EUR Currency = "EUR" // Euro
)

// EmptyCurrency represents the zero value for the Currency type.
var EmptyCurrency Currency

// validCurrencies holds the set of supported currency codes for validation.
var validCurrencies = map[Currency]struct{}{
	BRL: {},
	USD: {},
	EUR: {},
}

// NewCurrency creates a new Currency from a string code.
// The input is trimmed and converted to uppercase for consistent validation.
// Returns an error if the code is not in the list of valid currencies.
//
// Examples:
//   brl, err := NewCurrency("BRL")
//   usd, err := NewCurrency(" usd ") // Input is trimmed and uppercased
//   xxx, err := NewCurrency("XXX")   // Returns an error
func NewCurrency(value string) (Currency, error) {
	c := Currency(strings.ToUpper(strings.TrimSpace(value)))

	if c.IsZero() {
		return EmptyCurrency, nil
	}

	if !c.IsValid() {
		return EmptyCurrency, fault.New(
			"invalid currency code",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_code", value),
		)
	}
	return c, nil
}

// String returns the currency code as a string.
func (c Currency) String() string {
	return string(c)
}

// IsValid checks if the currency is in the list of supported currencies.
func (c Currency) IsValid() bool {
	_, ok := validCurrencies[c]
	return ok
}

// IsZero returns true if the currency is the zero value (EmptyCurrency).
func (c Currency) IsZero() bool {
	return c == EmptyCurrency
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Currency as a JSON string.
func (c Currency) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into a Currency, performing validation.
func (c *Currency) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*c = EmptyCurrency
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err,
			"currency must be a valid JSON string",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_json", string(data)),
		)
	}

	curr, err := NewCurrency(s)
	if err != nil {
		return err
	}

	*c = curr
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the currency code as a string or nil if it's the zero value.
func (c Currency) Value() (driver.Value, error) {
	if c.IsZero() {
		return nil, nil
	}
	return c.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts string or []byte values and validates them as a Currency.
func (c *Currency) Scan(src interface{}) error {
	if src == nil {
		*c = EmptyCurrency
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New(
			"unsupported scan type for Currency",
			fault.WithCode(fault.Invalid),
			fault.WithContext("received_type", fmt.Sprintf("%T", src)),
		)
	}

	curr, err := NewCurrency(s)
	if err != nil {
		return err
	}

	*c = curr
	return nil
}
