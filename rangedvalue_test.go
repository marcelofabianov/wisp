package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type RangedValueSuite struct {
	suite.Suite
}

func TestRangedValueSuite(t *testing.T) {
	suite.Run(t, new(RangedValueSuite))
}

func (s *RangedValueSuite) TestNewRangedValue() {
	s.Run("should create a valid ranged value", func() {
		rv, err := wisp.NewRangedValue(50, 10, 100)
		s.Require().NoError(err)
		s.Equal(int64(50), rv.Current())
		s.Equal(int64(10), rv.Min())
		s.Equal(int64(100), rv.Max())
	})

	s.Run("should create a valid value at the boundaries", func() {
		rv, err := wisp.NewRangedValue(10, 10, 100)
		s.Require().NoError(err)
		s.True(rv.IsAtMin())

		rv, err = wisp.NewRangedValue(100, 10, 100)
		s.Require().NoError(err)
		s.True(rv.IsAtMax())
	})

	s.Run("should fail if current is outside the range", func() {
		_, err := wisp.NewRangedValue(9, 10, 100)
		s.Require().Error(err)

		_, err = wisp.NewRangedValue(101, 10, 100)
		s.Require().Error(err)
	})

	s.Run("should fail if min is greater than max", func() {
		_, err := wisp.NewRangedValue(15, 20, 10)
		s.Require().Error(err)
	})
}

func (s *RangedValueSuite) TestRangedValue_Add() {
	rv, _ := wisp.NewRangedValue(90, 0, 100)

	s.Run("should add amount successfully", func() {
		newRv, err := rv.Add(10)
		s.Require().NoError(err)
		s.Equal(int64(100), newRv.Current())
		s.True(newRv.IsAtMax())
	})

	s.Run("should fail if amount exceeds the maximum", func() {
		_, err := rv.Add(11)
		s.Require().Error(err)
		s.Equal(wisp.ErrValueExceedsMax, err)
	})

	s.Run("should not change the original value (immutability)", func() {
		rv.Add(5)
		s.Equal(int64(90), rv.Current())
	})
}

func (s *RangedValueSuite) TestRangedValue_Subtract() {
	rv, _ := wisp.NewRangedValue(10, 5, 100)

	s.Run("should subtract amount successfully", func() {
		newRv, err := rv.Subtract(5)
		s.Require().NoError(err)
		s.Equal(int64(5), newRv.Current())
		s.True(newRv.IsAtMin())
	})

	s.Run("should fail if amount subceeds the minimum", func() {
		_, err := rv.Subtract(6)
		s.Require().Error(err)
		s.Equal(wisp.ErrValueSubceedsMin, err)
	})
}

func (s *RangedValueSuite) TestRangedValue_Set() {
	rv, _ := wisp.NewRangedValue(50, 10, 100)

	s.Run("should set a new value within the range", func() {
		newRv, err := rv.Set(75)
		s.Require().NoError(err)
		s.Equal(int64(75), newRv.Current())
	})

	s.Run("should not change the original value (immutability)", func() {
		rv.Set(80)
		s.Equal(int64(50), rv.Current())
	})

	s.Run("should fail to set a value outside the range", func() {
		_, err := rv.Set(101)
		s.Require().Error(err)

		_, err = rv.Set(9)
		s.Require().Error(err)
	})
}

func (s *RangedValueSuite) TestRangedValue_JSON() {
	s.Run("should marshal and unmarshal correctly", func() {
		rv, _ := wisp.NewRangedValue(50, 0, 1000)
		data, err := json.Marshal(rv)
		s.Require().NoError(err)
		s.JSONEq(`{"current": 50, "min": 0, "max": 1000}`, string(data))

		var unmarshaledRv wisp.RangedValue
		err = json.Unmarshal(data, &unmarshaledRv)
		s.Require().NoError(err)
		s.Equal(rv, unmarshaledRv)
	})

	s.Run("should fail to unmarshal invalid range", func() {
		invalidJSON := `{"current": 5, "min": 10, "max": 100}`
		var unmarshaledRv wisp.RangedValue
		err := json.Unmarshal([]byte(invalidJSON), &unmarshaledRv)
		s.Require().Error(err)
	})
}
