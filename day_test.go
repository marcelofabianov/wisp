package atomic_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/atomic"
)

type DaySuite struct {
	suite.Suite
}

func TestDaySuite(t *testing.T) {
	suite.Run(t, new(DaySuite))
}

func (s *DaySuite) TestNewDay() {
	testCases := []struct {
		name        string
		input       int
		expected    atomic.Day
		expectError bool
	}{
		{name: "should create day for lower bound", input: 1, expected: 1, expectError: false},
		{name: "should create day for mid-range", input: 15, expected: 15, expectError: false},
		{name: "should create day for upper bound", input: 31, expected: 31, expectError: false},
		{name: "should fail for value below lower bound", input: 0, expectError: true},
		{name: "should fail for negative value", input: -5, expectError: true},
		{name: "should fail for value above upper bound", input: 32, expectError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			day, err := atomic.NewDay(tc.input)
			if tc.expectError {
				s.Require().Error(err)
				s.Equal(atomic.ZeroDay, day)
				faultErr, ok := err.(*fault.Error)
				s.Require().True(ok)
				s.Equal(fault.Invalid, faultErr.Code)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expected, day)
			}
		})
	}
}

func (s *DaySuite) TestDay_ZeroStateAndInt() {
	s.Run("should correctly identify zero and non-zero states", func() {
		day10, _ := atomic.NewDay(10)
		s.False(day10.IsZero())
		s.Equal(10, day10.Int())

		s.True(atomic.ZeroDay.IsZero())
		s.Equal(0, atomic.ZeroDay.Int())
	})
}

func (s *DaySuite) TestDay_DateCalculations() {
	today := time.Date(2025, time.September, 15, 10, 0, 0, 0, time.UTC)

	s.Run("HasPassed", func() {
		dayBefore, _ := atomic.NewDay(10)
		daySame, _ := atomic.NewDay(15)
		dayAfter, _ := atomic.NewDay(20)

		s.True(dayBefore.HasPassed(today))
		s.False(daySame.HasPassed(today))
		s.False(dayAfter.HasPassed(today))
		s.False(atomic.ZeroDay.HasPassed(today))
	})

	s.Run("DaysUntil", func() {
		dayBefore, _ := atomic.NewDay(10)
		daySame, _ := atomic.NewDay(15)
		dayAfter, _ := atomic.NewDay(25)

		// days left in Sept (30-15=15) + 10 = 25
		s.Equal(25, dayBefore.DaysUntil(today))
		s.Equal(0, daySame.DaysUntil(today))
		s.Equal(10, dayAfter.DaysUntil(today))
		s.Equal(0, atomic.ZeroDay.DaysUntil(today))
	})

	s.Run("DaysOverdue", func() {
		dayBefore, _ := atomic.NewDay(10)
		daySame, _ := atomic.NewDay(15)
		dayAfter, _ := atomic.NewDay(25)

		// days in previous month (Aug=31). (31-25) + 15 = 21
		s.Equal(21, dayAfter.DaysOverdue(today))
		s.Equal(0, daySame.DaysOverdue(today))
		s.Equal(5, dayBefore.DaysOverdue(today))
		s.Equal(0, atomic.ZeroDay.DaysOverdue(today))
	})
}

func (s *DaySuite) TestDay_JSONMarshaling() {
	s.Run("should marshal and unmarshal a valid day", func() {
		day, _ := atomic.NewDay(28)
		data, err := json.Marshal(day)
		s.Require().NoError(err)
		s.Equal(`28`, string(data))

		var unmarshaledDay atomic.Day
		err = json.Unmarshal(data, &unmarshaledDay)
		s.Require().NoError(err)
		s.Equal(day, unmarshaledDay)
	})

	s.Run("should unmarshal null as ZeroDay", func() {
		var day atomic.Day
		err := json.Unmarshal([]byte("null"), &day)
		s.Require().NoError(err)
		s.True(day.IsZero())
	})

	s.Run("should fail to unmarshal an invalid day value", func() {
		var day atomic.Day
		err := json.Unmarshal([]byte(`40`), &day)
		s.Require().Error(err)
		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok)
		s.Equal(fault.Invalid, faultErr.Code)
	})
}

func (s *DaySuite) TestDay_DatabaseInterface() {
	s.Run("Value", func() {
		day, _ := atomic.NewDay(20)
		val, err := day.Value()
		s.Require().NoError(err)
		s.Equal(int64(20), val)

		nilVal, err := atomic.ZeroDay.Value()
		s.Require().NoError(err)
		s.Nil(nilVal)
	})

	s.Run("Scan", func() {
		testCases := []struct {
			name        string
			src         interface{}
			expected    atomic.Day
			expectError bool
		}{
			{name: "should scan a valid int64", src: int64(18), expected: atomic.Day(18)},
			{name: "should scan nil as ZeroDay", src: nil, expected: atomic.ZeroDay},
			{name: "should fail to scan an out-of-bounds int64", src: int64(32), expectError: true},
			{name: "should fail to scan zero", src: int64(0), expectError: true},
			{name: "should fail to scan an incompatible type", src: "25", expectError: true},
		}

		for _, tc := range testCases {
			s.Run(tc.name, func() {
				var day atomic.Day
				err := day.Scan(tc.src)

				if tc.expectError {
					s.Require().Error(err)
					faultErr, ok := err.(*fault.Error)
					s.Require().True(ok)
					s.Equal(fault.Invalid, faultErr.Code)
				} else {
					s.Require().NoError(err)
					s.Equal(tc.expected, day)
				}
			})
		}
	})
}
