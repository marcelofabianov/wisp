package atomic_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/atomic"
)

const validUUIDString = "7a4a8862-8354-4b53-9b64-42a984a37218"
const nilUUIDString = "00000000-0000-0000-0000-000000000000"

type UUIDSuite struct {
	suite.Suite
}

func TestUUIDSuite(t *testing.T) {
	suite.Run(t, new(UUIDSuite))
}

func (s *UUIDSuite) TestNewUUID_Success() {
	s.Run("should create a new, non-nil UUID successfully", func() {
		id, err := atomic.NewUUID()

		s.Require().NoError(err)
		s.Require().NotEqual(atomic.Nil, id)
		s.Require().False(id.IsNil())

		_, parseErr := uuid.Parse(id.String())
		s.NoError(parseErr, "the generated UUID string should be valid")
	})
}

func (s *UUIDSuite) TestParseUUID() {
	s.Run("should parse a valid UUID string successfully", func() {
		id, err := atomic.ParseUUID(validUUIDString)

		s.Require().NoError(err)
		s.Equal(validUUIDString, id.String())
	})

	s.Run("should return an error for an invalid UUID string", func() {
		id, err := atomic.ParseUUID("not-a-valid-uuid")

		s.Require().Error(err)
		s.Equal(atomic.Nil, id)

		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok, "error should be of type *fault.Error")
		s.Equal(fault.Invalid, faultErr.Code)
	})
}

func (s *UUIDSuite) TestMustParseUUID() {
	s.Run("should not panic for a valid UUID string", func() {
		var id atomic.UUID
		s.NotPanics(func() {
			id = atomic.MustParseUUID(validUUIDString)
		})
		s.Equal(validUUIDString, id.String())
	})

	s.Run("should panic for an invalid UUID string", func() {
		s.Panics(func() {
			atomic.MustParseUUID("not-a-valid-uuid")
		})
	})
}

func (s *UUIDSuite) TestUUID_IsNil() {
	s.Run("should return true for a nil UUID", func() {
		s.True(atomic.Nil.IsNil())
	})

	s.Run("should return true for a zero-value UUID", func() {
		var zeroID atomic.UUID
		s.True(zeroID.IsNil())
	})

	s.Run("should return false for a non-nil UUID", func() {
		id := atomic.MustParseUUID(validUUIDString)
		s.False(id.IsNil())
	})
}

func (s *UUIDSuite) TestUUID_String() {
	s.Run("should return the correct string for a non-nil UUID", func() {
		id := atomic.MustParseUUID(validUUIDString)
		s.Equal(validUUIDString, id.String())
	})

	s.Run("should return the nil UUID string for a nil UUID", func() {
		s.Equal(nilUUIDString, atomic.Nil.String())
	})
}

func (s *UUIDSuite) TestUUID_MarshalText() {
	s.Run("should marshal a valid UUID to text", func() {
		id := atomic.MustParseUUID(validUUIDString)
		text, err := id.MarshalText()

		s.Require().NoError(err)
		s.Equal([]byte(validUUIDString), text)
	})

	s.Run("should marshal a nil UUID to text", func() {
		text, err := atomic.Nil.MarshalText()

		s.Require().NoError(err)
		s.Equal([]byte(nilUUIDString), text)
	})
}

func (s *UUIDSuite) TestUUID_UnmarshalText() {
	s.Run("should unmarshal a valid text into a UUID", func() {
		var id atomic.UUID
		err := id.UnmarshalText([]byte(validUUIDString))

		s.Require().NoError(err)
		s.Equal(validUUIDString, id.String())
	})

	s.Run("should return an error for invalid text", func() {
		var id atomic.UUID
		err := id.UnmarshalText([]byte("invalid-text"))

		s.Require().Error(err)
		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok, "error should be of type *fault.Error")
		s.Equal(fault.Invalid, faultErr.Code)
	})
}

func (s *UUIDSuite) TestUUID_Value() {
	s.Run("should return a string value for a valid UUID", func() {
		id := atomic.MustParseUUID(validUUIDString)
		val, err := id.Value()

		s.Require().NoError(err)
		s.Equal(validUUIDString, val)
	})

	s.Run("should return nil for a nil UUID", func() {
		val, err := atomic.Nil.Value()

		s.Require().NoError(err)
		s.Nil(val)
	})
}

func (s *UUIDSuite) TestUUID_Scan() {
	testCases := []struct {
		name    string
		src     interface{}
		wantErr bool
		wantVal string
	}{
		{
			name:    "should scan a valid string",
			src:     validUUIDString,
			wantErr: false,
			wantVal: validUUIDString,
		},
		{
			name:    "should scan a valid byte slice",
			src:     []byte(validUUIDString),
			wantErr: false,
			wantVal: validUUIDString,
		},
		{
			name:    "should scan a nil value",
			src:     nil,
			wantErr: false,
			wantVal: nilUUIDString,
		},
		{
			name:    "should return error for an invalid string",
			src:     "invalid-uuid-string",
			wantErr: true,
			wantVal: "",
		},
		{
			name:    "should return error for an incompatible type",
			src:     int64(12345),
			wantErr: true,
			wantVal: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var id atomic.UUID
			err := id.Scan(tc.src)

			if tc.wantErr {
				s.Require().Error(err)
				faultErr, ok := err.(*fault.Error)
				s.Require().True(ok, "error should be of type *fault.Error")
				s.Equal(fault.Invalid, faultErr.Code)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.wantVal, id.String())
			}
		})
	}
}
