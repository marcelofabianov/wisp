package atomic_test

import (
	"strings"
	"testing"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/atomic"
)

type EmailSuite struct {
	suite.Suite
}

func TestEmailSuite(t *testing.T) {
	suite.Run(t, new(EmailSuite))
}

func (s *EmailSuite) TestNewEmail() {
	testCases := []struct {
		name          string
		input         string
		expectedEmail atomic.Email
		expectError   bool
		expectedCode  fault.Code
	}{
		{
			name:          "should create a valid email",
			input:         "test@example.com",
			expectedEmail: "test@example.com",
			expectError:   false,
		},
		{
			name:          "should normalize email by trimming spaces and lowercasing",
			input:         "  Test.User@Example.COM  ",
			expectedEmail: "test.user@example.com",
			expectError:   false,
		},
		{
			name:          "should handle subdomains correctly",
			input:         "user@sub.domain.co.uk",
			expectedEmail: "user@sub.domain.co.uk",
			expectError:   false,
		},
		{
			name:          "should parse email from a display name format",
			input:         `"John Doe" <john.doe@work.org>`,
			expectedEmail: "john.doe@work.org",
			expectError:   false,
		},
		{
			name:         "should fail for empty string",
			input:        "",
			expectError:  true,
			expectedCode: fault.Invalid,
		},
		{
			name:         "should fail for blank string",
			input:        "   ",
			expectError:  true,
			expectedCode: fault.Invalid,
		},
		{
			name:         "should fail for invalid format",
			input:        "just-a-string",
			expectError:  true,
			expectedCode: fault.Invalid,
		},
		{
			name:         "should fail for missing local part",
			input:        "@example.com",
			expectError:  true,
			expectedCode: fault.Invalid,
		},
		{
			name:         "should fail for email exceeding max length",
			input:        strings.Repeat("a", 245) + "@example.com", // 245 + 1 + 11 = 257
			expectError:  true,
			expectedCode: fault.Invalid,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			email, err := atomic.NewEmail(tc.input)

			if tc.expectError {
				s.Require().Error(err)
				s.Equal(atomic.EmptyEmail, email)
				faultErr, ok := err.(*fault.Error)
				s.Require().True(ok, "error should be of type *fault.Error")
				s.Equal(tc.expectedCode, faultErr.Code)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expectedEmail, email)
			}
		})
	}
}

func (s *EmailSuite) TestMustNewEmail() {
	s.Run("should not panic for a valid email", func() {
		var email atomic.Email
		s.NotPanics(func() {
			email = atomic.MustNewEmail("test@example.com")
		})
		s.Equal(atomic.Email("test@example.com"), email)
	})

	s.Run("should panic for an invalid email", func() {
		s.Panics(func() {
			atomic.MustNewEmail("not-a-valid-email")
		})
	})
}

func (s *EmailSuite) TestEmail_IsEmptyAndString() {
	s.Run("should correctly report empty status", func() {
		email := atomic.MustNewEmail("test@example.com")
		s.False(email.IsEmpty())
		s.Equal("test@example.com", email.String())

		s.True(atomic.EmptyEmail.IsEmpty())
		s.Equal("", atomic.EmptyEmail.String())
	})
}

func (s *EmailSuite) TestEmail_JSONMarshaling() {
	s.Run("should correctly marshal and unmarshal valid email", func() {
		email := atomic.MustNewEmail("user@domain.com")

		jsonData, err := email.MarshalJSON()
		s.Require().NoError(err)
		s.Equal(`"user@domain.com"`, string(jsonData))

		var unmarshaledEmail atomic.Email
		err = unmarshaledEmail.UnmarshalJSON(jsonData)
		s.Require().NoError(err)
		s.Equal(email, unmarshaledEmail)
	})

	s.Run("should return error when unmarshaling invalid JSON", func() {
		var email atomic.Email
		err := email.UnmarshalJSON([]byte(`not-a-string`))
		s.Require().Error(err)
		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok)
		s.Equal(fault.Invalid, faultErr.Code)
	})

	s.Run("should return error when unmarshaling invalid email format from JSON", func() {
		var email atomic.Email
		err := email.UnmarshalJSON([]byte(`"invalid-email-format"`))
		s.Require().Error(err)
		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok)
		s.Equal(fault.Invalid, faultErr.Code)
	})
}

func (s *EmailSuite) TestEmail_TextMarshaling() {
	s.Run("should correctly marshal and unmarshal valid email", func() {
		email := atomic.MustNewEmail("user@domain.com")

		textData, err := email.MarshalText()
		s.Require().NoError(err)
		s.Equal("user@domain.com", string(textData))

		var unmarshaledEmail atomic.Email
		err = unmarshaledEmail.UnmarshalText(textData)
		s.Require().NoError(err)
		s.Equal(email, unmarshaledEmail)
	})

	s.Run("should return error when unmarshaling invalid email format from text", func() {
		var email atomic.Email
		err := email.UnmarshalText([]byte("invalid-email-format"))
		s.Require().Error(err)
	})
}

func (s *EmailSuite) TestEmail_Value() {
	s.Run("should return string for a valid email", func() {
		email := atomic.MustNewEmail("test@example.com")
		val, err := email.Value()
		s.Require().NoError(err)
		s.Equal("test@example.com", val)
	})

	s.Run("should return nil for an empty email", func() {
		val, err := atomic.EmptyEmail.Value()
		s.Require().NoError(err)
		s.Nil(val)
	})
}

func (s *EmailSuite) TestEmail_Scan() {
	testCases := []struct {
		name        string
		src         interface{}
		expected    atomic.Email
		expectError bool
	}{
		{
			name:     "should scan a valid string",
			src:      "scan.test@example.com",
			expected: "scan.test@example.com",
		},
		{
			name:     "should scan a valid byte slice",
			src:      []byte("scan.bytes@example.com"),
			expected: "scan.bytes@example.com",
		},
		{
			name:     "should scan nil into an empty email",
			src:      nil,
			expected: atomic.EmptyEmail,
		},
		{
			name:        "should fail to scan an invalid string",
			src:         "invalid-email",
			expectError: true,
		},
		{
			name:        "should fail to scan an incompatible type",
			src:         12345,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var email atomic.Email
			err := email.Scan(tc.src)

			if tc.expectError {
				s.Require().Error(err)
				faultErr, ok := err.(*fault.Error)
				s.Require().True(ok)
				s.Equal(fault.Invalid, faultErr.Code)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expected, email)
			}
		})
	}
}
