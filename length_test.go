package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type LengthSuite struct {
	suite.Suite
}

func TestLengthSuite(t *testing.T) {
	suite.Run(t, new(LengthSuite))
}

func (s *LengthSuite) TestNewLength() {
	s.Run("should create length from meters", func() {
		l, err := wisp.NewLength(2.5, wisp.Meter)
		s.Require().NoError(err)
		val, _ := l.In(wisp.Centimeter)
		s.InDelta(250, val, 0.001)
	})

	s.Run("should create length from feet", func() {
		// 1 foot = 0.3048 meters
		l, err := wisp.NewLength(1, wisp.Foot)
		s.Require().NoError(err)
		val, _ := l.In(wisp.Centimeter)
		s.InDelta(30.48, val, 0.001)
	})

	s.Run("should fail for negative length", func() {
		_, err := wisp.NewLength(-5, wisp.Meter)
		s.Require().Error(err)
	})
}

func (s *LengthSuite) TestLength_Conversions() {
	l, _ := wisp.NewLength(1, wisp.Meter)

	valCm, _ := l.In(wisp.Centimeter)
	s.InDelta(100, valCm, 0.001)

	valKm, _ := l.In(wisp.Kilometer)
	s.InDelta(0.001, valKm, 0.001)

	valIn, _ := l.In(wisp.Inch)
	s.InDelta(39.3701, valIn, 0.001)
}

func (s *LengthSuite) TestLength_Arithmetic() {
	l1, _ := wisp.NewLength(50, wisp.Centimeter)
	l2, _ := wisp.NewLength(2, wisp.Meter)

	sum := l1.Add(l2)
	sumInM, _ := sum.In(wisp.Meter)
	s.InDelta(2.5, sumInM, 0.001)

	diff := l2.Subtract(l1)
	diffInM, _ := diff.In(wisp.Meter)
	s.InDelta(1.5, diffInM, 0.001)
}

func (s *LengthSuite) TestLength_JSON_SQL() {
	l, _ := wisp.NewLength(3.5, wisp.Meter)

	s.Run("JSON Marshaling and Unmarshaling", func() {
		data, err := json.Marshal(l)
		s.Require().NoError(err)
		s.JSONEq(`{"value": 3.5, "unit": "m"}`, string(data))

		var unmarshaledL wisp.Length
		err = json.Unmarshal(data, &unmarshaledL)
		s.Require().NoError(err)
		s.True(l.Equals(unmarshaledL))
	})

	s.Run("SQL Interface", func() {
		val, err := l.Value()
		s.Require().NoError(err)
		s.Equal(int64(3500000), val) // 3.5m = 3,500,000 micrometers

		var scannedL wisp.Length
		err = scannedL.Scan(int64(500000)) // 500,000 micrometers = 0.5m
		s.Require().NoError(err)

		m, _ := scannedL.In(wisp.Meter)
		s.InDelta(0.5, m, 0.001)
	})
}
