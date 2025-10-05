package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/marcelofabianov/fault"
)

// Version is a value object representing an entity's version number, used for optimistic locking.
// It is a simple integer that increments with each modification to an entity.
// This helps prevent lost updates in concurrent environments.
//
// The zero value is ZeroVersion, and the typical starting version is 1.
//
// Example:
//   v := wisp.InitialVersion() // 1
//   v = v.Increment() // 2
type Version int

// ZeroVersion represents the zero value for the Version type.
var ZeroVersion Version

// NewVersion creates a new Version.
// It returns an error if the provided integer is negative.
func NewVersion(v int) (Version, error) {
	if v < 0 {
		return ZeroVersion, fault.New(
			"version cannot be negative",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", v),
		)
	}
	return Version(v), nil
}

// InitialVersion returns the standard starting version for a new entity, which is 1.
func InitialVersion() Version {
	return Version(1)
}

// Increment returns a new Version that is one greater than the current one.
func (v Version) Increment() Version {
	return v + 1
}

// Previous returns the preceding version.
// It will not go below zero.
func (v Version) Previous() Version {
	if v <= 1 {
		return ZeroVersion
	}
	return v - 1
}

// IsZero returns true if the Version is the zero value.
func (v Version) IsZero() bool {
	return v == ZeroVersion
}

// Equals checks if two Version instances are equal.
func (v Version) Equals(other Version) bool {
	return v == other
}

// IsGreaterThan checks if this version is greater than another.
func (v Version) IsGreaterThan(other Version) bool {
	return v > other
}

// IsLessThan checks if this version is less than another.
func (v Version) IsLessThan(other Version) bool {
	return v < other
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Version as a JSON number.
func (v Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(v))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON number into a Version, with validation.
func (v *Version) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*v = ZeroVersion
		return nil
	}

	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return fault.Wrap(err,
			"version must be a valid JSON number",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_json", string(data)),
		)
	}

	if i < 0 {
		return fault.New(
			"version cannot be negative",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", i),
		)
	}

	*v = Version(i)
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the Version as an int64.
func (v Version) Value() (driver.Value, error) {
	return int64(v), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a numeric type from the database and converts it into a Version.
func (v *Version) Scan(src interface{}) error {
	if src == nil {
		*v = ZeroVersion
		return nil
	}

	var intVal int64
	switch s := src.(type) {
	case int64:
		intVal = s
	case []byte:
		parsed, err := strconv.ParseInt(string(s), 10, 64)
		if err != nil {
			return fault.Wrap(err,
				"failed to convert bytes to version number",
				fault.WithCode(fault.Invalid),
				fault.WithContext("input_bytes", string(s)),
			)
		}
		intVal = parsed
	default:
		return fault.New(
			"incompatible type for Version scan",
			fault.WithCode(fault.Invalid),
			fault.WithContext("received_type", fmt.Sprintf("%T", src)),
		)
	}

	if intVal < 0 {
		return fault.New(
			"version from database cannot be negative",
			fault.WithCode(fault.Invalid),
			fault.WithContext("source_value", intVal),
		)
	}

	*v = Version(intVal)
	return nil
}

// Int returns the version number as a standard int.
func (v Version) Int() int {
	return int(v)
}
