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

// TimeOfDay represents a specific time of day (hour and minute), independent of any date or timezone.
// It is stored as the number of minutes from midnight, which simplifies comparisons and calculations.
// This value object is useful for representing schedules, business hours, or any time-based logic.
//
// The zero value is ZeroTimeOfDay, representing 00:00.
//
// Examples:
//   t, err := NewTimeOfDay(9, 30) // 09:30
//   t, err := ParseTimeOfDay("17:00") // 17:00
type TimeOfDay struct {
	minutesFromMidnight int
}

// ZeroTimeOfDay represents the zero value for TimeOfDay (00:00).
var ZeroTimeOfDay = TimeOfDay{}

// NewTimeOfDay creates a new TimeOfDay from an hour and minute.
// It returns an error if the hour is not between 0-23 or the minute is not between 0-59.
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

// ParseTimeOfDay creates a new TimeOfDay by parsing a string in HH:MM format.
// It returns an error if the string is not in the correct format.
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

// Hour returns the hour component of the time (0-23).
func (t TimeOfDay) Hour() int {
	return t.minutesFromMidnight / minutesInHour
}

// Minute returns the minute component of the time (0-59).
func (t TimeOfDay) Minute() int {
	return t.minutesFromMidnight % minutesInHour
}

// IsZero returns true if the TimeOfDay is the zero value (00:00).
func (t TimeOfDay) IsZero() bool {
	return t.minutesFromMidnight == 0
}

// Before checks if this TimeOfDay is before another.
func (t TimeOfDay) Before(other TimeOfDay) bool {
	return t.minutesFromMidnight < other.minutesFromMidnight
}

// After checks if this TimeOfDay is after another.
func (t TimeOfDay) After(other TimeOfDay) bool {
	return t.minutesFromMidnight > other.minutesFromMidnight
}

// MustNewTimeOfDay is like NewTimeOfDay but panics if the time is invalid.
// This is useful for initializing constant-like time values.
func MustNewTimeOfDay(hour, minute int) TimeOfDay {
	tod, err := NewTimeOfDay(hour, minute)
	if err != nil {
		panic(err)
	}
	return tod
}

// String returns the time formatted as an HH:MM string.
func (t TimeOfDay) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the TimeOfDay as an HH:MM formatted JSON string.
func (t TimeOfDay) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string in HH:MM format into a TimeOfDay.
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

// Value implements the driver.Valuer interface for database storage.
// It returns the time as the total number of minutes from midnight.
func (t TimeOfDay) Value() (driver.Value, error) {
	return int64(t.minutesFromMidnight), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts an integer (minutes from midnight) from the database and converts it into a TimeOfDay.
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
