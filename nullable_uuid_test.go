package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type NullableUUIDSuite struct {
	suite.Suite
}

func TestNullableUUIDSuite(t *testing.T) {
	suite.Run(t, new(NullableUUIDSuite))
}

func (s *NullableUUIDSuite) TestNewNullableUUID() {
	s.Run("should create a valid NullableUUID from a non-nil UUID", func() {
		id, _ := wisp.NewUUID()
		nu := wisp.NewNullableUUID(id)
		s.True(nu.Valid)
		s.Equal(id, nu.UUID)
		s.False(nu.IsZero())
	})

	s.Run("should create an invalid NullableUUID from a nil UUID", func() {
		nu := wisp.NewNullableUUID(wisp.Nil)
		s.False(nu.Valid)
		s.True(nu.IsZero())
	})
}

func (s *NullableUUIDSuite) TestNullableUUID_JSON() {
	id, _ := wisp.NewUUID()
	nu := wisp.NewNullableUUID(id)

	s.Run("should marshal a valid UUID to a JSON string", func() {
		data, err := json.Marshal(nu)
		expectedJSON, _ := json.Marshal(id)
		s.Require().NoError(err)
		s.JSONEq(string(expectedJSON), string(data))
	})

	s.Run("should marshal an invalid UUID to JSON null", func() {
		nilUUID := wisp.NewNullableUUID(wisp.Nil)
		data, err := json.Marshal(nilUUID)
		s.Require().NoError(err)
		s.Equal("null", string(data))
	})

	s.Run("should unmarshal a UUID string correctly", func() {
		jsonUUID, _ := json.Marshal(id)
		var unmarshaledNU wisp.NullableUUID
		err := json.Unmarshal(jsonUUID, &unmarshaledNU)
		s.Require().NoError(err)
		s.Equal(nu, unmarshaledNU)
	})

	s.Run("should unmarshal null correctly", func() {
		var unmarshaledNU wisp.NullableUUID
		err := json.Unmarshal([]byte("null"), &unmarshaledNU)
		s.Require().NoError(err)
		s.False(unmarshaledNU.Valid)
	})
}

func (s *NullableUUIDSuite) TestNullableUUID_SQL() {
	id, _ := wisp.NewUUID()
	nuValid := wisp.NewNullableUUID(id)
	nuInvalid := wisp.NewNullableUUID(wisp.Nil)

	s.Run("Value", func() {
		val, err := nuValid.Value()
		s.Require().NoError(err)
		s.Equal(id.String(), val)

		nilVal, err := nuInvalid.Value()
		s.Require().NoError(err)
		s.Nil(nilVal)
	})

	s.Run("Scan", func() {
		var scannedNU wisp.NullableUUID

		err := scannedNU.Scan(id.String())
		s.Require().NoError(err)
		s.True(scannedNU.Valid)
		s.Equal(id, scannedNU.UUID)

		err = scannedNU.Scan(nil)
		s.Require().NoError(err)
		s.False(scannedNU.Valid)
	})
}
