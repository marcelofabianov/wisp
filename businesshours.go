package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/marcelofabianov/fault"
)

// BusinessHours represents a weekly schedule of opening and closing times.
// It is a value object that maps each DayOfWeek to a specific TimeRange, defining when a business is open.
// This is useful for services, stores, or any entity with regular operating hours.
//
// The zero value is EmptyBusinessHours, representing a schedule where the business is always closed.
//
// Example:
//   schedule := map[wisp.DayOfWeek]wisp.TimeRange{
//       wisp.Monday: wisp.MustNewTimeRange(wisp.NewTimeOfDay(9, 0), wisp.NewTimeOfDay(17, 0)),
//       wisp.Tuesday: wisp.MustNewTimeRange(wisp.NewTimeOfDay(9, 0), wisp.NewTimeOfDay(17, 0)),
//   }
//   bh, _ := wisp.NewBusinessHours(schedule)
//   isOpen := bh.IsOpen(time.Now()) // Checks if the current time falls within business hours
type BusinessHours struct {
	schedule map[DayOfWeek]TimeRange
}

// EmptyBusinessHours represents a business that is always closed.
var EmptyBusinessHours = BusinessHours{schedule: make(map[DayOfWeek]TimeRange)}

// NewBusinessHours creates a new BusinessHours object from a schedule map.
// The schedule maps a DayOfWeek to a TimeRange.
// It returns an error if any DayOfWeek in the map is invalid.
func NewBusinessHours(schedule map[DayOfWeek]TimeRange) (BusinessHours, error) {
	if len(schedule) == 0 {
		return EmptyBusinessHours, nil
	}

	newSchedule := make(map[DayOfWeek]TimeRange, len(schedule))
	for day, timeRange := range schedule {
		if day < Sunday || day > Saturday {
			return EmptyBusinessHours, fault.New("invalid DayOfWeek key in schedule", fault.WithCode(fault.Invalid))
		}
		newSchedule[day] = timeRange
	}
	return BusinessHours{schedule: newSchedule}, nil
}

// IsOpen checks if the business is open at a specific time `t`.
// It determines the day of the week from `t` and checks if the time of day falls within the scheduled TimeRange for that day.
// It returns false if there is no schedule for that day.
func (bh BusinessHours) IsOpen(t time.Time) bool {
	day := DayOfWeek(t.Weekday())

	timeRange, ok := bh.schedule[day]
	if !ok {
		return false
	}

	timeOfDay, err := NewTimeOfDay(t.Hour(), t.Minute())
	if err != nil {
		return false
	}

	return timeRange.Contains(timeOfDay)
}

// IsZero returns true if the BusinessHours schedule is empty.
func (bh BusinessHours) IsZero() bool {
	return len(bh.schedule) == 0
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the BusinessHours schedule into a JSON object where keys are lowercase day names (e.g., "monday").
func (bh BusinessHours) MarshalJSON() ([]byte, error) {
	stringKeyMap := make(map[string]TimeRange)
	if bh.schedule != nil {
		for day, timeRange := range bh.schedule {
			stringKeyMap[strings.ToLower(day.String())] = timeRange
		}
	}
	return json.Marshal(stringKeyMap)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object into a BusinessHours schedule.
// The JSON keys are expected to be lowercase day names.
func (bh *BusinessHours) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*bh = EmptyBusinessHours
		return nil
	}

	var stringKeyMap map[string]TimeRange
	if err := json.Unmarshal(data, &stringKeyMap); err != nil {
		return fault.Wrap(err, "invalid JSON format for BusinessHours", fault.WithCode(fault.Invalid))
	}

	if len(stringKeyMap) == 0 {
		*bh = EmptyBusinessHours
		return nil
	}

	newSchedule := make(map[DayOfWeek]TimeRange)
	for dayStr, timeRange := range stringKeyMap {
		day, err := ParseDayOfWeek(dayStr)
		if err != nil {
			return fault.Wrap(err, "invalid day of week string in BusinessHours JSON", fault.WithCode(fault.Invalid))
		}
		newSchedule[day] = timeRange
	}

	*bh = BusinessHours{schedule: newSchedule}
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the BusinessHours schedule as a JSON byte array.
func (bh BusinessHours) Value() (driver.Value, error) {
	return bh.MarshalJSON()
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a JSON byte array or string from the database and converts it into a BusinessHours schedule.
func (bh *BusinessHours) Scan(src interface{}) error {
	if src == nil {
		*bh = EmptyBusinessHours
		return nil
	}

	var bytes []byte
	switch v := src.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fault.New("unsupported scan type for BusinessHours", fault.WithCode(fault.Invalid))
	}

	if len(bytes) == 0 || string(bytes) == "{}" {
		*bh = EmptyBusinessHours
		return nil
	}

	return bh.UnmarshalJSON(bytes)
}
