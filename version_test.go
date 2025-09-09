package atomic_test

import (
	"encoding/json"
	"testing"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/atomic"
)

type VersionSuite struct {
	suite.Suite
}

func TestVersionSuite(t *testing.T) {
	suite.Run(t, new(VersionSuite))
}

func (s *VersionSuite) TestVersion_InitialAndZeroStates() {
	s.Run("should return 1 for initial version", func() {
		initial := atomic.InitialVersion()
		s.Equal(atomic.Version(1), initial)
		s.Equal(1, initial.Int())
		s.False(initial.IsZero())
	})

	s.Run("should be 0 for zero version", func() {
		s.Equal(atomic.Version(0), atomic.ZeroVersion)
		s.Equal(0, atomic.ZeroVersion.Int())
		s.True(atomic.ZeroVersion.IsZero())
	})
}

func (s *VersionSuite) TestVersion_Behaviors() {
	s.Run("should increment the version", func() {
		v := atomic.Version(5)
		v.Increment()
		s.Equal(atomic.Version(6), v)
	})

	s.Run("should return the previous version", func() {
		v5 := atomic.Version(5)
		s.Equal(atomic.Version(4), v5.Previous())

		v1 := atomic.Version(1)
		s.Equal(atomic.ZeroVersion, v1.Previous())

		vZ := atomic.ZeroVersion
		s.Equal(atomic.ZeroVersion, vZ.Previous())
	})
}

func (s *VersionSuite) TestVersion_JSONMarshaling() {
	s.Run("should marshal and unmarshal a valid version", func() {
		v := atomic.Version(42)
		data, err := json.Marshal(v)
		s.Require().NoError(err)
		s.Equal(`42`, string(data))

		var unmarshaledV atomic.Version
		err = json.Unmarshal(data, &unmarshaledV)
		s.Require().NoError(err)
		s.Equal(v, unmarshaledV)
	})

	s.Run("should unmarshal null as ZeroVersion", func() {
		var v atomic.Version
		err := json.Unmarshal([]byte("null"), &v)
		s.Require().NoError(err)
		s.True(v.IsZero())
	})

	s.Run("should fail to unmarshal non-numeric JSON", func() {
		var v atomic.Version
		err := json.Unmarshal([]byte(`"abc"`), &v)
		s.Require().Error(err)
		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok)
		s.Equal(fault.Invalid, faultErr.Code)
	})

	s.Run("should fail to unmarshal a negative number", func() {
		var v atomic.Version
		err := json.Unmarshal([]byte(`-10`), &v)
		s.Require().Error(err)
		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok)
		s.Equal(fault.Invalid, faultErr.Code)
	})
}

func (s *VersionSuite) TestVersion_Value() {
	s.Run("should return int64 for a valid version", func() {
		v := atomic.Version(123)
		val, err := v.Value()
		s.Require().NoError(err)
		s.Equal(int64(123), val)
	})

	s.Run("should return int64(0) for a zero version", func() {
		val, err := atomic.ZeroVersion.Value()
		s.Require().NoError(err)
		s.Equal(int64(0), val)
	})
}

func (s *VersionSuite) TestVersion_Scan() {
	testCases := []struct {
		name        string
		src         interface{}
		expected    atomic.Version
		expectError bool
	}{
		{
			name:     "should scan a valid int64",
			src:      int64(99),
			expected: atomic.Version(99),
		},
		{
			name:     "should scan a valid byte slice",
			src:      []byte("42"),
			expected: atomic.Version(42),
		},
		{
			name:     "should scan nil into ZeroVersion",
			src:      nil,
			expected: atomic.ZeroVersion,
		},
		{
			name:        "should fail to scan a negative int64",
			src:         int64(-1),
			expectError: true,
		},
		{
			name:        "should fail to scan negative bytes",
			src:         []byte("-5"),
			expectError: true,
		},
		{
			name:        "should fail to scan invalid bytes",
			src:         []byte("not-a-number"),
			expectError: true,
		},
		{
			name:        "should fail to scan an incompatible type",
			src:         float64(123.45),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var v atomic.Version
			err := v.Scan(tc.src)

			if tc.expectError {
				s.Require().Error(err)
				faultErr, ok := err.(*fault.Error)
				s.Require().True(ok)
				s.Equal(fault.Invalid, faultErr.Code)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expected, v)
			}
		})
	}
}
