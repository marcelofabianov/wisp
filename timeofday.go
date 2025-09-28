package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/marcelofabianov/fault"
)

const (
	minutesInHour = 60
	minutesInDay  = 24 * minutesInHour
)

type TimeOfDay struct {
	minutesFromMidnight int
}

var ZeroTimeOfDay = TimeOfDay{}

func NewTimeOfDay(hour, minute int) (TimeOfDay, error) {
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return ZeroTimeOfDay, fault.New(
			"invalid time components provided",
			fault.WithCode(fault.Invalid),
			fault.WithContext("hour", hour),
			fault.WithContext("minute", minute),
		)
	}
	totalMinutes := hour*minutesInHour + minute
	return TimeOfDay{minutesFromMidnight: totalMinutes}, nil
}

func ParseTimeOfDay(s string) (TimeOfDay, error) {
	trimmed := strings.TrimSpace(s)

	parts := strings.Split(trimmed, ":")
	if len(parts) != 2 {
		return ZeroTimeOfDay, fault.New("time must be in HH:MM format", fault.WithCode(fault.Invalid), fault.WithContext("input", s))
	}

	hourStr := parts[0]
	minuteStr := parts[1]

	if len(hourStr) != 2 || len(minuteStr) != 2 {
		return ZeroTimeOfDay, fault.New("time must use two digits for hour and minute (HH:MM)", fault.WithCode(fault.Invalid), fault.WithContext("input", s))
	}

	h, err := strconv.Atoi(hourStr)
	if err != nil {
		return ZeroTimeOfDay, fault.Wrap(err, "hour part is not a valid number", fault.WithCode(fault.Invalid), fault.WithContext("hour_part", hourStr))
	}

	m, err := strconv.Atoi(minuteStr)
	if err != nil {
		return ZeroTimeOfDay, fault.Wrap(err, "minute part is not a valid number", fault.WithCode(fault.Invalid), fault.WithContext("minute_part", minuteStr))
	}

	return NewTimeOfDay(h, m)
}

func (t TimeOfDay) Hour() int {
	return t.minutesFromMidnight / minutesInHour
}

func (t TimeOfDay) Minute() int {
	return t.minutesFromMidnight % minutesInHour
}

func (t TimeOfDay) IsZero() bool {
	return t.minutesFromMidnight == 0
}

func (t TimeOfDay) Before(other TimeOfDay) bool {
	return t.minutesFromMidnight < other.minutesFromMidnight
}

func (t TimeOfDay) After(other TimeOfDay) bool {
	return t.minutesFromMidnight > other.minutesFromMidnight
}

func MustNewTimeOfDay(hour, minute int) TimeOfDay {
	tod, err := NewTimeOfDay(hour, minute)
	if err != nil {
		panic(err)
	}
	return tod
}

func (t TimeOfDay) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

func (t TimeOfDay) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TimeOfDay) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "TimeOfDay must be a valid JSON string", fault.WithCode(fault.Invalid))
	}
	tod, err := ParseTimeOfDay(s)
	if err != nil {
		return err
	}
	*t = tod
	return nil
}

func (t TimeOfDay) Value() (driver.Value, error) {
	return int64(t.minutesFromMidnight), nil
}

func (t *TimeOfDay) Scan(src interface{}) error {
	if src == nil {
		*t = ZeroTimeOfDay
		return nil
	}
	var min int64
	switch v := src.(type) {
	case int64:
		min = v
	default:
		return fault.New("unsupported scan type for TimeOfDay", fault.WithCode(fault.Invalid))
	}
	if min < 0 || min >= minutesInDay {
		return fault.New("value out of range for TimeOfDay", fault.WithCode(fault.Invalid), fault.WithContext("value", min))
	}
	*t = TimeOfDay{minutesFromMidnight: int(min)}
	return nil
}
