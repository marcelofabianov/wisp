package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/marcelofabianov/fault"
)

var registeredTimezones = make(map[string]struct{})

type Timezone struct {
	location *time.Location
}

var ZeroTimezone = Timezone{}

func RegisterTimezones(names ...string) error {
	for _, name := range names {
		_, err := time.LoadLocation(name)
		if err != nil {
			return fault.Wrap(err, "failed to validate timezone for registration", fault.WithContext("name", name))
		}
	}

	for _, name := range names {
		registeredTimezones[name] = struct{}{}
	}

	return nil
}

func ClearRegisteredTimezones() {
	registeredTimezones = make(map[string]struct{})
}

func IsTimezoneRegistered(name string) bool {
	_, ok := registeredTimezones[name]
	return ok
}

func NewTimezone(name string) (Timezone, error) {
	if name == "" {
		return ZeroTimezone, fault.New("timezone name cannot be empty", fault.WithCode(fault.Invalid))
	}

	if !IsTimezoneRegistered(name) {
		return ZeroTimezone, fault.New(
			"timezone is not registered in the allowed list",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_name", name),
		)
	}

	loc, _ := time.LoadLocation(name)

	return Timezone{location: loc}, nil
}

func (tz Timezone) Location() *time.Location {
	return tz.location
}

func (tz Timezone) Convert(t time.Time) time.Time {
	if tz.IsZero() {
		return t
	}
	return t.In(tz.location)
}

func (tz Timezone) String() string {
	if tz.IsZero() {
		return ""
	}
	return tz.location.String()
}

func (tz Timezone) IsZero() bool {
	return tz.location == nil
}

func (tz Timezone) Equals(other Timezone) bool {
	if tz.IsZero() || other.IsZero() {
		return tz.IsZero() == other.IsZero()
	}
	return tz.location.String() == other.location.String()
}

func (tz Timezone) MarshalJSON() ([]byte, error) {
	if tz.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(tz.String())
}

func (tz *Timezone) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*tz = ZeroTimezone
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "Timezone must be a valid JSON string or null", fault.WithCode(fault.Invalid))
	}

	newTz, err := NewTimezone(s)
	if err != nil {
		return err
	}
	*tz = newTz
	return nil
}

func (tz Timezone) Value() (driver.Value, error) {
	if tz.IsZero() {
		return nil, nil
	}
	return tz.String(), nil
}

func (tz *Timezone) Scan(src interface{}) error {
	if src == nil {
		*tz = ZeroTimezone
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for Timezone", fault.WithCode(fault.Invalid))
	}

	newTz, err := NewTimezone(s)
	if err != nil {
		return err
	}
	*tz = newTz
	return nil
}
