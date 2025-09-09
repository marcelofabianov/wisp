package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/marcelofabianov/fault"
	wisp "github.com/marcelofabianov/wisp"
	"github.com/stretchr/testify/suite"
)

type AuditUserSuite struct {
	suite.Suite
}

func TestAuditUserSuite(t *testing.T) {
	suite.Run(t, new(AuditUserSuite))
}

func (s *AuditUserSuite) TestNewAuditUser() {
	testCases := []struct {
		name        string
		input       string
		expected    wisp.AuditUser
		expectError bool
	}{
		{name: "should create a valid user from an email", input: "test@example.com", expected: "test@example.com"},
		{name: "should create a valid user from 'system' literal", input: "system", expected: wisp.SystemAuditUser},
		{name: "should create and normalize from 'SYSTEM' literal", input: " SYSTEM ", expected: wisp.SystemAuditUser},
		{name: "should create an empty user from an empty string", input: "", expected: wisp.EmptyAuditUser},
		{name: "should fail for an invalid string that is not an email", input: "some_user", expectError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			user, err := wisp.NewAuditUser(tc.input)
			if tc.expectError {
				s.Require().Error(err)
				s.Equal(wisp.EmptyAuditUser, user)
				faultErr, ok := err.(*fault.Error)
				s.Require().True(ok)
				s.Equal(fault.Invalid, faultErr.Code)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expected, user)
			}
		})
	}
}

func (s *AuditUserSuite) TestAuditUser_Methods() {
	emailUser, _ := wisp.NewAuditUser("test@example.com")
	systemUser, _ := wisp.NewAuditUser("system")
	emptyUser := wisp.EmptyAuditUser

	s.Run("State Checks", func() {
		s.True(emailUser.IsEmail())
		s.False(emailUser.IsSystem())
		s.False(emailUser.IsZero())

		s.False(systemUser.IsEmail())
		s.True(systemUser.IsSystem())
		s.False(systemUser.IsZero())

		s.False(emptyUser.IsEmail())
		s.False(emptyUser.IsSystem())
		s.True(emptyUser.IsZero())
	})

	s.Run("Email", func() {
		email, ok := emailUser.Email()
		s.True(ok)
		s.Equal(wisp.Email("test@example.com"), email)

		_, ok = systemUser.Email()
		s.False(ok)

		_, ok = emptyUser.Email()
		s.False(ok)
	})
}

func (s *AuditUserSuite) TestAuditUser_JSONMarshaling() {
	s.Run("should marshal and unmarshal a valid user email", func() {
		user, _ := wisp.NewAuditUser("test@example.com")
		data, err := json.Marshal(user)
		s.Require().NoError(err)
		s.Equal(`"test@example.com"`, string(data))

		var unmarshaledUser wisp.AuditUser
		err = json.Unmarshal(data, &unmarshaledUser)
		s.Require().NoError(err)
		s.Equal(user, unmarshaledUser)
	})

	s.Run("should fail to unmarshal an invalid user string", func() {
		var user wisp.AuditUser
		err := json.Unmarshal([]byte(`"invalid-user"`), &user)
		s.Require().Error(err)
	})
}

func (s *AuditUserSuite) TestAuditUser_DatabaseInterface() {
	s.Run("Value", func() {
		systemUser, _ := wisp.NewAuditUser("system")
		val, err := systemUser.Value()
		s.Require().NoError(err)
		s.Equal("system", val)

		nilVal, err := wisp.EmptyAuditUser.Value()
		s.Require().NoError(err)
		s.Nil(nilVal)
	})

	s.Run("Scan", func() {
		var user wisp.AuditUser
		err := user.Scan("test@example.com")
		s.Require().NoError(err)
		s.Equal(wisp.AuditUser("test@example.com"), user)

		err = user.Scan(nil)
		s.Require().NoError(err)
		s.True(user.IsZero())

		err = user.Scan("invalid-user")
		s.Require().Error(err)
	})
}
