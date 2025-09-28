package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/marcelofabianov/fault"
)

type DayOfWeek time.Weekday

const (
	Sunday    DayOfWeek = DayOfWeek(time.Sunday)
	Monday    DayOfWeek = DayOfWeek(time.Monday)
	Tuesday   DayOfWeek = DayOfWeek(time.Tuesday)
	Wednesday DayOfWeek = DayOfWeek(time.Wednesday)
	Thursday  DayOfWeek = DayOfWeek(time.Thursday)
	Friday    DayOfWeek = DayOfWeek(time.Friday)
	Saturday  DayOfWeek = DayOfWeek(time.Saturday)
)

var dayOfWeekMap = map[string]DayOfWeek{
	"sunday":    Sunday,
	"monday":    Monday,
	"tuesday":   Tuesday,
	"wednesday": Wednesday,
	"thursday":  Thursday,
	"friday":    Friday,
	"saturday":  Saturday,
}

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

func (d DayOfWeek) Weekday() time.Weekday {
	return time.Weekday(d)
}

func (d DayOfWeek) IsWeekend() bool {
	return d == Saturday || d == Sunday
}

func (d DayOfWeek) IsWeekday() bool {
	return !d.IsWeekend()
}

func (d DayOfWeek) String() string {
	return d.Weekday().String()
}

func (d DayOfWeek) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToLower(d.String()))
}

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

func (d DayOfWeek) Value() (driver.Value, error) {
	return int64(d), nil
}

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
