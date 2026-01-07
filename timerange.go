package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

// TimeRange represents a period between a start and an end time of day.
// It is a value object that ensures the start time is strictly before the end time.
// This is useful for defining periods like business hours or appointment slots.
//
// The range is exclusive of the end time, i.e., [start, end).
//
// Examples:
//   start, _ := NewTimeOfDay(9, 0)
//   end, _ := NewTimeOfDay(17, 0)
//   tr, err := NewTimeRange(start, end) // Represents 09:00 to 17:00
type TimeRange struct {
	start TimeOfDay
	end   TimeOfDay
}

// ZeroTimeRange represents the zero value for the TimeRange type.
var ZeroTimeRange = TimeRange{}

// NewTimeRange creates a new TimeRange from a start and end TimeOfDay.
// It returns an error if the start time is not before the end time.
func NewTimeRange(start, end TimeOfDay) (TimeRange, error) {
	if !start.Before(end) {
		return ZeroTimeRange, fault.New(
			"start time must be before end time in a time range",
			fault.WithCode(fault.Invalid),
			fault.WithContext("start", start.String()),
			fault.WithContext("end", end.String()),
		)
	}
	return TimeRange{start: start, end: end}, nil
}

// Start returns the start time of the range.
func (tr TimeRange) Start() TimeOfDay {
	return tr.start
}

// End returns the end time of the range.
func (tr TimeRange) End() TimeOfDay {
	return tr.end
}

// IsZero returns true if the TimeRange is the zero value.
func (tr TimeRange) IsZero() bool {
	return tr.start.IsZero() && tr.end.IsZero()
}

// Contains checks if a given TimeOfDay is within the time range.
// The check is inclusive of the start time and exclusive of the end time: [start, end).
func (tr TimeRange) Contains(t TimeOfDay) bool {
	return !t.Before(tr.start) && t.Before(tr.end)
}

// String returns a formatted string representation of the time range, like "HH:MM-HH:MM".
func (tr TimeRange) String() string {
	return fmt.Sprintf("%s-%s", tr.start.String(), tr.end.String())
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the TimeRange into a JSON object with "start" and "end" fields.
func (tr TimeRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}{
		Start: tr.start.String(),
		End:   tr.end.String(),
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object with "start" and "end" fields into a TimeRange.
func (tr *TimeRange) UnmarshalJSON(data []byte) error {
	dto := &struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}{}
	if err := json.Unmarshal(data, &dto); err != nil {
		return fault.Wrap(err, "invalid JSON format for TimeRange", fault.WithCode(fault.Invalid))
	}

	start, err := ParseTimeOfDay(dto.Start)
	if err != nil {
		return fault.Wrap(err, "invalid start time for TimeRange")
	}
	end, err := ParseTimeOfDay(dto.End)
	if err != nil {
		return fault.Wrap(err, "invalid end time for TimeRange")
	}

	timeRange, err := NewTimeRange(start, end)
	if err != nil {
		return err
	}
	*tr = timeRange
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the TimeRange as a JSON string or nil if it's the zero value.
func (tr TimeRange) Value() (driver.Value, error) {
	if tr.IsZero() {
		return nil, nil
	}

	data, err := tr.MarshalJSON()
	if err != nil {
		return nil, fault.Wrap(err,
			"failed to marshal time range for database storage",
			fault.WithCode(fault.Internal),
		)
	}

	return string(data), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts string or []byte values containing JSON and validates them as TimeRange.
func (tr *TimeRange) Scan(src interface{}) error {
	if src == nil {
		*tr = ZeroTimeRange
		return nil
	}

	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fault.New(
			"unsupported scan type for TimeRange",
			fault.WithCode(fault.Invalid),
			fault.WithContext("received_type", fmt.Sprintf("%T", src)),
		)
	}

	if err := tr.UnmarshalJSON(data); err != nil {
		return err
	}

	return nil
}
