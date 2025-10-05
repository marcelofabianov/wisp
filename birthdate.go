package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/marcelofabianov/fault"
)

// defaultLegalAge is the default age considered to be of legal age.
// It can be configured globally using SetLegalAge.
var defaultLegalAge = 18

// SetLegalAge configures the global default legal age used by the IsOfAge method.
// This allows customization based on different legal requirements (e.g., 21 in some regions).
// The age must be a positive integer.
func SetLegalAge(age int) {
	if age > 0 {
		defaultLegalAge = age
	}
}

// BirthDate represents a person's date of birth.
// It is a value object that wraps a wisp.Date and ensures the date is not in the future.
// It provides methods to calculate age and check for legal age.
//
// The zero value is ZeroBirthDate.
//
// Examples:
//   bd, err := NewBirthDate(1990, time.January, 1)
//   age := bd.Age(Today()) // Calculates age based on the current date
//   isAdult := bd.IsOfAge(Today())
type BirthDate struct {
	date Date
}

// ZeroBirthDate represents the zero value for the BirthDate type.
var ZeroBirthDate BirthDate

// NewBirthDate creates a new BirthDate from a year, month, and day.
// It returns an error if the date is invalid or in the future.
func NewBirthDate(year int, month time.Month, day int) (BirthDate, error) {
	d, err := NewDate(year, month, day)
	if err != nil {
		return ZeroBirthDate, err
	}

	if d.After(Today()) {
		return ZeroBirthDate, fault.New(
			"birth date cannot be in the future",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_date", d.String()),
		)
	}

	return BirthDate{date: d}, nil
}

// ParseBirthDate creates a new BirthDate by parsing a string in YYYY-MM-DD format.
// It returns an error if the string is not a valid date or is in the future.
func ParseBirthDate(value string) (BirthDate, error) {
	d, err := ParseDate(value)
	if err != nil {
		return ZeroBirthDate, err
	}

	return NewBirthDate(d.Year(), d.Month(), d.Day())
}

// Date returns the underlying wisp.Date value.
func (bd BirthDate) Date() Date {
	return bd.date
}

// IsZero returns true if the BirthDate is the zero value.
func (bd BirthDate) IsZero() bool {
	return bd.date.IsZero()
}

// Age calculates the person's age in years as of a given reference date (`today`).
func (bd BirthDate) Age(today Date) int {
	if bd.IsZero() {
		return 0
	}
	age := today.Year() - bd.date.Year()
	if today.Month() < bd.date.Month() || (today.Month() == bd.date.Month() && today.Day() < bd.date.Day()) {
		age--
	}
	return age
}

// IsOfAge checks if the person has reached the legal age as of a given reference date (`today`).
// The legal age is determined by the global `defaultLegalAge`, which can be set via `SetLegalAge`.
func (bd BirthDate) IsOfAge(today Date) bool {
	if bd.IsZero() {
		return false
	}
	return bd.Age(today) >= defaultLegalAge
}

// AnniversaryThisYear returns the date of the birthday anniversary for the current year of a given reference date (`today`).
func (bd BirthDate) AnniversaryThisYear(today Date) Date {
	if bd.IsZero() {
		return ZeroDate
	}

	anniversaryTime := time.Date(today.Year(), bd.date.Month(), bd.date.Day(), 0, 0, 0, 0, time.UTC)
	return Date{t: anniversaryTime}
}

// HasAnniversaryPassed checks if the birthday for the current year has already passed as of a given reference date (`today`).
func (bd BirthDate) HasAnniversaryPassed(today Date) bool {
	if bd.IsZero() {
		return false
	}
	return today.After(bd.AnniversaryThisYear(today))
}

// String returns the birth date formatted as a YYYY-MM-DD string.
func (bd BirthDate) String() string {
	return bd.date.String()
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the BirthDate as a YYYY-MM-DD string.
func (bd BirthDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(bd.date)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string in YYYY-MM-DD format into a BirthDate, with validation.
func (bd *BirthDate) UnmarshalJSON(data []byte) error {
	var d Date
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	bdObj, err := NewBirthDate(d.Year(), d.Month(), d.Day())
	if err != nil {
		return err
	}
	*bd = bdObj
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the BirthDate as a time.Time value.
func (bd BirthDate) Value() (driver.Value, error) {
	return bd.date.Value()
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a time.Time from the database and converts it into a BirthDate, with validation.
func (bd *BirthDate) Scan(src interface{}) error {
	var d Date
	if err := d.Scan(src); err != nil {
		return err
	}
	if d.IsZero() {
		*bd = ZeroBirthDate
		return nil
	}

	bdObj, err := NewBirthDate(d.Year(), d.Month(), d.Day())
	if err != nil {
		return err
	}
	*bd = bdObj
	return nil
}
