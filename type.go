package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/marcelofabianov/fault"
)

type Type string

var validTypes = make(map[Type]struct{})
var EmptyType Type

func RegisterTypes(types ...Type) {
	for _, t := range types {
		normalized := Type(strings.ToLower(strings.TrimSpace(string(t))))
		if normalized != "" {
			validTypes[normalized] = struct{}{}
		}
	}
}

func NewType(value string) (Type, error) {
	normalized := Type(strings.ToLower(strings.TrimSpace(value)))
	if normalized == EmptyType {
		return EmptyType, nil
	}

	if !normalized.IsValid() {
		return EmptyType, fault.New(
			"identifier is not registered as a valid type",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_type", value),
		)
	}
	return normalized, nil
}

func ClearRegisteredTypes() {
	validTypes = make(map[Type]struct{})
}

func (t Type) String() string {
	return string(t)
}

func (t Type) IsValid() bool {
	_, ok := validTypes[t]
	return ok
}

func (t Type) IsZero() bool {
	return t == EmptyType
}

func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "Type must be a valid JSON string", fault.WithCode(fault.Invalid))
	}
	newT, err := NewType(s)
	if err != nil {
		return err
	}
	*t = newT
	return nil
}

func (t Type) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.String(), nil
}

func (t *Type) Scan(src interface{}) error {
	if src == nil {
		*t = EmptyType
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for Type", fault.WithCode(fault.Invalid))
	}

	newT, err := NewType(s)
	if err != nil {
		return err
	}
	*t = newT
	return nil
}
