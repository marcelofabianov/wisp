package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

type Latitude float64

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

func (l Latitude) Float64() float64 {
	return float64(l)
}

func (l Latitude) String() string {
	return fmt.Sprintf("%f", l)
}

func (l Latitude) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Float64())
}

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

func (l Latitude) Value() (driver.Value, error) {
	return l.Float64(), nil
}

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
