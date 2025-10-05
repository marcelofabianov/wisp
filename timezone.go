package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/marcelofabianov/fault"
)

// registeredTimezones holds the set of IANA timezone names that are allowed in the application.
var registeredTimezones = make(map[string]struct{})

// Timezone is a value object representing an IANA timezone (e.g., "America/Sao_Paulo", "UTC").
// It ensures that only valid and explicitly registered timezones are used throughout the application.
// Before a timezone can be used, it must be added to a global registry via `RegisterTimezones`.
// This prevents the use of arbitrary or incorrect timezone strings.
//
// The zero value is ZeroTimezone.
//
// Example:
//   wisp.RegisterTimezones("America/Sao_Paulo", "UTC")
//   tz, err := wisp.NewTimezone("America/Sao_Paulo")
//   nowInSP := tz.Convert(time.Now())
type Timezone struct {
	location *time.Location
}

// ZeroTimezone represents the zero value for the Timezone type.
var ZeroTimezone = Timezone{}

// RegisterTimezones adds one or more IANA timezone names to the global registry of allowed timezones.
// It validates each name by attempting to load it. If any name is invalid, it returns an error and no timezones are registered.
// This function should be called during application startup to define the set of supported timezones.
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

// ClearRegisteredTimezones removes all previously registered timezones from the global registry.
// This is primarily useful for testing purposes to ensure a clean state between tests.
func ClearRegisteredTimezones() {
	registeredTimezones = make(map[string]struct{})
}

// IsTimezoneRegistered checks if a given timezone name is in the global registry.
func IsTimezoneRegistered(name string) bool {
	_, ok := registeredTimezones[name]
	return ok
}

// NewTimezone creates a new Timezone from an IANA timezone name.
// It returns an error if the name is empty or has not been registered via `RegisterTimezones`.
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

// Location returns the underlying *time.Location value.
func (tz Timezone) Location() *time.Location {
	return tz.location
}

// Convert returns a new time.Time instance converted to this timezone.
// If the timezone is zero, it returns the original time unchanged.
func (tz Timezone) Convert(t time.Time) time.Time {
	if tz.IsZero() {
		return t
	}
	return t.In(tz.location)
}

// String returns the IANA name of the timezone.
func (tz Timezone) String() string {
	if tz.IsZero() {
		return ""
	}
	return tz.location.String()
}

// IsZero returns true if the Timezone is the zero value.
func (tz Timezone) IsZero() bool {
	return tz.location == nil
}

// Equals checks if two Timezone instances are the same.
func (tz Timezone) Equals(other Timezone) bool {
	if tz.IsZero() || other.IsZero() {
		return tz.IsZero() == other.IsZero()
	}
	return tz.location.String() == other.location.String()
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Timezone as its IANA name string.
func (tz Timezone) MarshalJSON() ([]byte, error) {
	if tz.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(tz.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into a Timezone, validating against the registry.
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

// Value implements the driver.Valuer interface for database storage.
// It returns the Timezone as its IANA name string.
func (tz Timezone) Value() (driver.Value, error) {
	if tz.IsZero() {
		return nil, nil
	}
	return tz.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a string from the database and converts it into a Timezone, with validation.
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
