package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type PortNumberSuite struct {
	suite.Suite
}

func TestPortNumberSuite(t *testing.T) {
	suite.Run(t, new(PortNumberSuite))
}

func (s *PortNumberSuite) TestNewPortNumber() {
	testCases := []struct {
		name        string
		input       int
		expectError bool
	}{
		{name: "should create a valid well-known port", input: 80},
		{name: "should create a valid ephemeral port", input: 49152},
		{name: "should create a valid port at lower bound", input: 1},
		{name: "should create a valid port at upper bound", input: 65535},
		{name: "should fail for port zero", input: 0, expectError: true},
		{name: "should fail for negative port", input: -1, expectError: true},
		{name: "should fail for port above upper bound", input: 65536, expectError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			port, err := wisp.NewPortNumber(tc.input)
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Equal(uint16(tc.input), port.Uint16())
			}
		})
	}
}

func (s *PortNumberSuite) TestPortNumber_JSON_SQL() {
	port, _ := wisp.NewPortNumber(8080)

	s.Run("JSON Marshaling and Unmarshaling", func() {
		data, err := json.Marshal(port)
		s.Require().NoError(err)
		s.JSONEq(`8080`, string(data))

		var unmarshaledPort wisp.PortNumber
		err = json.Unmarshal(data, &unmarshaledPort)
		s.Require().NoError(err)
		s.Equal(port, unmarshaledPort)

		err = json.Unmarshal([]byte(`0`), &unmarshaledPort)
		s.Require().Error(err)
	})

	s.Run("SQL Interface", func() {
		val, err := port.Value()
		s.Require().NoError(err)
		s.Equal(int64(8080), val)

		var scannedPort wisp.PortNumber
		err = scannedPort.Scan(int64(3000))
		s.Require().NoError(err)
		s.Equal(uint16(3000), scannedPort.Uint16())

		err = scannedPort.Scan(int64(0))
		s.Require().Error(err)
	})
}
