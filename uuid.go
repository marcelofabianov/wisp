package wisp

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/marcelofabianov/fault"
)

// UUID is a wrapper around google/uuid.UUID that provides additional validation and integration.
// It uses UUID version 7 by default, which is time-ordered and suitable for database keys.
// UUID v7 provides better database performance compared to v4 due to its sequential nature.
//
// The UUID type implements interfaces for JSON marshaling, database storage, and text encoding.
// It ensures type safety by preventing direct instantiation with invalid values.
//
// Example:
//   id, err := NewUUID()                    // Generate new v7 UUID
//   parsed, err := ParseUUID("123e4567-...") // Parse from string
//   fmt.Println(id.String())                // Output: formatted UUID string
type UUID uuid.UUID

// Nil represents the zero value for UUID type (all bits zero).
var Nil UUID

// NewUUID generates a new UUID version 7 (time-ordered).
// UUID v7 is preferred over v4 for database primary keys because it's time-ordered,
// which provides better database performance for indexes and reduces fragmentation.
//
// Returns an error if the system cannot generate a UUID (e.g., insufficient entropy).
//
// Example:
//   id, err := NewUUID()
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(id.String()) // Output: "01234567-89ab-7def-8123-456789abcdef"
func NewUUID() (UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Nil, fault.Wrap(err,
			"failed to generate v7 UUID",
			fault.WithCode(fault.Internal),
			fault.WithContext("operation", "wisp.NewUUID"),
		)
	}
	return UUID(id), nil
}

// MustNewUUID is like NewUUID but panics if it cannot generate a UUID.
// This is useful in initialization code where failure should be fatal.
//
// Example:
//   var defaultID = MustNewUUID() // Will panic if generation fails
func MustNewUUID() UUID {
	id, err := NewUUID()
	if err != nil {
		panic(err)
	}
	return id
}

// ParseUUID parses a UUID from its string representation.
// It accepts the standard UUID format: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// Both uppercase and lowercase letters are accepted.
//
// Returns an error if the input string is not a valid UUID format.
//
// Example:
//   id, err := ParseUUID("123e4567-e89b-12d3-a456-426614174000")
//   if err != nil {
//       log.Printf("Invalid UUID: %v", err)
//   }
func ParseUUID(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return Nil, fault.Wrap(err,
			"failed to parse UUID string",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input", s),
			fault.WithContext("operation", "wisp.ParseUUID"),
		)
	}
	return UUID(id), nil
}

// MustParseUUID is like ParseUUID but panics if the string is not a valid UUID.
// This is useful when you need to parse a UUID that should always be valid.
//
// Example:
//   id := MustParseUUID("123e4567-e89b-12d3-a456-426614174000")
func MustParseUUID(s string) UUID {
	id, err := ParseUUID(s)
	if err != nil {
		panic(err)
	}
	return id
}

// String returns the string representation of the UUID in canonical format.
// The format is: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" (36 characters with hyphens).
// Returns "00000000-0000-0000-0000-000000000000" for Nil UUID.
func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// IsNil reports whether the UUID is the Nil UUID (all bits zero).
// This is equivalent to checking if u == Nil but more readable.
func (u UUID) IsNil() bool {
	return u == Nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// It returns the UUID in canonical string format as bytes.
func (u UUID) MarshalText() ([]byte, error) {
	return uuid.UUID(u).MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses the UUID from text in canonical format.
func (u *UUID) UnmarshalText(text []byte) error {
	var underlyingUUID uuid.UUID
	if err := underlyingUUID.UnmarshalText(text); err != nil {
		return fault.Wrap(err,
			"invalid text representation for UUID",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_text", string(text)),
			fault.WithContext("operation", "wisp.UUID.UnmarshalText"),
		)
	}
	*u = UUID(underlyingUUID)
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the UUID as a string or nil if it's the Nil UUID.
// The database will store the UUID in canonical string format.
func (u UUID) Value() (driver.Value, error) {
	if u == Nil {
		return nil, nil
	}
	return u.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts various types that can be converted to UUID (string, []byte, etc.).
// The source value is parsed and validated as a proper UUID.
func (u *UUID) Scan(src interface{}) error {
	var underlyingUUID uuid.UUID
	if err := underlyingUUID.Scan(src); err != nil {
		return fault.Wrap(err,
			"failed to scan database value into UUID",
			fault.WithCode(fault.Invalid),
			fault.WithContext("source_type", fmt.Sprintf("%T", src)),
			fault.WithContext("operation", "wisp.UUID.Scan"),
		)
	}
	*u = UUID(underlyingUUID)
	return nil
}
