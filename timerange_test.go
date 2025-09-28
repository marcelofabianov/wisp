package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type TimeRangeSuite struct {
	suite.Suite
}

func TestTimeRangeSuite(t *testing.T) {
	suite.Run(t, new(TimeRangeSuite))
}

func (s *TimeRangeSuite) TestNewTimeRange() {
	start := wisp.MustNewTimeOfDay(9, 0)
	end := wisp.MustNewTimeOfDay(17, 30)

	s.Run("should create a valid time range", func() {
		tr, err := wisp.NewTimeRange(start, end)
		s.Require().NoError(err)
		s.False(tr.IsZero())
		s.Equal(start, tr.Start())
		s.Equal(end, tr.End())
	})

	s.Run("should fail if start is not before end", func() {
		_, err := wisp.NewTimeRange(end, start)
		s.Require().Error(err)

		_, err = wisp.NewTimeRange(start, start)
		s.Require().Error(err)
	})
}

func (s *TimeRangeSuite) TestTimeRange_Contains() {
	start := wisp.MustNewTimeOfDay(9, 0)
	end := wisp.MustNewTimeOfDay(17, 0)
	tr, _ := wisp.NewTimeRange(start, end)

	testCases := []struct {
		name     string
		time     wisp.TimeOfDay
		expected bool
	}{
		{"time before range", wisp.MustNewTimeOfDay(8, 59), false},
		{"exact start time of range", wisp.MustNewTimeOfDay(9, 0), true},
		{"time within range", wisp.MustNewTimeOfDay(12, 30), true},
		{"exact end time of range", wisp.MustNewTimeOfDay(17, 0), false},
		{"time after range", wisp.MustNewTimeOfDay(17, 1), false},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Equal(tc.expected, tr.Contains(tc.time))
		})
	}
}

func (s *TimeRangeSuite) TestTimeRange_JSON() {
	start := wisp.MustNewTimeOfDay(9, 30)
	end := wisp.MustNewTimeOfDay(18, 0)
	tr, _ := wisp.NewTimeRange(start, end)

	s.Run("should marshal and unmarshal correctly", func() {
		data, err := json.Marshal(tr)
		s.Require().NoError(err)
		s.JSONEq(`{"start": "09:30", "end": "18:00"}`, string(data))

		var unmarshaledTR wisp.TimeRange
		err = json.Unmarshal(data, &unmarshaledTR)
		s.Require().NoError(err)
		s.Equal(tr, unmarshaledTR)
	})

	s.Run("should fail for invalid JSON data", func() {
		var unmarshaledTR wisp.TimeRange
		invalidJSON := `{"start": "18:00", "end": "09:00"}`
		err := json.Unmarshal([]byte(invalidJSON), &unmarshaledTR)
		s.Require().Error(err)
	})
}
