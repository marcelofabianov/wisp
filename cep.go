package atomic

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

type CEP string

var EmptyCEP CEP

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

func NewCEP(input string) (CEP, error) {
	return parseCEP(input)
}

func (c CEP) String() string {
	return string(c)
}

func (c CEP) IsZero() bool {
	return c == EmptyCEP
}

func (c CEP) Formatted() string {
	if len(c) != 8 {
		return c.String()
	}
	return fmt.Sprintf("%s-%s", c[0:5], c[5:8])
}

func (c CEP) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

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

func (c CEP) Value() (driver.Value, error) {
	if c.IsZero() {
		return nil, nil
	}
	return c.String(), nil
}

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
