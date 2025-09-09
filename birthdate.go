package atomic

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/marcelofabianov/fault"
)

var defaultLegalAge = 18

func SetLegalAge(age int) {
	if age > 0 {
		defaultLegalAge = age
	}
}

type BirthDate struct {
	date Date
}

var ZeroBirthDate BirthDate

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

func ParseBirthDate(value string) (BirthDate, error) {
	d, err := ParseDate(value)
	if err != nil {
		return ZeroBirthDate, err
	}

	return NewBirthDate(d.Year(), d.Month(), d.Day())
}

func (bd BirthDate) Date() Date {
	return bd.date
}

func (bd BirthDate) IsZero() bool {
	return bd.date.IsZero()
}

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

func (bd BirthDate) IsOfAge(today Date) bool {
	if bd.IsZero() {
		return false
	}
	return bd.Age(today) >= defaultLegalAge
}

func (bd BirthDate) AnniversaryThisYear(today Date) Date {
	if bd.IsZero() {
		return ZeroDate
	}

	anniversaryTime := time.Date(today.Year(), bd.date.Month(), bd.date.Day(), 0, 0, 0, 0, time.UTC)
	return Date{t: anniversaryTime}
}

func (bd BirthDate) HasAnniversaryPassed(today Date) bool {
	if bd.IsZero() {
		return false
	}
	return today.After(bd.AnniversaryThisYear(today))
}

func (bd BirthDate) String() string {
	return bd.date.String()
}

func (bd BirthDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(bd.date)
}

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

func (bd BirthDate) Value() (driver.Value, error) {
	return bd.date.Value()
}

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
