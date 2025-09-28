package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"

	"github.com/marcelofabianov/fault"
)

type WeightUnit string

const (
	Kilogram WeightUnit = "kg"
	Gram     WeightUnit = "g"
	Pound    WeightUnit = "lb"
	Ounce    WeightUnit = "oz"
)

const (
	gramsInAKilogram = 1000
	gramsInAPound    = 453.59237
	gramsInAnOunce   = 28.34952
	mgInAGram        = 1000
)

type Weight struct {
	milligrams int64
}

var ZeroWeight = Weight{}

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

func (w Weight) Add(other Weight) Weight {
	return Weight{milligrams: w.milligrams + other.milligrams}
}

func (w Weight) Subtract(other Weight) Weight {
	return Weight{milligrams: w.milligrams - other.milligrams}
}

func (w Weight) IsNegative() bool {
	return w.milligrams < 0
}

func (w Weight) Equals(other Weight) bool {
	return w.milligrams == other.milligrams
}

func (w Weight) String() string {
	kg, _ := w.In(Kilogram)
	return fmt.Sprintf("%.3f kg", kg)
}

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

func (w Weight) Value() (driver.Value, error) {
	return w.milligrams, nil
}

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
