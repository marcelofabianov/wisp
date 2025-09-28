package wisp_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type MinValueSuite struct {
	suite.Suite
}

func TestMinValueSuite(t *testing.T) {
	suite.Run(t, new(MinValueSuite))
}

func (s *MinValueSuite) TestNewMinValue() {
	s.Run("should create a valid min value", func() {
		mv, err := wisp.NewMinValue(20, 10)
		s.Require().NoError(err)
		s.Equal(int64(20), mv.Current())
		s.Equal(int64(10), mv.Min())
	})

	s.Run("should fail if current is less than min", func() {
		_, err := wisp.NewMinValue(9, 10)
		s.Require().Error(err)
	})
}

func (s *MinValueSuite) TestMinValue_Add() {
	mv, _ := wisp.NewMinValue(10, 10)

	s.Run("should add amount successfully", func() {
		newMv, err := mv.Add(5)
		s.Require().NoError(err)
		s.Equal(int64(15), newMv.Current())
	})
}

func (s *MinValueSuite) TestMinValue_Subtract() {
	mv, _ := wisp.NewMinValue(20, 10)

	s.Run("should subtract amount successfully", func() {
		newMv, err := mv.Subtract(10)
		s.Require().NoError(err)
		s.Equal(int64(10), newMv.Current())
		s.True(newMv.IsAtMin())
	})

	s.Run("should fail if amount subceeds the minimum", func() {
		_, err := mv.Subtract(11)
		s.Require().Error(err)
		s.Equal(wisp.ErrValueSubceedsMin, err)
	})

	s.Run("should not change the original value (immutability)", func() {
		mv.Subtract(5)
		s.Equal(int64(20), mv.Current())
	})
}

func (s *MinValueSuite) TestMinValue_Set() {
	mv, _ := wisp.NewMinValue(20, 10)

	s.Run("should set a new value within the range", func() {
		newMv, err := mv.Set(15)
		s.Require().NoError(err)
		s.Equal(int64(15), newMv.Current())
	})

	s.Run("should fail to set a value below the minimum", func() {
		_, err := mv.Set(9)
		s.Require().Error(err)
	})
}
