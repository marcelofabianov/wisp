package atomic_test

import (
	"encoding/json"
	"testing"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/atomic"
)

type CPFSuite struct {
	suite.Suite
	// A valid CPF for testing purposes.
	validCPFUnmasked  string
	validCPFFormatted string
}

func (s *CPFSuite) SetupSuite() {
	s.validCPFUnmasked = "86222616038"
	s.validCPFFormatted = "862.226.160-38"
}

func TestCPFSuite(t *testing.T) {
	suite.Run(t, new(CPFSuite))
}

func (s *CPFSuite) TestNewCPF() {
	testCases := []struct {
		name        string
		input       string
		expected    atomic.CPF
		expectError bool
	}{
		{name: "should create a valid CPF from unmasked string", input: s.validCPFUnmasked, expected: atomic.CPF(s.validCPFUnmasked)},
		{name: "should create a valid CPF from formatted string", input: s.validCPFFormatted, expected: atomic.CPF(s.validCPFUnmasked)},
		{name: "should create an empty CPF from an empty string", input: "", expected: atomic.EmptyCPF},
		{name: "should fail for CPF with invalid length", input: "123456789", expectError: true},
		{name: "should fail for CPF with all repeated digits", input: "11111111111", expectError: true},
		{name: "should fail for CPF with incorrect check digits", input: "12345678900", expectError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cpf, err := atomic.NewCPF(tc.input)
			if tc.expectError {
				s.Require().Error(err)
				s.Equal(atomic.EmptyCPF, cpf)
				faultErr, ok := err.(*fault.Error)
				s.Require().True(ok, "error should be of type *fault.Error")
				s.Equal(fault.Invalid, faultErr.Code)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expected, cpf)
			}
		})
	}
}

func (s *CPFSuite) TestCPF_Methods() {
	cpf, err := atomic.NewCPF(s.validCPFUnmasked)
	s.Require().NoError(err)

	s.Run("IsZero", func() {
		s.False(cpf.IsZero())
		s.True(atomic.EmptyCPF.IsZero())
	})

	s.Run("String", func() {
		s.Equal(s.validCPFUnmasked, cpf.String())
	})

	s.Run("Formatted", func() {
		s.Equal(s.validCPFFormatted, cpf.Formatted())
		s.Equal("", atomic.EmptyCPF.Formatted())
	})
}

func (s *CPFSuite) TestCPF_JSONMarshaling() {
	s.Run("should marshal and unmarshal a valid CPF", func() {
		cpf, _ := atomic.NewCPF(s.validCPFUnmasked)
		data, err := json.Marshal(cpf)
		s.Require().NoError(err)
		s.Equal(`"`+s.validCPFUnmasked+`"`, string(data))

		var unmarshaledCPF atomic.CPF
		err = json.Unmarshal(data, &unmarshaledCPF)
		s.Require().NoError(err)
		s.Equal(cpf, unmarshaledCPF)
	})

	s.Run("should fail to unmarshal an invalid CPF string", func() {
		var cpf atomic.CPF
		err := json.Unmarshal([]byte(`"11111111111"`), &cpf)
		s.Require().Error(err)
	})
}

func (s *CPFSuite) TestCPF_DatabaseInterface() {
	cpf, _ := atomic.NewCPF(s.validCPFUnmasked)

	s.Run("Value", func() {
		val, err := cpf.Value()
		s.Require().NoError(err)
		s.Equal(s.validCPFUnmasked, val)

		nilVal, err := atomic.EmptyCPF.Value()
		s.Require().NoError(err)
		s.Nil(nilVal)
	})

	s.Run("Scan", func() {
		s.Run("should scan a valid string", func() {
			var scannedCPF atomic.CPF
			err := scannedCPF.Scan(s.validCPFUnmasked)
			s.Require().NoError(err)
			s.Equal(cpf, scannedCPF)
		})

		s.Run("should scan nil as EmptyCPF", func() {
			var scannedCPF atomic.CPF
			err := scannedCPF.Scan(nil)
			s.Require().NoError(err)
			s.True(scannedCPF.IsZero())
		})

		s.Run("should fail to scan an invalid CPF string", func() {
			var scannedCPF atomic.CPF
			err := scannedCPF.Scan("99988877766")
			s.Require().Error(err)
		})
	})
}
