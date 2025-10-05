package wisp

import (
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

// DateRange represents a period between a start and an end date, inclusive.
// It is a value object that ensures the start date is not after the end date.
// Both the start and end dates are of the wisp.Date type.
//
// The zero value for DateRange is ZeroDateRange, where both start and end dates are zero.
//
// Examples:
//   start, _ := wisp.NewDate(2025, 1, 1)
//   end, _ := wisp.NewDate(2025, 1, 31)
//   dr, err := wisp.NewDateRange(start, end)
type DateRange struct {
	start Date
	end   Date
}

// ZeroDateRange represents the zero value for the DateRange type.
var ZeroDateRange DateRange

// NewDateRange creates a new DateRange from a start and end date.
// It returns an error if the start date is after the end date.
func NewDateRange(start, end Date) (DateRange, error) {
	if start.After(end) {
		return ZeroDateRange, fault.New(
			"start date cannot be after end date",
			fault.WithCode(fault.Invalid),
			fault.WithContext("start_date", start.String()),
			fault.WithContext("end_date", end.String()),
		)
	}
	return DateRange{start: start, end: end}, nil
}

// Start returns the start date of the range.
func (dr DateRange) Start() Date {
	return dr.start
}

// End returns the end date of the range.
func (dr DateRange) End() Date {
	return dr.end
}

// IsZero returns true if the DateRange is the zero value.
func (dr DateRange) IsZero() bool {
	return dr.start.IsZero() && dr.end.IsZero()
}

// Equals checks if two DateRange instances are equal by comparing their start and end dates.
func (dr DateRange) Equals(other DateRange) bool {
	return dr.start.Equals(other.start) && dr.end.Equals(other.end)
}

// Contains checks if a given date is within the date range (inclusive).
func (dr DateRange) Contains(d Date) bool {
	if dr.IsZero() || d.IsZero() {
		return false
	}
	return !d.Before(dr.start) && !d.After(dr.end)
}

// Overlaps checks if two date ranges have at least one day in common.
func (dr DateRange) Overlaps(other DateRange) bool {
	if dr.IsZero() || other.IsZero() {
		return false
	}

	return !dr.start.After(other.end) && !dr.end.Before(other.start)
}

// Days returns the total number of days in the range, inclusive.
// For example, a range from 2025-01-01 to 2025-01-03 has 3 days.
func (dr DateRange) Days() int {
	if dr.IsZero() {
		return 0
	}

	return int(dr.end.t.Sub(dr.start.t).Hours()/24) + 1
}

// String returns a formatted string representation of the date range, like "YYYY-MM-DD to YYYY-MM-DD".
func (dr DateRange) String() string {
	if dr.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s to %s", dr.start.String(), dr.end.String())
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the DateRange into a JSON object with "start" and "end" fields.
func (dr DateRange) MarshalJSON() ([]byte, error) {
	if dr.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(&struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}{
		Start: dr.start.String(),
		End:   dr.end.String(),
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object with "start" and "end" fields into a DateRange.
func (dr *DateRange) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*dr = ZeroDateRange
		return nil
	}

	dto := &struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}{}

	if err := json.Unmarshal(data, dto); err != nil {
		return fault.Wrap(err, "invalid JSON format for DateRange", fault.WithCode(fault.Invalid))
	}

	start, err := ParseDate(dto.Start)
	if err != nil {
		return fault.Wrap(err, "invalid start date for DateRange", fault.WithCode(fault.Invalid))
	}

	end, err := ParseDate(dto.End)
	if err != nil {
		return fault.Wrap(err, "invalid end date for DateRange", fault.WithCode(fault.Invalid))
	}

	dateRange, err := NewDateRange(start, end)
	if err != nil {
		return err
	}

	*dr = dateRange
	return nil
}
