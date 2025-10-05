package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/marcelofabianov/fault"
)

// DayOfWeek is a value object representing a day of the week (Sunday, Monday, etc.).
// It is an alias for time.Weekday, providing type safety and useful methods.
// This type ensures that only valid days of the week are used.
//
// It can be parsed from a string and implements JSON and database interfaces.
//
// Examples:
//   dow, err := ParseDayOfWeek("Monday")
//   isWeekend := dow.IsWeekend()
type DayOfWeek time.Weekday

// Constants representing each day of the week.
const (
	Sunday    DayOfWeek = DayOfWeek(time.Sunday)
	Monday    DayOfWeek = DayOfWeek(time.Monday)
	Tuesday   DayOfWeek = DayOfWeek(time.Tuesday)
	Wednesday DayOfWeek = DayOfWeek(time.Wednesday)
	Thursday  DayOfWeek = DayOfWeek(time.Thursday)
	Friday    DayOfWeek = DayOfWeek(time.Friday)
	Saturday  DayOfWeek = DayOfWeek(time.Saturday)
)

// dayOfWeekMap provides a lookup from a lowercase string to a DayOfWeek constant.
var dayOfWeekMap = map[string]DayOfWeek{
	"sunday":    Sunday,
	"monday":    Monday,
	"tuesday":   Tuesday,
	"wednesday": Wednesday,
	"thursday":  Thursday,
	"friday":    Friday,
	"saturday":  Saturday,
}

// ParseDayOfWeek creates a DayOfWeek from a string (e.g., "Monday").
// The input is case-insensitive.
// Returns an error if the string is not a valid day of the week.
func ParseDayOfWeek(s string) (DayOfWeek, error) {
	d, ok := dayOfWeekMap[strings.ToLower(strings.TrimSpace(s))]
	if !ok {
		return 0, fault.New(
			"invalid day of week string",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", s),
		)
	}
	return d, nil
}

// Weekday returns the underlying time.Weekday value.
func (d DayOfWeek) Weekday() time.Weekday {
	return time.Weekday(d)
}

// IsWeekend returns true if the day is Saturday or Sunday.
func (d DayOfWeek) IsWeekend() bool {
	return d == Saturday || d == Sunday
}

// IsWeekday returns true if the day is not a weekend day (Monday to Friday).
func (d DayOfWeek) IsWeekday() bool {
	return !d.IsWeekend()
}

// String returns the English name of the day of the week (e.g., "Monday").
func (d DayOfWeek) String() string {
	return d.Weekday().String()
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the DayOfWeek as a lowercase JSON string (e.g., "monday").
func (d DayOfWeek) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToLower(d.String()))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into a DayOfWeek, with validation.
func (d *DayOfWeek) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "DayOfWeek must be a valid JSON string", fault.WithCode(fault.Invalid))
	}

	day, err := ParseDayOfWeek(s)
	if err != nil {
		return err
	}
	*d = day
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the DayOfWeek as an integer (where Sunday=0, Monday=1, etc.).
func (d DayOfWeek) Value() (driver.Value, error) {
	return int64(d), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts an integer from the database and converts it into a DayOfWeek.
func (d *DayOfWeek) Scan(src interface{}) error {
	if src == nil {
		*d = 0 // Sunday
		return nil
	}

	var i int64
	switch v := src.(type) {
	case int64:
		i = v
	default:
		return fault.New("unsupported scan type for DayOfWeek", fault.WithCode(fault.Invalid))
	}

	if i < 0 || i > 6 {
		return fault.New("value out of range for DayOfWeek", fault.WithCode(fault.Invalid), fault.WithContext("value", i))
	}

	*d = DayOfWeek(i)
	return nil
}
