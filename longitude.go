package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

// Longitude is a value object representing a geographic longitude.
// It ensures that the value is always within the valid range of -180 to +180 degrees.
//
// The value is stored as a float64.
//
// Example:
//   lon, err := NewLongitude(-46.633308)
type Longitude float64

// NewLongitude creates a new Longitude.
// It returns an error if the value is outside the valid range of -180 to +180.
func NewLongitude(value float64) (Longitude, error) {
	if value < -180.0 || value > 180.0 {
		return 0, fault.New(
			"longitude must be between -180 and 180",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", value),
		)
	}
	return Longitude(value), nil
}

// Float64 returns the longitude value as a float64.
func (l Longitude) Float64() float64 {
	return float64(l)
}

// String returns the string representation of the longitude value.
func (l Longitude) String() string {
	return fmt.Sprintf("%f", l)
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Longitude to its float64 representation.
func (l Longitude) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Float64())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON number into a Longitude, with validation.
func (l *Longitude) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return fault.Wrap(err, "Longitude must be a valid JSON number", fault.WithCode(fault.Invalid))
	}

	lon, err := NewLongitude(f)
	if err != nil {
		return err
	}
	*l = lon
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the Longitude as a float64.
func (l Longitude) Value() (driver.Value, error) {
	return l.Float64(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a numeric type from the database and converts it into a Longitude, with validation.
func (l *Longitude) Scan(src interface{}) error {
	if src == nil {
		*l = 0
		return nil
	}

	var f float64
	switch v := src.(type) {
	case float64:
		f = v
	case float32:
		f = float64(v)
	case int64:
		f = float64(v)
	case []byte:
		if _, err := fmt.Sscanf(string(v), "%f", &f); err != nil {
			return fault.Wrap(err, "failed to scan bytes into Longitude", fault.WithCode(fault.Invalid))
		}
	default:
		return fault.New("unsupported scan type for Longitude", fault.WithCode(fault.Invalid))
	}

	lon, err := NewLongitude(f)
	if err != nil {
		return err
	}
	*l = lon
	return nil
}
