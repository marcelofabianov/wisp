package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marcelofabianov/fault"
)

// Day represents a day of the month, as an integer from 1 to 31.
// It is a value object that ensures the day is within a valid range.
// This type is useful for representing recurring monthly dates, such as a billing day.
//
// The zero value is ZeroDay.
//
// Examples:
//   d, err := NewDay(15) // Represents the 15th day of the month
type Day int

// ZeroDay represents the zero value for the Day type.
var ZeroDay Day

// validateDay checks if the integer value is a valid day of the month (1-31).
func validateDay(value int) error {
	if value < 1 || value > 31 {
		return fault.New(
			"day must be between 1 and 31",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", value),
		)
	}
	return nil
}

// NewDay creates a new Day.
// It returns an error if the value is not between 1 and 31.
func NewDay(value int) (Day, error) {
	if err := validateDay(value); err != nil {
		return ZeroDay, err
	}
	return Day(value), nil
}

// Int returns the integer representation of the Day.
func (d Day) Int() int {
	return int(d)
}

// IsZero returns true if the Day is the zero value.
func (d Day) IsZero() bool {
	return d == ZeroDay
}

// HasPassed checks if this day of the month has already passed in the context of a given reference time (`today`).
func (d Day) HasPassed(today time.Time) bool {
	if d.IsZero() {
		return false
	}
	return d.Int() < today.Day()
}

// DaysUntil calculates the number of days from a reference date (`today`) until the next occurrence of this day.
// It accounts for month boundaries.
func (d Day) DaysUntil(today time.Time) int {
	if d.IsZero() {
		return 0
	}

	day := d.Int()
	todayDay := today.Day()

	if day >= todayDay {
		return day - todayDay
	}

	daysInMonth := time.Date(today.Year(), today.Month()+1, 0, 0, 0, 0, 0, today.Location()).Day()
	return (daysInMonth - todayDay) + day
}

// DaysOverdue calculates the number of days that have passed since the last occurrence of this day, relative to a reference date (`today`).
func (d Day) DaysOverdue(today time.Time) int {
	if d.IsZero() {
		return 0
	}

	day := d.Int()
	todayDay := today.Day()

	if day <= todayDay {
		return todayDay - day
	}

	prevMonth := today.AddDate(0, -1, 0)
	daysInPrevMonth := time.Date(prevMonth.Year(), prevMonth.Month()+1, 0, 0, 0, 0, 0, today.Location()).Day()
	return (daysInPrevMonth - day) + todayDay
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Day as a JSON number.
func (d Day) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Int())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON number into a Day, with validation.
func (d *Day) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*d = ZeroDay
		return nil
	}

	var day int
	if err := json.Unmarshal(data, &day); err != nil {
		return fault.Wrap(err,
			"day must be a valid JSON number",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_json", string(data)),
		)
	}

	if err := validateDay(day); err != nil {
		return err
	}

	*d = Day(day)
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the Day as an int64.
func (d Day) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return int64(d.Int()), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts an int64 from the database and converts it into a Day, with validation.
func (d *Day) Scan(src interface{}) error {
	if src == nil {
		*d = ZeroDay
		return nil
	}

	var day int64
	switch v := src.(type) {
	case int64:
		day = v
	default:
		return fault.New(
			"unsupported scan type for Day",
			fault.WithCode(fault.Invalid),
			fault.WithContext("received_type", fmt.Sprintf("%T", src)),
		)
	}

	dayAsInt := int(day)
	if err := validateDay(dayAsInt); err != nil {
		return err
	}

	*d = Day(dayAsInt)
	return nil
}
