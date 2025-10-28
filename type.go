package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/marcelofabianov/fault"
)

// Type represents a registered classification or category within a finite set defined by the domain.
type Type string

var validTypes = make(map[Type]struct{})
var EmptyType Type

// RegisterTypes allows the consuming application to define a set of valid type identifiers.
// Normalizes types to uppercase and trimmed.
// This function should be called during application startup.
func RegisterTypes(types ...Type) {
	for _, t := range types {
		normalized := Type(strings.ToUpper(strings.TrimSpace(string(t))))
		if normalized != "" {
			validTypes[normalized] = struct{}{}
		}
	}
}

// NewType creates a new Type from a string, ensuring it's a valid, registered type.
func NewType(value string) (Type, error) {
	normalized := Type(strings.ToUpper(strings.TrimSpace(value)))
	if normalized == EmptyType {
		// Allow creation of an empty type if input is empty/blank
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

// ClearRegisteredTypes is a helper function to clear all registered types,
// primarily useful in test environments.
func ClearRegisteredTypes() {
	validTypes = make(map[Type]struct{})
}

// String returns the string representation of the Type.
func (t Type) String() string {
	return string(t)
}

// IsValid checks if the Type has been registered as a valid type identifier.
func (t Type) IsValid() bool {
	_, ok := validTypes[t]
	return ok
}

// IsZero returns true if the type is empty.
func (t Type) IsZero() bool {
	return t == EmptyType
}

// MarshalJSON implements the json.Marshaler interface.
func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "Type must be a valid JSON string", fault.WithCode(fault.Invalid))
	}
	newT, err := NewType(s)
	if err != nil {
		return err // NewType already wraps the error correctly
	}
	*t = newT
	return nil
}

// Value implements the driver.Valuer interface.
func (t Type) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.String(), nil
}

// Scan implements the sql.Scanner interface.
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
		return err // NewType already wraps the error correctly
	}
	*t = newT
	return nil
}
