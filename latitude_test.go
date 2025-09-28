package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type LatitudeSuite struct {
	suite.Suite
}

func TestLatitudeSuite(t *testing.T) {
	suite.Run(t, new(LatitudeSuite))
}

func (s *LatitudeSuite) TestNewLatitude() {
	testCases := []struct {
		name        string
		input       float64
		expectError bool
	}{
		{name: "should create a valid positive latitude", input: 23.5505},
		{name: "should create a valid negative latitude", input: -46.6333},
		{name: "should create a valid latitude at lower bound", input: -90.0},
		{name: "should create a valid latitude at upper bound", input: 90.0},
		{name: "should fail for latitude below lower bound", input: -90.1, expectError: true},
		{name: "should fail for latitude above upper bound", input: 90.1, expectError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			lat, err := wisp.NewLatitude(tc.input)
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.input, lat.Float64())
			}
		})
	}
}

func (s *LatitudeSuite) TestLatitude_JSON_SQL() {
	lat, _ := wisp.NewLatitude(-23.5505)

	s.Run("JSON Marshaling and Unmarshaling", func() {
		data, err := json.Marshal(lat)
		s.Require().NoError(err)
		s.JSONEq(`-23.5505`, string(data))

		var unmarshaledLat wisp.Latitude
		err = json.Unmarshal(data, &unmarshaledLat)
		s.Require().NoError(err)
		s.Equal(lat, unmarshaledLat)

		err = json.Unmarshal([]byte(`-100.0`), &unmarshaledLat)
		s.Require().Error(err)
	})

	s.Run("SQL Interface", func() {
		val, err := lat.Value()
		s.Require().NoError(err)
		s.Equal(-23.5505, val)

		var scannedLat wisp.Latitude
		err = scannedLat.Scan(float64(-23.5505))
		s.Require().NoError(err)
		s.Equal(lat, scannedLat)
	})
}
