package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

// Latitude is a value object representing a geographic latitude.
// It ensures that the value is always within the valid range of -90 to +90 degrees.
//
// The value is stored as a float64.
//
// Example:
//   lat, err := NewLatitude(-23.55052)
type Latitude float64

// NewLatitude creates a new Latitude.
// It returns an error if the value is outside the valid range of -90 to +90.
func NewLatitude(value float64) (Latitude, error) {
	if value < -90.0 || value > 90.0 {
		return 0, fault.New(
			"latitude must be between -90 and 90",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", value),
		)
	}
	return Latitude(value), nil
}

// Float64 returns the latitude value as a float64.
func (l Latitude) Float64() float64 {
	return float64(l)
}

// String returns the string representation of the latitude value.
func (l Latitude) String() string {
	return fmt.Sprintf("%f", l)
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Latitude to its float64 representation.
func (l Latitude) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Float64())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON number into a Latitude, with validation.
func (l *Latitude) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return fault.Wrap(err, "Latitude must be a valid JSON number", fault.WithCode(fault.Invalid))
	}

	lat, err := NewLatitude(f)
	if err != nil {
		return err
	}
	*l = lat
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the Latitude as a float64.
func (l Latitude) Value() (driver.Value, error) {
	return l.Float64(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a numeric type from the database and converts it into a Latitude, with validation.
func (l *Latitude) Scan(src interface{}) error {
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
			return fault.Wrap(err, "failed to scan bytes into Latitude", fault.WithCode(fault.Invalid))
		}
	default:
		return fault.New("unsupported scan type for Latitude", fault.WithCode(fault.Invalid))
	}

	lat, err := NewLatitude(f)
	if err != nil {
		return err
	}
	*l = lat
	return nil
}
