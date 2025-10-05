package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/marcelofabianov/fault"
)

// CNPJ represents a Brazilian legal entity identification number (Cadastro Nacional da Pessoa Jurídica).
// It validates the format and verifies check digits according to Brazilian government standards.
// The value is stored without formatting (digits only) but can be displayed with proper formatting.
//
// Examples:
//   - Input: "12.345.678/0001-90" or "12345678000190"
//   - Storage: "12345678000190"
//   - Formatted output: "12.345.678/0001-90"
//
// A CNPJ is considered valid when:
//   - It contains exactly 14 digits
//   - It's not a sequence of repeated digits (e.g., "11111111111111")
//   - Both check digits are mathematically correct according to the official algorithm
type CNPJ string

// EmptyCNPJ represents the zero value for CNPJ type.
var EmptyCNPJ CNPJ

func parseCNPJ(input string) (CNPJ, error) {
	if input == "" {
		return EmptyCNPJ, nil
	}

	sanitized := nonDigitRegex.ReplaceAllString(input, "")

	if len(sanitized) != 14 {
		return EmptyCNPJ, fault.New("CNPJ must have 14 digits", fault.WithCode(fault.Invalid), fault.WithContext("input", input))
	}

	// Check for invalid known sequences (e.g., "11111111111111")
	allSame := true
	for i := 1; i < 14; i++ {
		if sanitized[i] != sanitized[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return EmptyCNPJ, fault.New("invalid CNPJ sequence of repeated digits", fault.WithCode(fault.Invalid), fault.WithContext("input", input))
	}

	// Calculate check digits
	var d1, d2 int
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	// First check digit
	sum1 := 0
	for i := 0; i < 12; i++ {
		digit, _ := strconv.Atoi(string(sanitized[i]))
		sum1 += digit * weights1[i]
	}
	remainder1 := sum1 % 11
	if remainder1 < 2 {
		d1 = 0
	} else {
		d1 = 11 - remainder1
	}

	d1Str, _ := strconv.Atoi(string(sanitized[12]))
	if d1 != d1Str {
		return EmptyCNPJ, fault.New("invalid CNPJ check digit 1", fault.WithCode(fault.Invalid), fault.WithContext("input", input))
	}

	// Second check digit
	sum2 := 0
	for i := 0; i < 13; i++ {
		digit, _ := strconv.Atoi(string(sanitized[i]))
		sum2 += digit * weights2[i]
	}
	remainder2 := sum2 % 11
	if remainder2 < 2 {
		d2 = 0
	} else {
		d2 = 11 - remainder2
	}

	d2Str, _ := strconv.Atoi(string(sanitized[13]))
	if d2 != d2Str {
		return EmptyCNPJ, fault.New("invalid CNPJ check digit 2", fault.WithCode(fault.Invalid), fault.WithContext("input", input))
	}

	return CNPJ(sanitized), nil
}

// NewCNPJ creates a new CNPJ from the given input string.
// It accepts CNPJ in various formats (with or without dots, slash and dash) and validates it.
//
// The function performs the following validations:
//   - Removes all non-digit characters
//   - Checks if it has exactly 14 digits
//   - Validates that it's not a sequence of repeated digits
//   - Verifies both check digits using the official algorithm
//
// Examples:
//   cnpj, err := NewCNPJ("12.345.678/0001-90")  // Valid formatted
//   cnpj, err := NewCNPJ("12345678000190")      // Valid unformatted
//   cnpj, err := NewCNPJ("")                   // Returns EmptyCNPJ
//   cnpj, err := NewCNPJ("11111111111111")     // Error: repeated digits
//   cnpj, err := NewCNPJ("123456789")          // Error: too short
func NewCNPJ(input string) (CNPJ, error) {
	return parseCNPJ(input)
}

// String returns the CNPJ as a string without formatting (digits only).
// For formatted output, use Formatted() method instead.
func (c CNPJ) String() string {
	return string(c)
}

// IsZero returns true if the CNPJ is the zero value (EmptyCNPJ).
func (c CNPJ) IsZero() bool {
	return c == EmptyCNPJ
}

// Formatted returns the CNPJ in the standard Brazilian format (XX.XXX.XXX/XXXX-XX).
// If the CNPJ is invalid or has wrong length, returns the unformatted string.
//
// Example:
//   cnpj := CNPJ("12345678000190")
//   fmt.Println(cnpj.Formatted()) // Output: "12.345.678/0001-90"
func (c CNPJ) Formatted() string {
	if len(c) != 14 {
		return c.String()
	}
	return fmt.Sprintf("%s.%s.%s/%s-%s", c[0:2], c[2:5], c[5:8], c[8:12], c[12:14])
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the CNPJ as a JSON string without formatting.
func (c CNPJ) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into a CNPJ, performing full validation.
func (c *CNPJ) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "CNPJ must be a valid JSON string", fault.WithCode(fault.Invalid))
	}
	cnpj, err := NewCNPJ(s)
	if err != nil {
		return err
	}
	*c = cnpj
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the CNPJ as a string or nil if zero value.
func (c CNPJ) Value() (driver.Value, error) {
	if c.IsZero() {
		return nil, nil
	}
	return c.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts string or []byte values and validates them as CNPJ.
func (c *CNPJ) Scan(src interface{}) error {
	if src == nil {
		*c = EmptyCNPJ
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for CNPJ", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}

	cnpj, err := NewCNPJ(s)
	if err != nil {
		return err
	}
	*c = cnpj
	return nil
}
