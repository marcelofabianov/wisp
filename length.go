package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"

	"github.com/marcelofabianov/fault"
)

type LengthUnit string

const (
	Meter      LengthUnit = "m"
	Centimeter LengthUnit = "cm"
	Millimeter LengthUnit = "mm"
	Kilometer  LengthUnit = "km"
	Inch       LengthUnit = "in"
	Foot       LengthUnit = "ft"
)

const (
	micrometersInAMeter = 1000000.0
	metersInAKilometer  = 1000.0
	metersInAFeoot      = 0.3048
	metersInAnInch      = 0.0254
)

type Length struct {
	micrometers int64
}

var ZeroLength = Length{}

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

func (l Length) Add(other Length) Length {
	return Length{micrometers: l.micrometers + other.micrometers}
}

func (l Length) Subtract(other Length) Length {
	return Length{micrometers: l.micrometers - other.micrometers}
}

func (l Length) IsNegative() bool {
	return l.micrometers < 0
}

func (l Length) Equals(other Length) bool {
	return l.micrometers == other.micrometers
}

func (l Length) String() string {
	m, _ := l.In(Meter)
	return fmt.Sprintf("%.3f m", m)
}

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

func (l Length) Value() (driver.Value, error) {
	return l.micrometers, nil
}

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
