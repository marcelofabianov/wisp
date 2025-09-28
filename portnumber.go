package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/marcelofabianov/fault"
)

const maxPortNumber = 65535

type PortNumber uint16

var ZeroPortNumber PortNumber

func NewPortNumber(value int) (PortNumber, error) {
	if value <= 0 || value > maxPortNumber {
		return ZeroPortNumber, fault.New(
			"port number must be between 1 and 65535",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", value),
		)
	}
	return PortNumber(value), nil
}

func (p PortNumber) Uint16() uint16 {
	return uint16(p)
}

func (p PortNumber) IsZero() bool {
	return p == ZeroPortNumber
}

func (p PortNumber) String() string {
	return strconv.Itoa(int(p))
}

func (p PortNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Uint16())
}

func (p *PortNumber) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return fault.Wrap(err, "PortNumber must be a valid JSON number", fault.WithCode(fault.Invalid))
	}

	port, err := NewPortNumber(i)
	if err != nil {
		return err
	}
	*p = port
	return nil
}

func (p PortNumber) Value() (driver.Value, error) {
	return int64(p.Uint16()), nil
}

func (p *PortNumber) Scan(src interface{}) error {
	if src == nil {
		*p = ZeroPortNumber
		return nil
	}

	var i int64
	switch v := src.(type) {
	case int64:
		i = v
	default:
		return fault.New("unsupported scan type for PortNumber", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}

	port, err := NewPortNumber(int(i))
	if err != nil {
		return err
	}
	*p = port
	return nil
}
