package atomic

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/marcelofabianov/fault"
)

type CPF string

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

func NewCPF(input string) (CPF, error) {
	return parseCPF(input)
}

func (c CPF) String() string {
	return string(c)
}

func (c CPF) IsZero() bool {
	return c == EmptyCPF
}

func (c CPF) Formatted() string {
	if len(c) != 11 {
		return c.String()
	}
	return fmt.Sprintf("%s.%s.%s-%s", c[0:3], c[3:6], c[6:9], c[9:11])
}

func (c CPF) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

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

func (c CPF) Value() (driver.Value, error) {
	if c.IsZero() {
		return nil, nil
	}
	return c.String(), nil
}

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
