package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

// PositiveInt is a value object ensuring an integer is always strictly greater than zero.
// This is useful for representing values like counts, quantities, or IDs that cannot be zero or negative.
//
// The zero value is ZeroPositiveInt.
//
// Example:
//   count, err := NewPositiveInt(10)
//
//   _, err = NewPositiveInt(0) // returns an error
type PositiveInt int

// ZeroPositiveInt represents the zero value for PositiveInt.
var ZeroPositiveInt PositiveInt

// NewPositiveInt creates a new PositiveInt.
// It returns an error if the value is not strictly greater than zero.
func NewPositiveInt(value int) (PositiveInt, error) {
	if value <= 0 {
		return ZeroPositiveInt, fault.New(
			"value must be a positive integer",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", value),
		)
	}
	return PositiveInt(value), nil
}

// Int returns the underlying integer value.
func (p PositiveInt) Int() int {
	return int(p)
}

// IsZero returns true if the PositiveInt is the zero value.
func (p PositiveInt) IsZero() bool {
	return p == ZeroPositiveInt
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the PositiveInt to its integer representation.
func (p PositiveInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Int())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON number into a PositiveInt, with validation.
func (p *PositiveInt) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return fault.Wrap(err, "PositiveInt must be a valid JSON number", fault.WithCode(fault.Invalid))
	}

	pi, err := NewPositiveInt(i)
	if err != nil {
		return err
	}
	*p = pi
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the PositiveInt as an int64.
func (p PositiveInt) Value() (driver.Value, error) {
	return int64(p.Int()), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts an int64 from the database and converts it into a PositiveInt, with validation.
func (p *PositiveInt) Scan(src interface{}) error {
	if src == nil {
		*p = ZeroPositiveInt
		return nil
	}

	var i int64
	switch v := src.(type) {
	case int64:
		i = v
	default:
		return fault.New("unsupported scan type for PositiveInt", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}

	pi, err := NewPositiveInt(int(i))
	if err != nil {
		return err
	}
	*p = pi
	return nil
}
