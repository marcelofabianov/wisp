package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/marcelofabianov/fault"
)

type CNPJ string

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

func NewCNPJ(input string) (CNPJ, error) {
	return parseCNPJ(input)
}

func (c CNPJ) String() string {
	return string(c)
}

func (c CNPJ) IsZero() bool {
	return c == EmptyCNPJ
}

func (c CNPJ) Formatted() string {
	if len(c) != 14 {
		return c.String()
	}
	return fmt.Sprintf("%s.%s.%s/%s-%s", c[0:2], c[2:5], c[5:8], c[8:12], c[12:14])
}

func (c CNPJ) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

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

func (c CNPJ) Value() (driver.Value, error) {
	if c.IsZero() {
		return nil, nil
	}
	return c.String(), nil
}

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
