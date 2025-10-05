package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/marcelofabianov/fault"
)

// CPF represents a Brazilian individual taxpayer identification number (Cadastro de Pessoa Física).
// It validates the format and verifies check digits according to Brazilian government standards.
// The value is stored without formatting (digits only) but can be displayed with proper formatting.
//
// Examples:
//   - Input: "123.456.789-09" or "12345678909"
//   - Storage: "12345678909"
//   - Formatted output: "123.456.789-09"
//
// A CPF is considered valid when:
//   - It contains exactly 11 digits
//   - It's not a sequence of repeated digits (e.g., "11111111111")
//   - Both check digits are mathematically correct according to the official algorithm
type CPF string

// EmptyCPF represents the zero value for CPF type.
var EmptyCPF CPF

func parseCPF(input string) (CPF, error) {
	if input == "" {
		return EmptyCPF, nil
	}

	sanitized := nonDigitRegex.ReplaceAllString(input, "")

	if len(sanitized) != 11 {
		return EmptyCPF, fault.New("CPF must have 11 digits", fault.WithCode(fault.Invalid), fault.WithContext("input", input))
	}

	// Check for invalid known sequences (e.g., "11111111111")
	allSame := true
	for i := 1; i < 11; i++ {
		if sanitized[i] != sanitized[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return EmptyCPF, fault.New("invalid CPF sequence of repeated digits", fault.WithCode(fault.Invalid), fault.WithContext("input", input))
	}

	// Calculate check digits
	var d1, d2 int

	// First check digit
	sum1 := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(sanitized[i]))
		sum1 += digit * (10 - i)
	}
	remainder1 := sum1 % 11
	if remainder1 < 2 {
		d1 = 0
	} else {
		d1 = 11 - remainder1
	}

	d1Str, _ := strconv.Atoi(string(sanitized[9]))
	if d1 != d1Str {
		return EmptyCPF, fault.New("invalid CPF check digit 1", fault.WithCode(fault.Invalid), fault.WithContext("input", input))
	}

	// Second check digit
	sum2 := 0
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(sanitized[i]))
		sum2 += digit * (11 - i)
	}
	remainder2 := sum2 % 11
	if remainder2 < 2 {
		d2 = 0
	} else {
		d2 = 11 - remainder2
	}

	d2Str, _ := strconv.Atoi(string(sanitized[10]))
	if d2 != d2Str {
		return EmptyCPF, fault.New("invalid CPF check digit 2", fault.WithCode(fault.Invalid), fault.WithContext("input", input))
	}

	return CPF(sanitized), nil
}

// NewCPF creates a new CPF from the given input string.
// It accepts CPF in various formats (with or without dots and dash) and validates it.
//
// The function performs the following validations:
//   - Removes all non-digit characters
//   - Checks if it has exactly 11 digits
//   - Validates that it's not a sequence of repeated digits
//   - Verifies both check digits using the official algorithm
//
// Examples:
//   cpf, err := NewCPF("123.456.789-09")  // Valid formatted
//   cpf, err := NewCPF("12345678909")     // Valid unformatted
//   cpf, err := NewCPF("")               // Returns EmptyCPF
//   cpf, err := NewCPF("11111111111")    // Error: repeated digits
//   cpf, err := NewCPF("123456789")      // Error: too short
func NewCPF(input string) (CPF, error) {
	return parseCPF(input)
}

// String returns the CPF as a string without formatting (digits only).
// For formatted output, use Formatted() method instead.
func (c CPF) String() string {
	return string(c)
}

// IsZero returns true if the CPF is the zero value (EmptyCPF).
func (c CPF) IsZero() bool {
	return c == EmptyCPF
}

// Formatted returns the CPF in the standard Brazilian format (XXX.XXX.XXX-XX).
// If the CPF is invalid or has wrong length, returns the unformatted string.
//
// Example:
//   cpf := CPF("12345678909")
//   fmt.Println(cpf.Formatted()) // Output: "123.456.789-09"
func (c CPF) Formatted() string {
	if len(c) != 11 {
		return c.String()
	}
	return fmt.Sprintf("%s.%s.%s-%s", c[0:3], c[3:6], c[6:9], c[9:11])
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the CPF as a JSON string without formatting.
func (c CPF) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into a CPF, performing full validation.
func (c *CPF) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "CPF must be a valid JSON string", fault.WithCode(fault.Invalid))
	}
	cpf, err := NewCPF(s)
	if err != nil {
		return err
	}
	*c = cpf
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the CPF as a string or nil if zero value.
func (c CPF) Value() (driver.Value, error) {
	if c.IsZero() {
		return nil, nil
	}
	return c.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts string or []byte values and validates them as CPF.
func (c *CPF) Scan(src interface{}) error {
	if src == nil {
		*c = EmptyCPF
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for CPF", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}

	cpf, err := NewCPF(s)
	if err != nil {
		return err
	}
	*c = cpf
	return nil
}
