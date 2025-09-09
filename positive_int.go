package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

// PositiveInt is a value object ensuring an integer is always greater than zero.
type PositiveInt int

var ZeroPositiveInt PositiveInt

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

func (p PositiveInt) Int() int {
	return int(p)
}

func (p PositiveInt) IsZero() bool {
	return p == ZeroPositiveInt
}

func (p PositiveInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Int())
}

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

func (p PositiveInt) Value() (driver.Value, error) {
	return int64(p.Int()), nil
}

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
