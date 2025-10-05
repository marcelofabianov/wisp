package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/marcelofabianov/fault"
)

// registeredMIMETypes holds the global set of allowed MIME types.
var registeredMIMETypes = make(map[MIMEType]struct{})

// MIMEType is a value object representing a standard MIME type (e.g., "application/json", "image/jpeg").
// It ensures that only explicitly registered MIME types are used, which is crucial for security
// and for controlling the types of content processed by an application.
//
// Before use, MIME types must be added to a global registry via `RegisterMIMETypes`.
// The value is stored in a normalized (lowercase) "type/subtype" format.
//
// Example:
//   wisp.RegisterMIMETypes("application/json", "image/jpeg")
//   mt, err := wisp.NewMIMEType("application/json")
type MIMEType string

// EmptyMIMEType represents the zero value for the MIMEType type.
var EmptyMIMEType MIMEType

// RegisterMIMETypes adds one or more MIME types to the global registry.
// It normalizes them to lowercase and validates the "type/subtype" format.
// This function should be called at application startup.
func RegisterMIMETypes(types ...string) {
	for _, t := range types {
		normalized := strings.ToLower(strings.TrimSpace(t))
		if normalized != "" {
			parts := strings.Split(normalized, "/")
			if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
				registeredMIMETypes[MIMEType(normalized)] = struct{}{}
			}
		}
	}
}

// ClearRegisteredMIMETypes removes all MIME types from the global registry.
// This is primarily for testing purposes.
func ClearRegisteredMIMETypes() {
	registeredMIMETypes = make(map[MIMEType]struct{})
}

// NewMIMEType creates a new MIMEType from a string.
// It normalizes the input and validates it against the "type/subtype" format and the global registry.
// Returns an error if the input is empty, malformed, or not registered.
func NewMIMEType(input string) (MIMEType, error) {
	normalized := strings.ToLower(strings.TrimSpace(input))

	if normalized == "" {
		return EmptyMIMEType, fault.New("mime type input cannot be empty", fault.WithCode(fault.Invalid))
	}

	parts := strings.Split(normalized, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return EmptyMIMEType, fault.New(
			"mime type must follow the 'type/subtype' format",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", input),
		)
	}

	mt := MIMEType(normalized)
	if !mt.IsRegistered() {
		return EmptyMIMEType, fault.New(
			"mime type is not registered in the allowed list",
			fault.WithCode(fault.Invalid),
			fault.WithContext("mime_type", normalized),
		)
	}

	return mt, nil
}

// IsRegistered checks if the MIMEType is in the global registry.
func (mt MIMEType) IsRegistered() bool {
	_, ok := registeredMIMETypes[mt]
	return ok
}

// Type returns the primary type part of the MIME type (e.g., "application").
func (mt MIMEType) Type() string {
	if mt.IsZero() {
		return ""
	}
	parts := strings.Split(string(mt), "/")
	return parts[0]
}

// SubType returns the subtype part of the MIME type (e.g., "json").
func (mt MIMEType) SubType() string {
	if mt.IsZero() {
		return ""
	}
	parts := strings.Split(string(mt), "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

// String returns the full MIME type string in "type/subtype" format.
func (mt MIMEType) String() string {
	return string(mt)
}

// IsZero returns true if the MIMEType is the zero value.
func (mt MIMEType) IsZero() bool {
	return mt == EmptyMIMEType
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the MIMEType to its string representation.
func (mt MIMEType) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into a MIMEType, with validation.
func (mt *MIMEType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "MIMEType must be a valid JSON string", fault.WithCode(fault.Invalid))
	}
	newMT, err := NewMIMEType(s)
	if err != nil {
		return err
	}
	*mt = newMT
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the MIMEType as a string.
func (mt MIMEType) Value() (driver.Value, error) {
	if mt.IsZero() {
		return nil, nil
	}
	return mt.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a string or byte slice from the database and converts it into a MIMEType, with validation.
func (mt *MIMEType) Scan(src interface{}) error {
	if src == nil {
		*mt = EmptyMIMEType
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for MIMEType", fault.WithCode(fault.Invalid))
	}

	newMT, err := NewMIMEType(s)
	if err != nil {
		return err
	}
	*mt = newMT
	return nil
}
