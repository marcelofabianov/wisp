package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type LongitudeSuite struct {
	suite.Suite
}

func TestLongitudeSuite(t *testing.T) {
	suite.Run(t, new(LongitudeSuite))
}

func (s *LongitudeSuite) TestNewLongitude() {
	testCases := []struct {
		name        string
		input       float64
		expectError bool
	}{
		{name: "should create a valid positive longitude", input: 45.123},
		{name: "should create a valid negative longitude", input: -120.456},
		{name: "should create a valid longitude at lower bound", input: -180.0},
		{name: "should create a valid longitude at upper bound", input: 180.0},
		{name: "should fail for longitude below lower bound", input: -180.1, expectError: true},
		{name: "should fail for longitude above upper bound", input: 180.1, expectError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			lon, err := wisp.NewLongitude(tc.input)
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.input, lon.Float64())
			}
		})
	}
}

func (s *LongitudeSuite) TestLongitude_JSON_SQL() {
	lon, _ := wisp.NewLongitude(-49.2738)

	s.Run("JSON Marshaling and Unmarshaling", func() {
		data, err := json.Marshal(lon)
		s.Require().NoError(err)
		s.JSONEq(`-49.2738`, string(data))

		var unmarshaledLon wisp.Longitude
		err = json.Unmarshal(data, &unmarshaledLon)
		s.Require().NoError(err)
		s.Equal(lon, unmarshaledLon)

		err = json.Unmarshal([]byte(`200.0`), &unmarshaledLon)
		s.Require().Error(err)
	})

	s.Run("SQL Interface", func() {
		val, err := lon.Value()
		s.Require().NoError(err)
		s.Equal(-49.2738, val)

		var scannedLon wisp.Longitude
		err = scannedLon.Scan(float64(-49.2738))
		s.Require().NoError(err)
		s.Equal(lon, scannedLon)

		err = scannedLon.Scan("not a float")
		s.Require().Error(err)
	})
}
