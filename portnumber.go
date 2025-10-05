package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/marcelofabianov/fault"
)

// maxPortNumber is the highest valid TCP/UDP port number.
const maxPortNumber = 65535

// PortNumber is a value object representing a network port number.
// It ensures that the port is within the valid range for TCP/UDP (1-65535).
//
// The zero value is ZeroPortNumber.
//
// Example:
//   p, err := NewPortNumber(8080)
type PortNumber uint16

// ZeroPortNumber represents the zero value for the PortNumber type.
var ZeroPortNumber PortNumber

// NewPortNumber creates a new PortNumber.
// It returns an error if the value is not within the valid port range (1-65535).
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

// Uint16 returns the port number as a uint16.
func (p PortNumber) Uint16() uint16 {
	return uint16(p)
}

// IsZero returns true if the PortNumber is the zero value.
func (p PortNumber) IsZero() bool {
	return p == ZeroPortNumber
}

// String returns the string representation of the port number.
func (p PortNumber) String() string {
	return strconv.Itoa(int(p))
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the PortNumber as a JSON number.
func (p PortNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Uint16())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON number into a PortNumber, with validation.
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

// Value implements the driver.Valuer interface for database storage.
// It returns the PortNumber as an int64.
func (p PortNumber) Value() (driver.Value, error) {
	return int64(p.Uint16()), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts an int64 from the database and converts it into a PortNumber, with validation.
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
