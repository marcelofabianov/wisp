package atomic

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marcelofabianov/fault"
)

type NullableTime struct {
	Time  time.Time
	Valid bool
}

var EmptyNullableTime = NullableTime{}

func NewNullableTime(t time.Time) NullableTime {
	return NullableTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

func (nt NullableTime) IsZero() bool {
	return !nt.Valid
}

func (nt NullableTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}

	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return fault.Wrap(err, "NullableTime must be a valid JSON time string or null", fault.WithCode(fault.Invalid))
	}

	nt.Time = t
	nt.Valid = true
	return nil
}

func (nt NullableTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt *NullableTime) Scan(src interface{}) error {
	if src == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}

	switch v := src.(type) {
	case time.Time:
		nt.Time, nt.Valid = v, true
		return nil
	default:
		return fault.New("unsupported scan type for NullableTime", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}
}
