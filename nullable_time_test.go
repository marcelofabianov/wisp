package atomic_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/atomic"
)

type NullableTimeSuite struct {
	suite.Suite
}

func TestNullableTimeSuite(t *testing.T) {
	suite.Run(t, new(NullableTimeSuite))
}

func (s *NullableTimeSuite) TestNewNullableTime() {
	s.Run("should create a valid NullableTime from a non-zero time", func() {
		now := time.Now()
		nt := atomic.NewNullableTime(now)
		s.True(nt.Valid)
		s.Equal(now, nt.Time)
		s.False(nt.IsZero())
	})

	s.Run("should create an invalid NullableTime from a zero time", func() {
		nt := atomic.NewNullableTime(time.Time{})
		s.False(nt.Valid)
		s.True(nt.IsZero())
	})
}

func (s *NullableTimeSuite) TestNullableTime_JSONMarshaling() {
	s.Run("should marshal a valid time to a JSON string", func() {
		t := time.Date(2025, 9, 9, 12, 30, 0, 0, time.UTC)
		nt := atomic.NewNullableTime(t)
		data, err := json.Marshal(nt)

		expectedJSON, _ := json.Marshal(t)
		s.Require().NoError(err)
		s.JSONEq(string(expectedJSON), string(data))
	})

	s.Run("should marshal an invalid time to JSON null", func() {
		nt := atomic.EmptyNullableTime
		data, err := json.Marshal(nt)
		s.Require().NoError(err)
		s.Equal("null", string(data))
	})

	s.Run("should unmarshal a time string correctly", func() {
		t := time.Date(2025, 9, 9, 12, 30, 0, 0, time.UTC)
		expectedNT := atomic.NewNullableTime(t)
		jsonTime, _ := json.Marshal(t)

		var nt atomic.NullableTime
		err := json.Unmarshal(jsonTime, &nt)
		s.Require().NoError(err)
		s.Equal(expectedNT, nt)
	})

	s.Run("should unmarshal null correctly", func() {
		var nt atomic.NullableTime
		err := json.Unmarshal([]byte("null"), &nt)
		s.Require().NoError(err)
		s.False(nt.Valid)
		s.True(nt.IsZero())
	})
}

func (s *NullableTimeSuite) TestNullableTime_DatabaseInterface() {
	s.Run("Value", func() {
		now := time.Now()
		ntValid := atomic.NewNullableTime(now)
		val, err := ntValid.Value()
		s.Require().NoError(err)
		s.Equal(now, val)

		ntInvalid := atomic.EmptyNullableTime
		nilVal, err := ntInvalid.Value()
		s.Require().NoError(err)
		s.Nil(nilVal)
	})

	s.Run("Scan", func() {
		now := time.Now()
		var nt atomic.NullableTime

		err := nt.Scan(now)
		s.Require().NoError(err)
		s.True(nt.Valid)
		s.Equal(now, nt.Time)

		err = nt.Scan(nil)
		s.Require().NoError(err)
		s.False(nt.Valid)
	})
}
