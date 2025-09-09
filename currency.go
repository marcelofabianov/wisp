package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/marcelofabianov/fault"
)

type Currency string

const (
	BRL Currency = "BRL"
	USD Currency = "USD"
	EUR Currency = "EUR"
)

var EmptyCurrency Currency

var validCurrencies = map[Currency]struct{}{
	BRL: {},
	USD: {},
	EUR: {},
}

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

func (c Currency) String() string {
	return string(c)
}

func (c Currency) IsValid() bool {
	_, ok := validCurrencies[c]
	return ok
}

func (c Currency) IsZero() bool {
	return c == EmptyCurrency
}

func (c Currency) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

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

func (c Currency) Value() (driver.Value, error) {
	if c.IsZero() {
		return nil, nil
	}
	return c.String(), nil
}

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
