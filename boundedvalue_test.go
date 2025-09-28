package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type BoundedValueSuite struct {
	suite.Suite
}

func TestBoundedValueSuite(t *testing.T) {
	suite.Run(t, new(BoundedValueSuite))
}

func (s *BoundedValueSuite) TestNewBoundedValue() {
	s.Run("should create a valid bounded value", func() {
		bv, err := wisp.NewBoundedValue(50, 100)
		s.Require().NoError(err)
		s.Equal(int64(50), bv.Current())
		s.Equal(int64(100), bv.Max())
		s.Equal(int64(50), bv.AvailableSpace())
	})

	s.Run("should fail if current is greater than max", func() {
		_, err := wisp.NewBoundedValue(101, 100)
		s.Require().Error(err)
	})

	s.Run("should fail if current is negative", func() {
		_, err := wisp.NewBoundedValue(-1, 100)
		s.Require().Error(err)
	})

	s.Run("should fail if max is negative", func() {
		_, err := wisp.NewBoundedValue(0, -100)
		s.Require().Error(err)
	})
}

func (s *BoundedValueSuite) TestBoundedValue_Add() {
	bv, _ := wisp.NewBoundedValue(90, 100)

	s.Run("should add amount successfully", func() {
		newBv, err := bv.Add(10)
		s.Require().NoError(err)
		s.Equal(int64(100), newBv.Current())
		s.True(newBv.IsFull())
	})

	s.Run("should fail if amount exceeds the maximum", func() {
		_, err := bv.Add(11)
		s.Require().Error(err)
		s.Equal(wisp.ErrValueExceedsMax, err)
	})

	s.Run("should not change the original value (immutability)", func() {
		bv.Add(5)
		s.Equal(int64(90), bv.Current())
	})

	s.Run("should fail for negative amount", func() {
		_, err := bv.Add(-5)
		s.Require().Error(err)
	})
}

func (s *BoundedValueSuite) TestBoundedValue_Subtract() {
	bv, _ := wisp.NewBoundedValue(10, 100)

	s.Run("should subtract amount successfully", func() {
		newBv, err := bv.Subtract(5)
		s.Require().NoError(err)
		s.Equal(int64(5), newBv.Current())
	})

	s.Run("should fail if subtracting more than current value", func() {
		_, err := bv.Subtract(11)
		s.Require().Error(err)
	})

	s.Run("should fail for negative amount", func() {
		_, err := bv.Subtract(-5)
		s.Require().Error(err)
	})
}

func (s *BoundedValueSuite) TestBoundedValue_Set() {
	bv, _ := wisp.NewBoundedValue(50, 100)

	s.Run("should set a new value within the range", func() {
		newBv, err := bv.Set(75)
		s.Require().NoError(err)
		s.Equal(int64(75), newBv.Current())
	})

	s.Run("should fail to set a value outside the [0, max] range", func() {
		_, err := bv.Set(101)
		s.Require().Error(err)

		_, err = bv.Set(-1)
		s.Require().Error(err)
	})
}

func (s *BoundedValueSuite) TestBoundedValue_JSON() {
	s.Run("should marshal and unmarshal correctly", func() {
		bv, _ := wisp.NewBoundedValue(25, 200)
		data, err := json.Marshal(bv)
		s.Require().NoError(err)
		s.JSONEq(`{"current": 25, "max": 200}`, string(data))

		var unmarshaledBv wisp.BoundedValue
		err = json.Unmarshal(data, &unmarshaledBv)
		s.Require().NoError(err)
		s.Equal(bv.Current(), unmarshaledBv.Current())
		s.Equal(bv.Max(), unmarshaledBv.Max())
	})
}
