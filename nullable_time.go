package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marcelofabianov/fault"
)

// NullableTime is a wrapper for time.Time that can be null.
// This is useful for optional timestamps in a database, such as `deleted_at` or `processed_at`,
// where the absence of a time is meaningful.
// It implements the necessary interfaces for JSON and database serialization/deserialization.
//
// The zero value is an invalid NullableTime (Valid=false).
//
// Example:
//   var deletedAt wisp.NullableTime
//   if shouldDelete {
//       deletedAt = wisp.NewNullableTime(time.Now())
//   }
type NullableTime struct {
	Time  time.Time
	Valid bool
}

// EmptyNullableTime represents the zero value for NullableTime, which is an invalid (null) time.
var EmptyNullableTime = NullableTime{}

// NewNullableTime creates a new NullableTime.
// If the provided time.Time is zero, the NullableTime is considered invalid (null).
func NewNullableTime(t time.Time) NullableTime {
	return NullableTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

// IsZero returns true if the NullableTime is invalid (null).
// This is an alias for !nt.Valid to provide a consistent IsZero interface.
func (nt NullableTime) IsZero() bool {
	return !nt.Valid
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the NullableTime to a JSON time string, or `null` if it is invalid.
func (nt NullableTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON time string or `null` into a NullableTime.
func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}

	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return fault.Wrap(err, "NullableTime must be a valid JSON time string or null", fault.WithCode(fault.Invalid))
	}

	nt.Time = t
	nt.Valid = true
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the time.Time value, or `nil` if the NullableTime is invalid.
func (nt NullableTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a time.Time or `nil` from the database and converts it into a NullableTime.
func (nt *NullableTime) Scan(src interface{}) error {
	if src == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}

	switch v := src.(type) {
	case time.Time:
		nt.Time, nt.Valid = v, true
		return nil
	default:
		return fault.New("unsupported scan type for NullableTime", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}
}
