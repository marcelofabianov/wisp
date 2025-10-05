package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"

	"github.com/marcelofabianov/fault"
)

// WeightUnit defines the supported units of mass.
type WeightUnit string

// Constants for supported weight units.
const (
	Kilogram WeightUnit = "kg"
	Gram     WeightUnit = "g"
	Pound    WeightUnit = "lb"
	Ounce    WeightUnit = "oz"
)

// Conversion factors to grams.
const (
	gramsInAKilogram = 1000
	gramsInAPound    = 453.59237
	gramsInAnOunce   = 28.34952
	mgInAGram        = 1000
)

// Weight is a value object representing a physical weight.
// It stores the value internally in milligrams to maintain precision and avoid floating-point errors
// during conversions and calculations. It supports common units like kilograms, grams, pounds, and ounces.
//
// The zero value is ZeroWeight.
//
// Example:
//   w, err := NewWeight(1.5, Kilogram)
//   pounds, _ := w.In(Pound) // Converts the weight to pounds
type Weight struct {
	milligrams int64
}

// ZeroWeight represents the zero value for the Weight type.
var ZeroWeight = Weight{}

// NewWeight creates a new Weight from a float value and a unit.
// It converts the input value to milligrams for internal storage.
// Returns an error if the value is negative or the unit is not supported.
func NewWeight(value float64, unit WeightUnit) (Weight, error) {
	if value < 0 {
		return ZeroWeight, fault.New("weight value cannot be negative", fault.WithCode(fault.Invalid))
	}

	var grams float64
	switch unit {
	case Kilogram:
		grams = value * gramsInAKilogram
	case Gram:
		grams = value
	case Pound:
		grams = value * gramsInAPound
	case Ounce:
		grams = value * gramsInAnOunce
	default:
		return ZeroWeight, fault.New("unsupported weight unit", fault.WithCode(fault.Invalid), fault.WithContext("unit", unit))
	}

	mg := int64(math.Round(grams * mgInAGram))

	return Weight{milligrams: mg}, nil
}

// In converts the stored weight to the specified unit.
// It returns the value as a float64.
// Returns an error if the target unit is not supported.
func (w Weight) In(unit WeightUnit) (float64, error) {
	grams := float64(w.milligrams) / mgInAGram

	switch unit {
	case Kilogram:
		return grams / gramsInAKilogram, nil
	case Gram:
		return grams, nil
	case Pound:
		return grams / gramsInAPound, nil
	case Ounce:
		return grams / gramsInAnOunce, nil
	}

	return 0, fault.New("unsupported weight unit for conversion", fault.WithCode(fault.Invalid), fault.WithContext("unit", unit))
}

// Add returns a new Weight that is the sum of this weight and another.
func (w Weight) Add(other Weight) Weight {
	return Weight{milligrams: w.milligrams + other.milligrams}
}

// Subtract returns a new Weight that is the difference between this weight and another.
func (w Weight) Subtract(other Weight) Weight {
	return Weight{milligrams: w.milligrams - other.milligrams}
}

// IsNegative returns true if the weight is negative.
func (w Weight) IsNegative() bool {
	return w.milligrams < 0
}

// Equals checks if two Weight instances are equal.
func (w Weight) Equals(other Weight) bool {
	return w.milligrams == other.milligrams
}

// String returns the weight formatted as kilograms (e.g., "1.500 kg").
func (w Weight) String() string {
	kg, _ := w.In(Kilogram)
	return fmt.Sprintf("%.3f kg", kg)
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Weight to a JSON object with its value in kilograms.
func (w Weight) MarshalJSON() ([]byte, error) {
	kg, _ := w.In(Kilogram)
	return json.Marshal(&struct {
		Value float64    `json:"value"`
		Unit  WeightUnit `json:"unit"`
	}{
		Value: kg,
		Unit:  Kilogram,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object with a value and unit into a Weight.
func (w *Weight) UnmarshalJSON(data []byte) error {
	dto := &struct {
		Value float64    `json:"value"`
		Unit  WeightUnit `json:"unit"`
	}{}

	if err := json.Unmarshal(data, &dto); err != nil {
		return fault.Wrap(err, "invalid JSON format for Weight", fault.WithCode(fault.Invalid))
	}

	weight, err := NewWeight(dto.Value, dto.Unit)
	if err != nil {
		return err
	}
	*w = weight
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the weight in milligrams as an int64.
func (w Weight) Value() (driver.Value, error) {
	return w.milligrams, nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts an int64 (milligrams) from the database and converts it into a Weight.
func (w *Weight) Scan(src interface{}) error {
	if src == nil {
		*w = ZeroWeight
		return nil
	}

	var mg int64
	switch v := src.(type) {
	case int64:
		mg = v
	default:
		return fault.New("unsupported scan type for Weight", fault.WithCode(fault.Invalid))
	}

	if mg < 0 {
		return fault.New("weight from database cannot be negative", fault.WithCode(fault.Invalid))
	}

	*w = Weight{milligrams: mg}
	return nil
}
