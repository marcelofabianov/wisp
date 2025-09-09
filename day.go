package atomic

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marcelofabianov/fault"
)

type Day int

var ZeroDay Day

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

func NewDay(value int) (Day, error) {
	if err := validateDay(value); err != nil {
		return ZeroDay, err
	}
	return Day(value), nil
}

func (d Day) Int() int {
	return int(d)
}

func (d Day) IsZero() bool {
	return d == ZeroDay
}

func (d Day) HasPassed(today time.Time) bool {
	if d.IsZero() {
		return false
	}
	return d.Int() < today.Day()
}

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

func (d Day) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Int())
}

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

func (d Day) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return int64(d.Int()), nil
}

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
