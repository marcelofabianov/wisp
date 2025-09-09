package wisp_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/marcelofabianov/fault"
	wisp "github.com/marcelofabianov/wisp"
	"github.com/stretchr/testify/suite"
)

type DateSuite struct {
	suite.Suite
}

func TestDateSuite(t *testing.T) {
	suite.Run(t, new(DateSuite))
}

func (s *DateSuite) TestNewDateAndParse() {
	s.Run("should create a valid date", func() {
		d, err := wisp.NewDate(2025, time.September, 9)
		s.Require().NoError(err)
		s.Equal(2025, d.Year())
		s.Equal(time.September, d.Month())
		s.Equal(9, d.Day())
		s.False(d.IsZero())
	})

	s.Run("should fail for an invalid date", func() {
		_, err := wisp.NewDate(2025, time.February, 30)
		s.Require().Error(err)
		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok)
		s.Equal(fault.Invalid, faultErr.Code)
	})

	s.Run("should parse a valid date string", func() {
		d, err := wisp.ParseDate("2025-09-09")
		s.Require().NoError(err)
		expected, _ := wisp.NewDate(2025, time.September, 9)
		s.True(d.Equals(expected))
	})

	s.Run("should fail to parse an invalid date string", func() {
		_, err := wisp.ParseDate("09-09-2025")
		s.Require().Error(err)
	})
}

func (s *DateSuite) TestDate_ComparisonAndManipulation() {
	d1, _ := wisp.NewDate(2025, time.January, 10)
	d2, _ := wisp.NewDate(2025, time.January, 20)
	d1Clone, _ := wisp.NewDate(2025, time.January, 10)

	s.Run("Comparison", func() {
		s.True(d1.Equals(d1Clone))
		s.False(d1.Equals(d2))
		s.True(d1.Before(d2))
		s.True(d2.After(d1))
	})

	s.Run("Manipulation", func() {
		s.Equal("2025-01-20", d1.AddDays(10).String())
		s.Equal("2025-03-10", d1.AddMonths(2).String())
		s.Equal("2030-01-10", d1.AddYears(5).String())
	})
}

func (s *DateSuite) TestDate_JSONMarshaling() {
	d, _ := wisp.NewDate(2025, time.September, 9)

	s.Run("should marshal and unmarshal correctly", func() {
		data, err := json.Marshal(d)
		s.Require().NoError(err)
		s.Equal(`"2025-09-09"`, string(data))

		var unmarshaledDate wisp.Date
		err = json.Unmarshal(data, &unmarshaledDate)
		s.Require().NoError(err)
		s.True(d.Equals(unmarshaledDate))
	})

	s.Run("should handle zero and null values", func() {
		data, err := json.Marshal(wisp.ZeroDate)
		s.Require().NoError(err)
		s.Equal("null", string(data))

		var unmarshaledDate wisp.Date
		err = json.Unmarshal([]byte("null"), &unmarshaledDate)
		s.Require().NoError(err)
		s.True(unmarshaledDate.IsZero())
	})
}

func (s *DateSuite) TestDate_DatabaseInterface() {
	d, _ := wisp.NewDate(2025, time.September, 9)

	s.Run("Value", func() {
		val, err := d.Value()
		s.Require().NoError(err)

		// The driver expects a time.Time for a DATE column
		t, ok := val.(time.Time)
		s.Require().True(ok)
		s.Equal(d.Year(), t.Year())
		s.Equal(d.Month(), t.Month())
		s.Equal(d.Day(), t.Day())
	})

	s.Run("Scan", func() {
		var scannedDate wisp.Date
		// Simulate the database returning a time.Time
		dbTime := time.Date(2025, 9, 9, 15, 30, 0, 0, time.Local)
		err := scannedDate.Scan(dbTime)
		s.Require().NoError(err)
		s.True(d.Equals(scannedDate), "Scan should truncate the time part")
	})

	s.Run("should handle nil from database", func() {
		var scannedDate wisp.Date
		err := scannedDate.Scan(nil)
		s.Require().NoError(err)
		s.True(scannedDate.IsZero())
	})
}
