package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/marcelofabianov/fault"
)

// defaultPrecision is the global default number of decimal places for new quantities.
var defaultPrecision = 3

// SetDefaultPrecision sets the global default precision for the Quantity type.
// This affects new quantities created without specifying a precision.
func SetDefaultPrecision(p int) {
	if p >= 0 {
		defaultPrecision = p
	}
}

// Quantity is a value object representing a numeric amount with a specific unit of measure.
// It is designed to handle decimal values with a defined precision by storing the value as a scaled integer,
// thus avoiding floating-point inaccuracies in calculations.
//
// The unit must be registered in the global `Unit` registry before use.
//
// Example:
//   wisp.RegisterUnits("BOX")
//   q, _ := wisp.NewQuantity(12.5, "BOX")
//   price, _ := wisp.NewMoney(1000, wisp.BRL) // R$10.00 per box
//   total, _ := q.MultiplyByMoney(price) // R$125.00
type Quantity struct {
	value     int64
	unit      Unit
	precision int
}

// newQuantity is the internal constructor for creating a Quantity with a specific precision.
func newQuantity(value float64, unit Unit, precision int) (Quantity, error) {
	if !unit.IsValid() {
		return Quantity{}, fault.New(
			"unit is not registered as a valid unit of measure",
			fault.WithCode(fault.Invalid),
			fault.WithContext("unit", unit),
		)
	}

	if precision < 0 {
		return Quantity{}, fault.New(
			"precision cannot be negative",
			fault.WithCode(fault.Invalid),
			fault.WithContext("precision", precision),
		)
	}

	factor := math.Pow10(precision)
	scaledValue := int64(math.RoundToEven(value * factor))

	return Quantity{
		value:     scaledValue,
		unit:      unit,
		precision: precision,
	}, nil
}

// NewQuantity creates a new Quantity with a value and a unit, using the default precision.
// The unit must be pre-registered.
// Returns an error if the unit is not valid.
func NewQuantity(value float64, unit Unit) (Quantity, error) {
	return newQuantity(value, unit, defaultPrecision)
}

// NewQuantityWithPrecision creates a new Quantity with a specific precision.
// The unit must be pre-registered.
// Returns an error if the unit is not valid or precision is negative.
func NewQuantityWithPrecision(value float64, unit Unit, precision int) (Quantity, error) {
	return newQuantity(value, unit, precision)
}

// IntValue returns the scaled integer value of the quantity.
func (q Quantity) IntValue() int64 {
	return q.value
}

// Unit returns the unit of measure for the quantity.
func (q Quantity) Unit() Unit {
	return q.unit
}

// Precision returns the number of decimal places the quantity supports.
func (q Quantity) Precision() int {
	return q.precision
}

// IsZero returns true if the Quantity is the zero value.
func (q Quantity) IsZero() bool {
	return q.value == 0 && q.unit == ""
}

// Float64 returns the quantity's value as a float64, unscaled.
func (q Quantity) Float64() float64 {
	if q.precision == 0 {
		return float64(q.value)
	}
	return float64(q.value) / math.Pow10(q.precision)
}

// Add returns a new Quantity that is the sum of this quantity and another.
// It returns an error if the units or precisions of the two quantities are different.
func (q Quantity) Add(other Quantity) (Quantity, error) {
	if q.unit != other.unit {
		return Quantity{}, fault.New(
			"cannot add quantities with different units",
			fault.WithCode(fault.DomainViolation),
			fault.WithContext("unit_a", q.unit),
			fault.WithContext("unit_b", other.unit),
		)
	}

	if q.precision != other.precision {
		return Quantity{}, fault.New(
			"cannot add quantities with different precisions",
			fault.WithCode(fault.DomainViolation),
			fault.WithContext("precision_a", q.precision),
			fault.WithContext("precision_b", other.precision),
		)
	}

	return Quantity{
		value:     q.value + other.value,
		unit:      q.unit,
		precision: q.precision,
	}, nil
}

// MultiplyByMoney calculates the total cost by multiplying the quantity by a price per unit.
// It returns a new Money instance representing the total value.
func (q Quantity) MultiplyByMoney(pricePerUnit Money) (Money, error) {
	if pricePerUnit.IsZero() {
		return Money{amount: 0, currency: pricePerUnit.Currency()}, nil
	}

	totalValue := q.Float64() * float64(pricePerUnit.Amount())
	roundedAmount := int64(math.RoundToEven(totalValue))

	return NewMoney(roundedAmount, pricePerUnit.Currency())
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Quantity to a JSON object with its float value and unit.
func (q Quantity) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Value float64 `json:"value"`
		Unit  Unit    `json:"unit"`
	}{
		Value: q.Float64(),
		Unit:  q.unit,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object into a Quantity, automatically detecting precision from the value.
func (q *Quantity) UnmarshalJSON(data []byte) error {
	dto := &struct {
		Value float64 `json:"value"`
		Unit  Unit    `json:"unit"`
	}{}

	if err := json.Unmarshal(data, &dto); err != nil {
		return fault.Wrap(err, "invalid JSON format for Quantity", fault.WithCode(fault.Invalid))
	}

	precision := 0
	if s := strings.Split(fmt.Sprintf("%v", dto.Value), "."); len(s) == 2 {
		precision = len(s[1])
	}

	qty, err := newQuantity(dto.Value, dto.Unit, precision)
	if err != nil {
		return err
	}
	*q = qty
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the Quantity as a JSON string or nil if it's the zero value.
func (q Quantity) Value() (driver.Value, error) {
	if q.IsZero() {
		return nil, nil
	}

	data, err := q.MarshalJSON()
	if err != nil {
		return nil, fault.Wrap(err,
			"failed to marshal quantity for database storage",
			fault.WithCode(fault.Internal),
		)
	}

	return string(data), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts string or []byte values containing JSON and validates them as Quantity.
func (q *Quantity) Scan(src interface{}) error {
	if src == nil {
		*q = Quantity{}
		return nil
	}

	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fault.New(
			"unsupported scan type for Quantity",
			fault.WithCode(fault.Invalid),
			fault.WithContext("received_type", fmt.Sprintf("%T", src)),
		)
	}

	if err := q.UnmarshalJSON(data); err != nil {
		return err
	}

	return nil
}
