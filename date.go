package atomic

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marcelofabianov/fault"
)

const (
	iso8601DateFormat = "2006-01-02"
)

type Date struct {
	t time.Time
}

var ZeroDate Date

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

func Today() Date {
	now := time.Now().UTC()
	return Date{t: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)}
}

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

func (d Date) Year() int {
	return d.t.Year()
}

func (d Date) Month() time.Month {
	return d.t.Month()
}

func (d Date) Day() int {
	return d.t.Day()
}

func (d Date) IsZero() bool {
	return d.t.IsZero()
}

func (d Date) Equals(other Date) bool {
	return d.t.Equal(other.t)
}

func (d Date) Before(other Date) bool {
	return d.t.Before(other.t)
}

func (d Date) After(other Date) bool {
	return d.t.After(other.t)
}

func (d Date) AddDays(days int) Date {
	return Date{t: d.t.AddDate(0, 0, days)}
}

func (d Date) AddMonths(months int) Date {
	return Date{t: d.t.AddDate(0, months, 0)}
}

func (d Date) AddYears(years int) Date {
	return Date{t: d.t.AddDate(years, 0, 0)}
}

func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return d.t.Format(iso8601DateFormat)
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(d.String())
}

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

func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.t, nil
}

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
