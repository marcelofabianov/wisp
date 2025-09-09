package atomic

import (
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

type DateRange struct {
	start Date
	end   Date
}

var ZeroDateRange DateRange

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

func (dr DateRange) Start() Date {
	return dr.start
}

func (dr DateRange) End() Date {
	return dr.end
}

func (dr DateRange) IsZero() bool {
	return dr.start.IsZero() && dr.end.IsZero()
}

func (dr DateRange) Equals(other DateRange) bool {
	return dr.start.Equals(other.start) && dr.end.Equals(other.end)
}

func (dr DateRange) Contains(d Date) bool {
	if dr.IsZero() || d.IsZero() {
		return false
	}
	return !d.Before(dr.start) && !d.After(dr.end)
}

func (dr DateRange) Overlaps(other DateRange) bool {
	if dr.IsZero() || other.IsZero() {
		return false
	}

	return !dr.start.After(other.end) && !dr.end.Before(other.start)
}

func (dr DateRange) Days() int {
	if dr.IsZero() {
		return 0
	}

	return int(dr.end.t.Sub(dr.start.t).Hours()/24) + 1
}

func (dr DateRange) String() string {
	if dr.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s to %s", dr.start.String(), dr.end.String())
}

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
