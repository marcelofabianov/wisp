package atomic_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/atomic"
)

type DateRangeSuite struct {
	suite.Suite
}

func TestDateRangeSuite(t *testing.T) {
	suite.Run(t, new(DateRangeSuite))
}

func (s *DateRangeSuite) TestNewDateRange() {
	start, _ := atomic.NewDate(2025, time.September, 10)
	end, _ := atomic.NewDate(2025, time.September, 20)

	s.Run("should create a valid date range", func() {
		dr, err := atomic.NewDateRange(start, end)
		s.Require().NoError(err)
		s.True(start.Equals(dr.Start()))
		s.True(end.Equals(dr.End()))
	})

	s.Run("should create a valid single-day range", func() {
		dr, err := atomic.NewDateRange(start, start)
		s.Require().NoError(err)
		s.Equal(1, dr.Days())
	})

	s.Run("should fail if start date is after end date", func() {
		_, err := atomic.NewDateRange(end, start)
		s.Require().Error(err)
	})
}

func (s *DateRangeSuite) TestDateRange_Methods() {
	start, _ := atomic.NewDate(2025, 9, 10)
	end, _ := atomic.NewDate(2025, 9, 20)
	dr, _ := atomic.NewDateRange(start, end)

	s.Run("Contains", func() {
		inside, _ := atomic.NewDate(2025, 9, 15)
		before, _ := atomic.NewDate(2025, 9, 9)
		s.True(dr.Contains(inside))
		s.True(dr.Contains(start))
		s.True(dr.Contains(end))
		s.False(dr.Contains(before))
	})

	s.Run("Overlaps", func() {
		mustNewRange := func(y1, m1, d1, y2, m2, d2 int) atomic.DateRange {
			dStart, err := atomic.NewDate(y1, time.Month(m1), d1)
			s.Require().NoError(err)
			dEnd, err := atomic.NewDate(y2, time.Month(m2), d2)
			s.Require().NoError(err)
			rng, err := atomic.NewDateRange(dStart, dEnd)
			s.Require().NoError(err)
			return rng
		}

		// Overlaps end
		r1 := mustNewRange(2025, 9, 15, 2025, 9, 25)
		// Overlaps start
		r2 := mustNewRange(2025, 9, 1, 2025, 9, 15)
		// No overlap
		r3 := mustNewRange(2025, 9, 1, 2025, 9, 9)
		// Fully contained
		r4 := mustNewRange(2025, 9, 12, 2025, 9, 18)

		s.True(dr.Overlaps(r1))
		s.True(dr.Overlaps(r2))
		s.False(dr.Overlaps(r3))
		s.True(dr.Overlaps(r4))
	})

	s.Run("Days", func() {
		s.Equal(11, dr.Days()) // 10, 11, ..., 20 -> 11 days
		oneDayRange, _ := atomic.NewDateRange(start, start)
		s.Equal(1, oneDayRange.Days())
	})
}

func (s *DateRangeSuite) TestDateRange_JSONMarshaling() {
	start, _ := atomic.NewDate(2025, 9, 10)
	end, _ := atomic.NewDate(2025, 9, 20)
	dr, _ := atomic.NewDateRange(start, end)

	s.Run("should marshal and unmarshal correctly", func() {
		data, err := json.Marshal(dr)
		s.Require().NoError(err)
		s.JSONEq(`{"start": "2025-09-10", "end": "2025-09-20"}`, string(data))

		var unmarshaledDR atomic.DateRange
		err = json.Unmarshal(data, &unmarshaledDR)
		s.Require().NoError(err)
		s.True(dr.Equals(unmarshaledDR))
	})

	s.Run("should fail to unmarshal an invalid range", func() {
		invalidJSON := `{"start": "2025-09-20", "end": "2025-09-10"}`
		var dr atomic.DateRange
		err := json.Unmarshal([]byte(invalidJSON), &dr)
		s.Require().Error(err)
	})
}
