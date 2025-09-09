package atomic_test

import (
	"encoding/json"
	"testing"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/atomic"
)

type UFSuite struct {
	suite.Suite
}

func TestUFSuite(t *testing.T) {
	suite.Run(t, new(UFSuite))
}

func (s *UFSuite) TestNewUF() {
	testCases := []struct {
		name        string
		input       string
		expected    atomic.UF
		expectError bool
	}{
		{name: "should create a valid UF from uppercase string", input: "SP", expected: "SP"},
		{name: "should create and normalize a lowercase UF", input: "go", expected: "GO"},
		{name: "should create and normalize a UF with spaces", input: "  rj  ", expected: "RJ"},
		{name: "should create an empty UF from an empty string", input: "", expected: atomic.EmptyUF},
		{name: "should fail for an invalid code", input: "XX", expectError: true},
		{name: "should fail for a string with more than 2 letters", input: "ABC", expectError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			uf, err := atomic.NewUF(tc.input)
			if tc.expectError {
				s.Require().Error(err)
				s.Equal(atomic.EmptyUF, uf)
				faultErr, ok := err.(*fault.Error)
				s.Require().True(ok)
				s.Equal(fault.Invalid, faultErr.Code)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expected, uf)
			}
		})
	}
}

func (s *UFSuite) TestUF_IsValidAndZero() {
	sp, _ := atomic.NewUF("SP")
	xx, _ := atomic.NewUF("XX")

	s.True(sp.IsValid())
	s.False(xx.IsValid())
	s.False(atomic.EmptyUF.IsValid())

	s.False(sp.IsZero())
	s.True(atomic.EmptyUF.IsZero())
}

func (s *UFSuite) TestUF_JSONMarshaling() {
	s.Run("should marshal and unmarshal a valid UF", func() {
		uf, _ := atomic.NewUF("MG")
		data, err := json.Marshal(uf)
		s.Require().NoError(err)
		s.Equal(`"MG"`, string(data))

		var unmarshaledUF atomic.UF
		err = json.Unmarshal(data, &unmarshaledUF)
		s.Require().NoError(err)
		s.Equal(uf, unmarshaledUF)
	})

	s.Run("should unmarshal null as EmptyUF", func() {
		var uf atomic.UF
		err := json.Unmarshal([]byte("null"), &uf)
		s.Require().NoError(err)
		s.True(uf.IsZero())
	})

	s.Run("should fail to unmarshal an invalid UF string", func() {
		var uf atomic.UF
		err := json.Unmarshal([]byte(`"ZZ"`), &uf)
		s.Require().Error(err)
	})
}

func (s *UFSuite) TestUF_DatabaseInterface() {
	uf, _ := atomic.NewUF("BA")

	s.Run("Value", func() {
		val, err := uf.Value()
		s.Require().NoError(err)
		s.Equal("BA", val)

		nilVal, err := atomic.EmptyUF.Value()
		s.Require().NoError(err)
		s.Nil(nilVal)
	})

	s.Run("Scan", func() {
		s.Run("should scan a valid string", func() {
			var scannedUF atomic.UF
			err := scannedUF.Scan("SC")
			s.Require().NoError(err)
			s.Equal(atomic.UF("SC"), scannedUF)
		})

		s.Run("should scan nil as EmptyUF", func() {
			var scannedUF atomic.UF
			err := scannedUF.Scan(nil)
			s.Require().NoError(err)
			s.True(scannedUF.IsZero())
		})

		s.Run("should fail to scan an invalid UF string", func() {
			var scannedUF atomic.UF
			err := scannedUF.Scan("XY")
			s.Require().Error(err)
		})
	})
}
