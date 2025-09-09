package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/suite"

	wisp "github.com/marcelofabianov/wisp"
)

type PhoneSuite struct {
	suite.Suite
}

func TestPhoneSuite(t *testing.T) {
	suite.Run(t, new(PhoneSuite))
}

func (s *PhoneSuite) TestNewPhone() {
	testCases := []struct {
		name        string
		input       string
		expected    wisp.Phone
		expectError bool
		errCode     fault.Code
	}{
		// Happy Paths
		{name: "should create a valid mobile phone from E.164 format", input: "5562982870053", expected: "5562982870053"},
		{name: "should create a valid mobile phone from formatted string", input: "+55 (62) 98287-0053", expected: "5562982870053"},
		{name: "should create a valid mobile phone assuming country code", input: "62982870053", expected: "5562982870053"},
		{name: "should create a valid landline phone", input: "(11) 4567-1234", expected: "551145671234"},
		{name: "should create an empty phone from an empty string", input: "", expected: wisp.EmptyPhone},

		// Error Paths
		{name: "should fail for number too short", input: "6298287", expectError: true, errCode: fault.Invalid},
		{name: "should fail for number too long", input: "5562982870053123", expectError: true, errCode: fault.Invalid},
		{name: "should fail for invalid DDD", input: "5523982870053", expectError: true, errCode: fault.Invalid},
		{name: "should fail for mobile number not starting with 9", input: "5562882870053", expectError: true, errCode: fault.Invalid},
		{name: "should fail for landline number with invalid prefix", input: "556212345678", expectError: true, errCode: fault.Invalid},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			phone, err := wisp.NewPhone(tc.input)

			if tc.expectError {
				s.Require().Error(err)
				s.Equal(wisp.EmptyPhone, phone)
				faultErr, ok := err.(*fault.Error)
				s.Require().True(ok)
				s.Equal(tc.errCode, faultErr.Code)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expected, phone)
			}
		})
	}
}

func (s *PhoneSuite) TestPhone_ComponentsAndStateChecks() {
	mobile, _ := wisp.NewPhone("5562982870053")
	landline, _ := wisp.NewPhone("551145671234")

	s.Run("Components", func() {
		s.Equal("55", mobile.CountryCode())
		s.Equal("62", mobile.AreaCode())
		s.Equal("982870053", mobile.Number())
		s.Equal("", wisp.EmptyPhone.AreaCode())
	})

	s.Run("State Checks", func() {
		s.True(mobile.IsMobile())
		s.False(mobile.IsLandline())
		s.False(mobile.IsZero())

		s.False(landline.IsMobile())
		s.True(landline.IsLandline())
		s.False(landline.IsZero())

		s.True(wisp.EmptyPhone.IsZero())
		s.False(wisp.EmptyPhone.IsMobile())
		s.False(wisp.EmptyPhone.IsLandline())
	})
}

func (s *PhoneSuite) TestPhone_Formatted() {
	mobile, _ := wisp.NewPhone("5562982870053")
	landline, _ := wisp.NewPhone("551145671234")

	s.Equal("+55 (62) 98287-0053", mobile.Formatted())
	s.Equal("+55 (11) 4567-1234", landline.Formatted())
	s.Equal("", wisp.EmptyPhone.Formatted())
}

func (s *PhoneSuite) TestPhone_JSONMarshaling() {
	s.Run("should marshal and unmarshal correctly", func() {
		phone, _ := wisp.NewPhone("+55 (62) 98287-0053")
		data, err := json.Marshal(phone)
		s.Require().NoError(err)
		s.Equal(`"5562982870053"`, string(data))

		var unmarshaledPhone wisp.Phone
		err = json.Unmarshal(data, &unmarshaledPhone)
		s.Require().NoError(err)
		s.Equal(phone, unmarshaledPhone)
	})

	s.Run("should fail to unmarshal an invalid phone number", func() {
		var phone wisp.Phone
		err := json.Unmarshal([]byte(`"123"`), &phone)
		s.Require().Error(err)
	})
}

func (s *PhoneSuite) TestPhone_DatabaseInterface() {
	s.Run("Value", func() {
		phone, _ := wisp.NewPhone("5562982870053")
		val, err := phone.Value()
		s.Require().NoError(err)
		s.Equal("5562982870053", val)

		nilVal, err := wisp.EmptyPhone.Value()
		s.Require().NoError(err)
		s.Nil(nilVal)
	})

	s.Run("Scan", func() {
		s.Run("should scan a valid string", func() {
			var phone wisp.Phone
			err := phone.Scan("551145671234")
			s.Require().NoError(err)
			s.Equal(wisp.Phone("551145671234"), phone)
		})

		s.Run("should scan nil as EmptyPhone", func() {
			var phone wisp.Phone
			err := phone.Scan(nil)
			s.Require().NoError(err)
			s.True(phone.IsZero())
		})

		s.Run("should fail to scan an invalid string", func() {
			var phone wisp.Phone
			err := phone.Scan("invalid-number")
			s.Require().Error(err)
		})

		s.Run("should fail to scan an incompatible type", func() {
			var phone wisp.Phone
			err := phone.Scan(12345)
			s.Require().Error(err)
		})
	})
}
