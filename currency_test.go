package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/suite"

	wisp "github.com/marcelofabianov/wisp"
)

type CurrencySuite struct {
	suite.Suite
}

func TestCurrencySuite(t *testing.T) {
	suite.Run(t, new(CurrencySuite))
}

func (s *CurrencySuite) TestNewCurrency() {
	testCases := []struct {
		name        string
		input       string
		expected    wisp.Currency
		expectError bool
	}{
		{name: "should create a valid uppercase currency", input: "BRL", expected: wisp.BRL},
		{name: "should create and normalize a lowercase currency", input: "usd", expected: wisp.USD},
		{name: "should create and normalize a mixed-case currency with spaces", input: "  eUr  ", expected: wisp.EUR},
		{name: "should handle empty string as EmptyCurrency", input: "", expected: wisp.EmptyCurrency},
		{name: "should handle blank string as EmptyCurrency", input: "   ", expected: wisp.EmptyCurrency},
		{name: "should fail for an unsupported currency code", input: "JPY", expectError: true},
		{name: "should fail for an invalid string", input: "invalid", expectError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			curr, err := wisp.NewCurrency(tc.input)
			if tc.expectError {
				s.Require().Error(err)
				s.Equal(wisp.EmptyCurrency, curr)
				faultErr, ok := err.(*fault.Error)
				s.Require().True(ok)
				s.Equal(fault.Invalid, faultErr.Code)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expected, curr)
			}
		})
	}
}

func (s *CurrencySuite) TestCurrency_IsValidAndZero() {
	s.Run("should correctly validate currency codes", func() {
		s.True(wisp.BRL.IsValid())
		s.True(wisp.USD.IsValid())
		s.True(wisp.EUR.IsValid())
		s.False(wisp.EmptyCurrency.IsValid())
		s.False(wisp.Currency("XYZ").IsValid())
	})

	s.Run("should correctly identify zero state", func() {
		s.True(wisp.EmptyCurrency.IsZero())
		s.False(wisp.BRL.IsZero())
	})
}

func (s *CurrencySuite) TestCurrency_JSONMarshaling() {
	s.Run("should marshal and unmarshal a valid currency", func() {
		curr := wisp.USD
		data, err := json.Marshal(curr)
		s.Require().NoError(err)
		s.Equal(`"USD"`, string(data))

		var unmarshaledCurr wisp.Currency
		err = json.Unmarshal(data, &unmarshaledCurr)
		s.Require().NoError(err)
		s.Equal(curr, unmarshaledCurr)
	})

	s.Run("should unmarshal null as EmptyCurrency", func() {
		var curr wisp.Currency
		err := json.Unmarshal([]byte("null"), &curr)
		s.Require().NoError(err)
		s.True(curr.IsZero())
	})

	s.Run("should unmarshal empty json string as EmptyCurrency", func() {
		var curr wisp.Currency
		err := json.Unmarshal([]byte(`""`), &curr)
		s.Require().NoError(err)
		s.True(curr.IsZero())
	})

	s.Run("should fail to unmarshal an invalid currency code", func() {
		var curr wisp.Currency
		err := json.Unmarshal([]byte(`"XYZ"`), &curr)
		s.Require().Error(err)
		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok)
		s.Equal(fault.Invalid, faultErr.Code)
	})
}

func (s *CurrencySuite) TestCurrency_DatabaseInterface() {
	s.Run("Value", func() {
		val, err := wisp.EUR.Value()
		s.Require().NoError(err)
		s.Equal("EUR", val)

		nilVal, err := wisp.EmptyCurrency.Value()
		s.Require().NoError(err)
		s.Nil(nilVal)
	})

	s.Run("Scan", func() {
		testCases := []struct {
			name        string
			src         interface{}
			expected    wisp.Currency
			expectError bool
		}{
			{name: "should scan a valid string", src: "BRL", expected: wisp.BRL},
			{name: "should scan and normalize a lowercase string", src: "usd", expected: wisp.USD},
			{name: "should scan a valid byte slice", src: []byte("EUR"), expected: wisp.EUR},
			{name: "should scan nil as EmptyCurrency", src: nil, expected: wisp.EmptyCurrency},
			{name: "should fail to scan an invalid code", src: "JPY", expectError: true},
			{name: "should fail to scan an incompatible type", src: 123, expectError: true},
		}

		for _, tc := range testCases {
			s.Run(tc.name, func() {
				var curr wisp.Currency
				err := curr.Scan(tc.src)

				if tc.expectError {
					s.Require().Error(err)
					faultErr, ok := err.(*fault.Error)
					s.Require().True(ok)
					s.Equal(fault.Invalid, faultErr.Code)
				} else {
					s.Require().NoError(err)
					s.Equal(tc.expected, curr)
				}
			})
		}
	})
}
