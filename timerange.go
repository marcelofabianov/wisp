package wisp

import (
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

type TimeRange struct {
	start TimeOfDay
	end   TimeOfDay
}

var ZeroTimeRange = TimeRange{}

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

func (tr TimeRange) Start() TimeOfDay {
	return tr.start
}

func (tr TimeRange) End() TimeOfDay {
	return tr.end
}

func (tr TimeRange) IsZero() bool {
	return tr.start.IsZero() && tr.end.IsZero()
}

func (tr TimeRange) Contains(t TimeOfDay) bool {
	return !t.Before(tr.start) && t.Before(tr.end)
}

func (tr TimeRange) String() string {
	return fmt.Sprintf("%s-%s", tr.start.String(), tr.end.String())
}

func (tr TimeRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}{
		Start: tr.start.String(),
		End:   tr.end.String(),
	})
}

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
