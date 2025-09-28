package wisp_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type BusinessHoursSuite struct {
	suite.Suite
	schedule map[wisp.DayOfWeek]wisp.TimeRange
	bh       wisp.BusinessHours
}

func (s *BusinessHoursSuite) SetupSuite() {
	start, _ := wisp.ParseTimeOfDay("09:00")
	end, _ := wisp.ParseTimeOfDay("18:00")
	weekendEnd, _ := wisp.ParseTimeOfDay("12:00")
	weekdayRange, _ := wisp.NewTimeRange(start, end)
	weekendRange, _ := wisp.NewTimeRange(start, weekendEnd)

	s.schedule = map[wisp.DayOfWeek]wisp.TimeRange{
		wisp.Monday:    weekdayRange,
		wisp.Tuesday:   weekdayRange,
		wisp.Wednesday: weekdayRange,
		wisp.Thursday:  weekdayRange,
		wisp.Friday:    weekdayRange,
		wisp.Saturday:  weekendRange,
	}

	s.bh, _ = wisp.NewBusinessHours(s.schedule)
}

func TestBusinessHoursSuite(t *testing.T) {
	suite.Run(t, new(BusinessHoursSuite))
}

func (s *BusinessHoursSuite) TestIsOpen() {
	testCases := []struct {
		name     string
		time     time.Time
		expected bool
	}{
		{"during a weekday, within hours", time.Date(2025, 9, 29, 10, 30, 0, 0, time.UTC), true}, // Monday 10:30
		{"during a weekday, before hours", time.Date(2025, 9, 29, 8, 59, 0, 0, time.UTC), false}, // Monday 08:59
		{"during a weekday, after hours", time.Date(2025, 9, 29, 18, 0, 0, 0, time.UTC), false},  // Monday 18:00 (exclusive)
		{"on a weekend, within hours", time.Date(2025, 10, 4, 11, 0, 0, 0, time.UTC), true},      // Saturday 11:00
		{"on a weekend, after hours", time.Date(2025, 10, 4, 12, 0, 0, 0, time.UTC), false},      // Saturday 12:00 (exclusive)
		{"on a closed day", time.Date(2025, 10, 5, 10, 0, 0, 0, time.UTC), false},                // Sunday 10:00
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Equal(tc.expected, s.bh.IsOpen(tc.time))
		})
	}
}

func (s *BusinessHoursSuite) TestBusinessHours_JSON() {
	s.Run("should marshal and unmarshal correctly", func() {
		data, err := json.Marshal(s.bh)
		s.Require().NoError(err)

		var unmarshaledBH wisp.BusinessHours
		err = json.Unmarshal(data, &unmarshaledBH)
		s.Require().NoError(err)

		s.True(unmarshaledBH.IsOpen(time.Date(2025, 9, 29, 10, 30, 0, 0, time.UTC)))
		s.False(unmarshaledBH.IsOpen(time.Date(2025, 10, 5, 10, 30, 0, 0, time.UTC)))
	})
}

func (s *BusinessHoursSuite) TestBusinessHours_SQL() {
	s.Run("should write to and scan from database representation", func() {
		val, err := s.bh.Value()
		s.Require().NoError(err)
		s.IsType([]byte{}, val)

		var scannedBH wisp.BusinessHours
		err = scannedBH.Scan(val)
		s.Require().NoError(err)

		s.True(scannedBH.IsOpen(time.Date(2025, 9, 29, 10, 30, 0, 0, time.UTC)))
	})
}
