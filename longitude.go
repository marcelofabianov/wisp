package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

type Longitude float64

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

func (l Longitude) Float64() float64 {
	return float64(l)
}

func (l Longitude) String() string {
	return fmt.Sprintf("%f", l)
}

func (l Longitude) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Float64())
}

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

func (l Longitude) Value() (driver.Value, error) {
	return l.Float64(), nil
}

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
