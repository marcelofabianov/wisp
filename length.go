package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"

	"github.com/marcelofabianov/fault"
)

// LengthUnit defines the supported units of length.
type LengthUnit string

// Constants for supported length units.
const (
	Meter      LengthUnit = "m"
	Centimeter LengthUnit = "cm"
	Millimeter LengthUnit = "mm"
	Kilometer  LengthUnit = "km"
	Inch       LengthUnit = "in"
	Foot       LengthUnit = "ft"
)

// Conversion factors to meters.
const (
	micrometersInAMeter = 1000000.0
	metersInAKilometer  = 1000.0
	metersInAFeoot      = 0.3048
	metersInAnInch      = 0.0254
)

// Length is a value object representing a physical length.
// It stores the value internally in micrometers to maintain precision and avoid floating-point errors
// during conversions and calculations. It supports common metric and imperial units.
//
// The zero value is ZeroLength.
//
// Example:
//   l, err := NewLength(1.8, Meter)
//   feet, _ := l.In(Foot) // Converts the length to feet
type Length struct {
	micrometers int64
}

// ZeroLength represents the zero value for the Length type.
var ZeroLength = Length{}

// NewLength creates a new Length from a float value and a unit.
// It converts the input value to micrometers for internal storage.
// Returns an error if the value is negative or the unit is not supported.
func NewLength(value float64, unit LengthUnit) (Length, error) {
	if value < 0 {
		return ZeroLength, fault.New("length value cannot be negative", fault.WithCode(fault.Invalid))
	}

	var meters float64
	switch unit {
	case Meter:
		meters = value
	case Centimeter:
		meters = value / 100.0
	case Millimeter:
		meters = value / 1000.0
	case Kilometer:
		meters = value * metersInAKilometer
	case Inch:
		meters = value * metersInAnInch
	case Foot:
		meters = value * metersInAFeoot
	default:
		return ZeroLength, fault.New("unsupported length unit", fault.WithCode(fault.Invalid), fault.WithContext("unit", unit))
	}

	micrometers := int64(math.Round(meters * micrometersInAMeter))

	return Length{micrometers: micrometers}, nil
}

// In converts the stored length to the specified unit.
// It returns the value as a float64.
// Returns an error if the target unit is not supported.
func (l Length) In(unit LengthUnit) (float64, error) {
	meters := float64(l.micrometers) / micrometersInAMeter

	switch unit {
	case Meter:
		return meters, nil
	case Centimeter:
		return meters * 100.0, nil
	case Millimeter:
		return meters * 1000.0, nil
	case Kilometer:
		return meters / metersInAKilometer, nil
	case Inch:
		return meters / metersInAnInch, nil
	case Foot:
		return meters / metersInAFeoot, nil
	}

	return 0, fault.New("unsupported length unit for conversion", fault.WithCode(fault.Invalid), fault.WithContext("unit", unit))
}

// Add returns a new Length that is the sum of this length and another.
func (l Length) Add(other Length) Length {
	return Length{micrometers: l.micrometers + other.micrometers}
}

// Subtract returns a new Length that is the difference between this length and another.
func (l Length) Subtract(other Length) Length {
	return Length{micrometers: l.micrometers - other.micrometers}
}

// IsNegative returns true if the length is negative.
func (l Length) IsNegative() bool {
	return l.micrometers < 0
}

// Equals checks if two Length instances are equal.
func (l Length) Equals(other Length) bool {
	return l.micrometers == other.micrometers
}

// String returns the length formatted as meters (e.g., "1.800 m").
func (l Length) String() string {
	m, _ := l.In(Meter)
	return fmt.Sprintf("%.3f m", m)
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Length to a JSON object with its value in meters.
func (l Length) MarshalJSON() ([]byte, error) {
	m, _ := l.In(Meter)
	return json.Marshal(&struct {
		Value float64    `json:"value"`
		Unit  LengthUnit `json:"unit"`
	}{
		Value: m,
		Unit:  Meter,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object with a value and unit into a Length.
func (l *Length) UnmarshalJSON(data []byte) error {
	dto := &struct {
		Value float64    `json:"value"`
		Unit  LengthUnit `json:"unit"`
	}{}

	if err := json.Unmarshal(data, &dto); err != nil {
		return fault.Wrap(err, "invalid JSON format for Length", fault.WithCode(fault.Invalid))
	}

	length, err := NewLength(dto.Value, dto.Unit)
	if err != nil {
		return err
	}
	*l = length
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the length in micrometers as an int64.
func (l Length) Value() (driver.Value, error) {
	return l.micrometers, nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts an int64 (micrometers) from the database and converts it into a Length.
func (l *Length) Scan(src interface{}) error {
	if src == nil {
		*l = ZeroLength
		return nil
	}

	var micrometers int64
	switch v := src.(type) {
	case int64:
		micrometers = v
	default:
		return fault.New("unsupported scan type for Length", fault.WithCode(fault.Invalid))
	}

	if micrometers < 0 {
		return fault.New("length from database cannot be negative", fault.WithCode(fault.Invalid))
	}

	*l = Length{micrometers: micrometers}
	return nil
}
