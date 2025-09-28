package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/marcelofabianov/fault"
)

type BusinessHours struct {
	schedule map[DayOfWeek]TimeRange
}

var EmptyBusinessHours = BusinessHours{schedule: make(map[DayOfWeek]TimeRange)}

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

func (bh BusinessHours) IsZero() bool {
	return len(bh.schedule) == 0
}

func (bh BusinessHours) MarshalJSON() ([]byte, error) {
	stringKeyMap := make(map[string]TimeRange)
	if bh.schedule != nil {
		for day, timeRange := range bh.schedule {
			stringKeyMap[strings.ToLower(day.String())] = timeRange
		}
	}
	return json.Marshal(stringKeyMap)
}

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

func (bh BusinessHours) Value() (driver.Value, error) {
	return bh.MarshalJSON()
}

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
