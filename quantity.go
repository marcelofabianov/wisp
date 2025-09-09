package atomic

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/marcelofabianov/fault"
)

var defaultPrecision = 3

func SetDefaultPrecision(p int) {
	if p >= 0 {
		defaultPrecision = p
	}
}

type Quantity struct {
	value     int64
	unit      Unit
	precision int
}

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

func NewQuantity(value float64, unit Unit) (Quantity, error) {
	return newQuantity(value, unit, defaultPrecision)
}

func NewQuantityWithPrecision(value float64, unit Unit, precision int) (Quantity, error) {
	return newQuantity(value, unit, precision)
}

func (q Quantity) Value() int64 {
	return q.value
}

func (q Quantity) Unit() Unit {
	return q.unit
}

func (q Quantity) Precision() int {
	return q.precision
}

func (q Quantity) IsZero() bool {
	return q.value == 0 && q.unit == ""
}

func (q Quantity) Float64() float64 {
	if q.precision == 0 {
		return float64(q.value)
	}
	return float64(q.value) / math.Pow10(q.precision)
}

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

func (q Quantity) MultiplyByMoney(pricePerUnit Money) (Money, error) {
	if pricePerUnit.IsZero() {
		return Money{amount: 0, currency: pricePerUnit.Currency()}, nil
	}

	totalValue := q.Float64() * float64(pricePerUnit.Amount())
	roundedAmount := int64(math.RoundToEven(totalValue))

	return NewMoney(roundedAmount, pricePerUnit.Currency())
}

func (q Quantity) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Value float64 `json:"value"`
		Unit  Unit    `json:"unit"`
	}{
		Value: q.Float64(),
		Unit:  q.unit,
	})
}

func (q *Quantity) UnmarshalJSON(data []byte) error {
	dto := &struct {
		Value float64 `json:"value"`
		Unit  Unit    `json:"unit"`
	}{}

	if err := json.Unmarshal(data, dto); err != nil {
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
