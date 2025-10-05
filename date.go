package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marcelofabianov/fault"
)

const (
	// iso8601DateFormat defines the standard string format for dates (YYYY-MM-DD).
	iso8601DateFormat = "2006-01-02"
)

// Date represents a calendar date (year, month, day) without time or timezone information.
// It is a value object that ensures date integrity and provides useful methods for comparison and manipulation.
// Internally, it is stored as a time.Time at midnight UTC.
//
// The zero value for Date is ZeroDate, which corresponds to a zero time.Time.
//
// Examples:
//   d, err := NewDate(2025, time.October, 5)
//   today := Today()
//   parsed, err := ParseDate("2025-10-05")
type Date struct {
	t time.Time
}

// ZeroDate represents the zero value for the Date type.
var ZeroDate Date

// NewDate creates a new Date from a year, month, and day.
// It validates that the provided components form a valid calendar date.
// For example, it will reject day 32 of a month.
func NewDate(year int, month time.Month, day int) (Date, error) {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	if t.Year() != year || t.Month() != month || t.Day() != day {
		return ZeroDate, fault.New(
			"invalid date provided",
			fault.WithCode(fault.Invalid),
			fault.WithContext("year", year),
			fault.WithContext("month", int(month)),
			fault.WithContext("day", day),
		)
	}

	return Date{t: t}, nil
}

// Today returns a new Date representing the current day in UTC.
func Today() Date {
	now := time.Now().UTC()
	return Date{t: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)}
}

// ParseDate creates a new Date by parsing a string in YYYY-MM-DD format.
// Returns an error if the string is not in the required format.
func ParseDate(value string) (Date, error) {
	t, err := time.Parse(iso8601DateFormat, value)
	if err != nil {
		return ZeroDate, fault.Wrap(err,
			"date must be in YYYY-MM-DD format",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input", value),
		)
	}
	return Date{t: t}, nil
}

// Year returns the year component of the date.
func (d Date) Year() int {
	return d.t.Year()
}

// Month returns the month component of the date.
func (d Date) Month() time.Month {
	return d.t.Month()
}

// Day returns the day component of the date.
func (d Date) Day() int {
	return d.t.Day()
}

// IsZero returns true if the Date is the zero value.
func (d Date) IsZero() bool {
	return d.t.IsZero()
}

// Equals checks if two Date instances represent the same calendar date.
func (d Date) Equals(other Date) bool {
	return d.t.Equal(other.t)
}

// Before checks if the Date is before another Date.
func (d Date) Before(other Date) bool {
	return d.t.Before(other.t)
}

// After checks if the Date is after another Date.
func (d Date) After(other Date) bool {
	return d.t.After(other.t)
}

// AddDays returns a new Date with the specified number of days added.
func (d Date) AddDays(days int) Date {
	return Date{t: d.t.AddDate(0, 0, days)}
}

// AddMonths returns a new Date with the specified number of months added.
func (d Date) AddMonths(months int) Date {
	return Date{t: d.t.AddDate(0, months, 0)}
}

// AddYears returns a new Date with the specified number of years added.
func (d Date) AddYears(years int) Date {
	return Date{t: d.t.AddDate(years, 0, 0)}
}

// String returns the date formatted as a YYYY-MM-DD string.
func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return d.t.Format(iso8601DateFormat)
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Date as a YYYY-MM-DD string or null if it's a zero value.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(d.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string in YYYY-MM-DD format into a Date.
func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*d = ZeroDate
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "Date must be a valid JSON string or null", fault.WithCode(fault.Invalid))
	}

	date, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = date
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the Date as a time.Time value or nil if it's a zero value.
func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.t, nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a time.Time from the database and converts it into a Date, ignoring the time part.
func (d *Date) Scan(src interface{}) error {
	if src == nil {
		*d = ZeroDate
		return nil
	}

	switch v := src.(type) {
	case time.Time:
		*d = Date{t: time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.UTC)}
		return nil
	default:
		return fault.New("unsupported scan type for Date", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}
}
