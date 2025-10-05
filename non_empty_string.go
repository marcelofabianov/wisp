package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/marcelofabianov/fault"
)

// NonEmptyString is a value object ensuring a string is not empty after trimming whitespace.
// It is a simple way to enforce that required string fields, like names or titles, are always present.
//
// The zero value is EmptyNonEmptyString.
//
// Example:
//   name, err := NewNonEmptyString("  My Product  ")
//   fmt.Println(name) // "My Product"
//
//   _, err = NewNonEmptyString("   ") // returns an error
type NonEmptyString string

// EmptyNonEmptyString represents the zero value for NonEmptyString.
var EmptyNonEmptyString NonEmptyString

// NewNonEmptyString creates a new NonEmptyString.
// It trims whitespace from the input and returns an error if the result is an empty string.
func NewNonEmptyString(value string) (NonEmptyString, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return EmptyNonEmptyString, fault.New(
			"string cannot be empty",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", value),
		)
	}
	return NonEmptyString(trimmed), nil
}

// String returns the underlying string value.
func (s NonEmptyString) String() string {
	return string(s)
}

// IsZero returns true if the NonEmptyString is the zero value.
func (s NonEmptyString) IsZero() bool {
	return s == EmptyNonEmptyString
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the NonEmptyString to its string representation.
func (s NonEmptyString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into a NonEmptyString, with validation.
func (s *NonEmptyString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fault.Wrap(err, "NonEmptyString must be a valid JSON string", fault.WithCode(fault.Invalid))
	}

	nes, err := NewNonEmptyString(str)
	if err != nil {
		return err
	}
	*s = nes
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the NonEmptyString as a string.
func (s NonEmptyString) Value() (driver.Value, error) {
	return s.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a string or byte slice from the database and converts it into a NonEmptyString, with validation.
func (s *NonEmptyString) Scan(src interface{}) error {
	if src == nil {
		*s = EmptyNonEmptyString
		return nil
	}

	var str string
	switch v := src.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fault.New("unsupported scan type for NonEmptyString", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}

	nes, err := NewNonEmptyString(str)
	if err != nil {
		return err
	}
	*s = nes
	return nil
}
