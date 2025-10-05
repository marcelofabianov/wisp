package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

// CEP represents a Brazilian postal code (Código de Endereçamento Postal).
// It is a value object that ensures the code consists of exactly 8 digits.
// The value is stored as a string of digits but can be formatted for display.
//
// Examples:
//   - Input: "12345-678" or "12345678"
//   - Stored as: "12345678"
//   - Formatted output: "12345-678"
type CEP string

// EmptyCEP represents the zero value for the CEP type.
var EmptyCEP CEP

// parseCEP contains the core logic for validating and sanitizing a CEP string.
func parseCEP(input string) (CEP, error) {
	if input == "" {
		return EmptyCEP, nil
	}

	sanitized := nonDigitRegex.ReplaceAllString(input, "")

	if len(sanitized) != 8 {
		return EmptyCEP, fault.New(
			"CEP must have 8 digits",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input", input),
		)
	}

	return CEP(sanitized), nil
}

// NewCEP creates a new CEP from a string.
// It sanitizes the input by removing non-digit characters and validates that it has exactly 8 digits.
// Returns an error if the CEP is invalid.
func NewCEP(input string) (CEP, error) {
	return parseCEP(input)
}

// String returns the CEP as a string of 8 digits.
func (c CEP) String() string {
	return string(c)
}

// IsZero returns true if the CEP is the zero value.
func (c CEP) IsZero() bool {
	return c == EmptyCEP
}

// Formatted returns the CEP in the standard Brazilian format (XXXXX-XXX).
func (c CEP) Formatted() string {
	if len(c) != 8 {
		return c.String()
	}
	return fmt.Sprintf("%s-%s", c[0:5], c[5:8])
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the CEP to its 8-digit string representation.
func (c CEP) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into a CEP, with validation.
func (c *CEP) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "CEP must be a valid JSON string", fault.WithCode(fault.Invalid))
	}
	cep, err := NewCEP(s)
	if err != nil {
		return err
	}
	*c = cep
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the CEP as an 8-digit string.
func (c CEP) Value() (driver.Value, error) {
	if c.IsZero() {
		return nil, nil
	}
	return c.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a string or byte slice from the database and converts it into a CEP, with validation.
func (c *CEP) Scan(src interface{}) error {
	if src == nil {
		*c = EmptyCEP
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for CEP", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}

	cep, err := NewCEP(s)
	if err != nil {
		return err
	}
	*c = cep
	return nil
}
