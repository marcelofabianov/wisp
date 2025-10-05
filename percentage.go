package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"

	"github.com/marcelofabianov/fault"
)

// Percentage represents a percentage value, stored as a scaled integer to avoid floating-point inaccuracies.
// It is designed for precise calculations in financial and business logic, such as discounts, interest rates, or tax rates.
//
// The value is stored scaled by a factor of 10,000 (representing 4 decimal places of precision).
// For example, a float value of 0.5 (representing 50%) is stored as the integer 5000.
//
// Examples:
//   - Input: 0.5 (50%)
//   - Stored as: 5000
//   - Formatted output: "50.00%"
//
// This approach ensures that arithmetic operations are precise.
type Percentage int64

// ZeroPercentage represents the zero value for the Percentage type.
var ZeroPercentage Percentage

// percentageFactor is the scaling factor used to store the percentage as an integer.
// A factor of 10,000 allows for 4 decimal places of precision (e.g., 0.0001 becomes 1).
const percentageFactor = 10000.0

// NewPercentageFromFloat creates a new Percentage from a float64 value.
// The float represents the percentage fraction (e.g., 0.5 for 50%).
// The value is scaled and rounded to the nearest even number to be stored as an integer.
//
// Returns an error if the input value is negative.
//
// Examples:
//   p, err := NewPercentageFromFloat(0.5)   // 50%
//   p, err := NewPercentageFromFloat(0.075) // 7.5%
//   p, err := NewPercentageFromFloat(-0.1)  // returns an error
func NewPercentageFromFloat(value float64) (Percentage, error) {
	if value < 0 {
		return ZeroPercentage, fault.New(
			"percentage value cannot be negative",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", value),
		)
	}

	scaledValue := math.RoundToEven(value * percentageFactor)
	return Percentage(scaledValue), nil
}

// IsNegative returns true if the percentage value is negative.
func (p Percentage) IsNegative() bool {
	return p < 0
}

// Float64 converts the scaled integer back to a float64 representation (e.g., 5000 becomes 0.5).
// This is useful for display or interoperability but should be used with caution in calculations
// due to potential floating-point inaccuracies.
func (p Percentage) Float64() float64 {
	return float64(p) / percentageFactor
}

// String returns a formatted string representation of the percentage (e.g., "50.00%").
func (p Percentage) String() string {
	return fmt.Sprintf("%.2f%%", p.Float64()*100)
}

// IsZero returns true if the percentage is the zero value.
func (p Percentage) IsZero() bool {
	return p == ZeroPercentage
}

// ApplyTo calculates the percentage of a given Money value.
// It returns a new Money instance representing the calculated amount.
// The result is rounded to the nearest smallest currency unit (e.g., cent).
//
// Example:
//   price := wisp.NewMoney(10000, wisp.BRL) // R$100.00
//   discount, _ := wisp.NewPercentageFromFloat(0.1) // 10%
//   discountAmount := discount.ApplyTo(price) // R$10.00 (1000 centavos)
func (p Percentage) ApplyTo(m Money) Money {
	if m.IsZero() || p.IsZero() {
		return Money{amount: 0, currency: m.Currency()}
	}

	result := float64(m.Amount()) * p.Float64()
	roundedAmount := int64(math.RoundToEven(result))

	return Money{
		amount:   roundedAmount,
		currency: m.Currency(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Percentage as its float64 representation.
func (p Percentage) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Float64())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON number (float64) into a Percentage, performing validation.
func (p *Percentage) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return fault.Wrap(err, "Percentage must be a valid JSON number", fault.WithCode(fault.Invalid))
	}
	perc, err := NewPercentageFromFloat(f)
	if err != nil {
		return err
	}
	*p = perc
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the scaled integer representation of the percentage.
func (p Percentage) Value() (driver.Value, error) {
	return int64(p), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts an int64 from the database and converts it into a Percentage.
func (p *Percentage) Scan(src interface{}) error {
	if src == nil {
		*p = ZeroPercentage
		return nil
	}

	var intVal int64
	switch v := src.(type) {
	case int64:
		intVal = v
	default:
		return fault.New("unsupported scan type for Percentage", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}

	if intVal < 0 {
		return fault.New("percentage from database cannot be negative", fault.WithCode(fault.Invalid), fault.WithContext("source_value", intVal))
	}

	*p = Percentage(intVal)
	return nil
}
