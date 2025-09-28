package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type DayOfWeekSuite struct {
	suite.Suite
}

func TestDayOfWeekSuite(t *testing.T) {
	suite.Run(t, new(DayOfWeekSuite))
}

func (s *DayOfWeekSuite) TestParseDayOfWeek() {
	s.Run("should parse valid day strings", func() {
		d, err := wisp.ParseDayOfWeek("monday")
		s.Require().NoError(err)
		s.Equal(wisp.Monday, d)

		d, err = wisp.ParseDayOfWeek("  SATURDAY  ")
		s.Require().NoError(err)
		s.Equal(wisp.Saturday, d)
	})

	s.Run("should fail for invalid day string", func() {
		_, err := wisp.ParseDayOfWeek("sonday")
		s.Require().Error(err)
	})
}

func (s *DayOfWeekSuite) TestDayOfWeek_Methods() {
	s.Run("IsWeekend", func() {
		s.True(wisp.Saturday.IsWeekend())
		s.True(wisp.Sunday.IsWeekend())
		s.False(wisp.Monday.IsWeekend())
	})

	s.Run("IsWeekday", func() {
		s.False(wisp.Saturday.IsWeekday())
		s.False(wisp.Sunday.IsWeekday())
		s.True(wisp.Monday.IsWeekday())
		s.True(wisp.Friday.IsWeekday())
	})

	s.Run("String", func() {
		s.Equal("Wednesday", wisp.Wednesday.String())
	})
}

func (s *DayOfWeekSuite) TestDayOfWeek_JSON() {
	s.Run("should marshal to lowercase string", func() {
		data, err := json.Marshal(wisp.Tuesday)
		s.Require().NoError(err)
		s.Equal(`"tuesday"`, string(data))
	})

	s.Run("should unmarshal from string", func() {
		var d wisp.DayOfWeek
		err := json.Unmarshal([]byte(`"friday"`), &d)
		s.Require().NoError(err)
		s.Equal(wisp.Friday, d)
	})
}

func (s *DayOfWeekSuite) TestDayOfWeek_SQL() {
	s.Run("Value", func() {
		val, err := wisp.Thursday.Value()
		s.Require().NoError(err)
		s.Equal(int64(4), val) // time.Thursday is 4
	})

	s.Run("Scan", func() {
		var d wisp.DayOfWeek
		err := d.Scan(int64(1)) // 1 is Monday
		s.Require().NoError(err)
		s.Equal(wisp.Monday, d)

		err = d.Scan(int64(7)) // out of range
		s.Require().Error(err)
	})
}
