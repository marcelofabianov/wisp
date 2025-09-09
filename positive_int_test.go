package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type PositiveIntSuite struct {
	suite.Suite
}

func TestPositiveIntSuite(t *testing.T) {
	suite.Run(t, new(PositiveIntSuite))
}

func (s *PositiveIntSuite) TestNewPositiveInt() {
	s.Run("should create a valid positive int", func() {
		pi, err := wisp.NewPositiveInt(10)
		s.Require().NoError(err)
		s.Equal(10, pi.Int())
		s.False(pi.IsZero())
	})

	s.Run("should fail for zero", func() {
		_, err := wisp.NewPositiveInt(0)
		s.Require().Error(err)
	})

	s.Run("should fail for a negative number", func() {
		_, err := wisp.NewPositiveInt(-5)
		s.Require().Error(err)
	})
}

func (s *PositiveIntSuite) TestPositiveInt_JSON() {
	pi, _ := wisp.NewPositiveInt(100)

	data, err := json.Marshal(pi)
	s.Require().NoError(err)
	s.Equal("100", string(data))

	var unmarshaledPI wisp.PositiveInt
	err = json.Unmarshal(data, &unmarshaledPI)
	s.Require().NoError(err)
	s.Equal(pi, unmarshaledPI)

	err = json.Unmarshal([]byte("-1"), &unmarshaledPI)
	s.Require().Error(err)
}

func (s *PositiveIntSuite) TestPositiveInt_SQL() {
	pi, _ := wisp.NewPositiveInt(100)

	val, err := pi.Value()
	s.Require().NoError(err)
	s.Equal(int64(100), val)

	var scannedPI wisp.PositiveInt
	err = scannedPI.Scan(int64(50))
	s.Require().NoError(err)
	s.Equal(50, scannedPI.Int())

	err = scannedPI.Scan(int64(0))
	s.Require().Error(err)
}
