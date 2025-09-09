package atomic_test

import (
	"encoding/json"
	"testing"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/atomic"
)

type CNPJSuite struct {
	suite.Suite
	validCNPJUnmasked  string
	validCNPJFormatted string
}

func (s *CNPJSuite) SetupSuite() {
	s.validCNPJUnmasked = "45543915000181"
	s.validCNPJFormatted = "45.543.915/0001-81"
}

func TestCNPJSuite(t *testing.T) {
	suite.Run(t, new(CNPJSuite))
}

func (s *CNPJSuite) TestNewCNPJ() {
	testCases := []struct {
		name        string
		input       string
		expected    atomic.CNPJ
		expectError bool
	}{
		{name: "should create a valid CNPJ from unmasked string", input: s.validCNPJUnmasked, expected: atomic.CNPJ(s.validCNPJUnmasked)},
		{name: "should create a valid CNPJ from formatted string", input: s.validCNPJFormatted, expected: atomic.CNPJ(s.validCNPJUnmasked)},
		{name: "should create an empty CNPJ from an empty string", input: "", expected: atomic.EmptyCNPJ},
		{name: "should fail for CNPJ with invalid length", input: "1234567890123", expectError: true},
		{name: "should fail for CNPJ with incorrect check digits", input: "11222333000100", expectError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cnpj, err := atomic.NewCNPJ(tc.input)
			if tc.expectError {
				s.Require().Error(err)
				s.Equal(atomic.EmptyCNPJ, cnpj)
				faultErr, ok := err.(*fault.Error)
				s.Require().True(ok, "error should be of type *fault.Error")
				s.Equal(fault.Invalid, faultErr.Code)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expected, cnpj)
			}
		})
	}
}

func (s *CNPJSuite) TestCNPJ_Methods() {
	cnpj, err := atomic.NewCNPJ(s.validCNPJUnmasked)
	s.Require().NoError(err)

	s.Run("IsZero", func() {
		s.False(cnpj.IsZero())
		s.True(atomic.EmptyCNPJ.IsZero())
	})

	s.Run("String", func() {
		s.Equal(s.validCNPJUnmasked, cnpj.String())
	})

	s.Run("Formatted", func() {
		s.Equal(s.validCNPJFormatted, cnpj.Formatted())
		s.Equal("", atomic.EmptyCNPJ.Formatted())
	})
}

func (s *CNPJSuite) TestCNPJ_JSONMarshaling() {
	s.Run("should marshal and unmarshal a valid CNPJ", func() {
		cnpj, _ := atomic.NewCNPJ(s.validCNPJUnmasked)
		data, err := json.Marshal(cnpj)
		s.Require().NoError(err)
		s.Equal(`"`+s.validCNPJUnmasked+`"`, string(data))

		var unmarshaledCNPJ atomic.CNPJ
		err = json.Unmarshal(data, &unmarshaledCNPJ)
		s.Require().NoError(err)
		s.Equal(cnpj, unmarshaledCNPJ)
	})

	s.Run("should fail to unmarshal an invalid CNPJ string", func() {
		var cnpj atomic.CNPJ
		err := json.Unmarshal([]byte(`"00000000000000"`), &cnpj)
		s.Require().Error(err)
	})
}

func (s *CNPJSuite) TestCNPJ_DatabaseInterface() {
	cnpj, _ := atomic.NewCNPJ(s.validCNPJUnmasked)

	s.Run("Value", func() {
		val, err := cnpj.Value()
		s.Require().NoError(err)
		s.Equal(s.validCNPJUnmasked, val)

		nilVal, err := atomic.EmptyCNPJ.Value()
		s.Require().NoError(err)
		s.Nil(nilVal)
	})

	s.Run("Scan", func() {
		s.Run("should scan a valid string", func() {
			var scannedCNPJ atomic.CNPJ
			err := scannedCNPJ.Scan(s.validCNPJUnmasked)
			s.Require().NoError(err)
			s.Equal(cnpj, scannedCNPJ)
		})

		s.Run("should scan nil as EmptyCNPJ", func() {
			var scannedCNPJ atomic.CNPJ
			err := scannedCNPJ.Scan(nil)
			s.Require().NoError(err)
			s.True(scannedCNPJ.IsZero())
		})

		s.Run("should fail to scan an invalid CNPJ string", func() {
			var scannedCNPJ atomic.CNPJ
			err := scannedCNPJ.Scan("11222333000100")
			s.Require().Error(err)
		})
	})
}
