package wisp

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/marcelofabianov/fault"
)

// NullableUUID represents a UUID that can be null/nil.
// It's similar to sql.NullString but specifically for UUIDs, providing type safety.
// This is particularly useful for optional foreign key relationships in databases.
//
// The Valid field indicates whether the UUID is valid (not null).
// When Valid is false, the UUID field should be ignored.
//
// Examples:
//   - Valid UUID: NullableUUID{UUID: someUUID, Valid: true}
//   - Null UUID: NullableUUID{Valid: false} or zero value
//   - JSON: {"user_id": "123e4567-..."} or {"user_id": null}
type NullableUUID struct {
	UUID  UUID // The UUID value (only meaningful when Valid is true)
	Valid bool // Whether the UUID is valid (not null)
}

// NewNullableUUID creates a new NullableUUID from the given UUID.
// The Valid field is automatically set based on whether the UUID is Nil.
//
// Examples:
//   id, _ := NewUUID()
//   nullable := NewNullableUUID(id)        // Valid: true
//   nilNullable := NewNullableUUID(Nil)    // Valid: false
func NewNullableUUID(id UUID) NullableUUID {
	return NullableUUID{
		UUID:  id,
		Valid: !id.IsNil(),
	}
}

// IsZero returns true if the NullableUUID is not valid (represents null).
// This follows the same pattern as other nullable types in the wisp package.
func (nu NullableUUID) IsZero() bool {
	return !nu.Valid
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the UUID as a JSON string if valid, or as null if not valid.
//
// Examples:
//   - Valid UUID: "123e4567-e89b-12d3-a456-426614174000"
//   - Invalid UUID: null
func (nu NullableUUID) MarshalJSON() ([]byte, error) {
	if !nu.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nu.UUID)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes from JSON, accepting either a valid UUID string or null.
// When the JSON value is null, Valid is set to false.
// When the JSON value is a string, it's parsed as a UUID and Valid is set to true.
func (nu *NullableUUID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nu.Valid = false
		return nil
	}

	var u UUID
	if err := json.Unmarshal(data, &u); err != nil {
		return fault.Wrap(err, "NullableUUID must be a valid JSON UUID string or null", fault.WithCode(fault.Invalid))
	}

	nu.UUID = u
	nu.Valid = true
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns nil if not valid, otherwise delegates to the underlying UUID's Value method.
func (nu NullableUUID) Value() (driver.Value, error) {
	if !nu.Valid {
		return nil, nil
	}
	return nu.UUID.Value()
}

// Scan implements the sql.Scanner interface for database retrieval.
// It handles NULL values by setting Valid to false, and non-NULL values
// by delegating to the underlying UUID's Scan method.
func (nu *NullableUUID) Scan(src interface{}) error {
	if src == nil {
		nu.UUID, nu.Valid = Nil, false
		return nil
	}

	var u UUID
	if err := u.Scan(src); err != nil {
		return err
	}

	nu.UUID = u
	nu.Valid = true
	return nil
}
