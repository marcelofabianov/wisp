package atomic_test

import (
	"encoding/json"
	"testing"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/atomic"
)

type CEPSuite struct {
	suite.Suite
}

func TestCEPSuite(t *testing.T) {
	suite.Run(t, new(CEPSuite))
}

func (s *CEPSuite) TestNewCEP() {
	testCases := []struct {
		name        string
		input       string
		expected    atomic.CEP
		expectError bool
	}{
		{name: "should create a valid CEP from unmasked string", input: "74835030", expected: "74835030"},
		{name: "should create a valid CEP from formatted string", input: "74835-030", expected: "74835030"},
		{name: "should create an empty CEP from an empty string", input: "", expected: atomic.EmptyCEP},
		{name: "should fail for CEP with less than 8 digits", input: "7483503", expectError: true},
		{name: "should fail for CEP with more than 8 digits", input: "748350301", expectError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cep, err := atomic.NewCEP(tc.input)
			if tc.expectError {
				s.Require().Error(err)
				s.Equal(atomic.EmptyCEP, cep)
				faultErr, ok := err.(*fault.Error)
				s.Require().True(ok)
				s.Equal(fault.Invalid, faultErr.Code)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expected, cep)
			}
		})
	}
}

func (s *CEPSuite) TestCEP_Methods() {
	cep, _ := atomic.NewCEP("74835030")

	s.Run("IsZero", func() {
		s.False(cep.IsZero())
		s.True(atomic.EmptyCEP.IsZero())
	})

	s.Run("String", func() {
		s.Equal("74835030", cep.String())
	})

	s.Run("Formatted", func() {
		s.Equal("74835-030", cep.Formatted())
		s.Equal("", atomic.EmptyCEP.Formatted())
	})
}

func (s *CEPSuite) TestCEP_JSONMarshaling() {
	s.Run("should marshal and unmarshal a valid CEP", func() {
		cep, _ := atomic.NewCEP("74835-030")
		data, err := json.Marshal(cep)
		s.Require().NoError(err)
		s.Equal(`"74835030"`, string(data))

		var unmarshaledCEP atomic.CEP
		err = json.Unmarshal(data, &unmarshaledCEP)
		s.Require().NoError(err)
		s.Equal(cep, unmarshaledCEP)
	})

	s.Run("should fail to unmarshal an invalid CEP string", func() {
		var cep atomic.CEP
		err := json.Unmarshal([]byte(`"12345"`), &cep)
		s.Require().Error(err)
	})
}

func (s *CEPSuite) TestCEP_DatabaseInterface() {
	cep, _ := atomic.NewCEP("74835030")

	s.Run("Value", func() {
		val, err := cep.Value()
		s.Require().NoError(err)
		s.Equal("74835030", val)

		nilVal, err := atomic.EmptyCEP.Value()
		s.Require().NoError(err)
		s.Nil(nilVal)
	})

	s.Run("Scan", func() {
		s.Run("should scan a valid string", func() {
			var scannedCEP atomic.CEP
			err := scannedCEP.Scan("12345678")
			s.Require().NoError(err)
			s.Equal(atomic.CEP("12345678"), scannedCEP)
		})

		s.Run("should scan nil as EmptyCEP", func() {
			var scannedCEP atomic.CEP
			err := scannedCEP.Scan(nil)
			s.Require().NoError(err)
			s.True(scannedCEP.IsZero())
		})

		s.Run("should fail to scan an invalid CEP string", func() {
			var scannedCEP atomic.CEP
			err := scannedCEP.Scan("123")
			s.Require().Error(err)
		})
	})
}
