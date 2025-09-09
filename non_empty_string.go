package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/marcelofabianov/fault"
)

// NonEmptyString is a value object ensuring a string is not empty after trimming whitespace.
type NonEmptyString string

var EmptyNonEmptyString NonEmptyString

func NewNonEmptyString(value string) (NonEmptyString, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return EmptyNonEmptyString, fault.New(
			"string cannot be empty",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", value),
		)
	}
	return NonEmptyString(trimmed), nil
}

func (s NonEmptyString) String() string {
	return string(s)
}

func (s NonEmptyString) IsZero() bool {
	return s == EmptyNonEmptyString
}

func (s NonEmptyString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *NonEmptyString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fault.Wrap(err, "NonEmptyString must be a valid JSON string", fault.WithCode(fault.Invalid))
	}

	nes, err := NewNonEmptyString(str)
	if err != nil {
		return err
	}
	*s = nes
	return nil
}

func (s NonEmptyString) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s *NonEmptyString) Scan(src interface{}) error {
	if src == nil {
		*s = EmptyNonEmptyString
		return nil
	}

	var str string
	switch v := src.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fault.New("unsupported scan type for NonEmptyString", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}

	nes, err := NewNonEmptyString(str)
	if err != nil {
		return err
	}
	*s = nes
	return nil
}
