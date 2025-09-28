package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type IPAddressSuite struct {
	suite.Suite
}

func TestIPAddressSuite(t *testing.T) {
	suite.Run(t, new(IPAddressSuite))
}

func (s *IPAddressSuite) TestNewIPAddress() {
	testCases := []struct {
		name           string
		input          string
		expectedString string
		isV4           bool
		isV6           bool
		expectError    bool
	}{
		{name: "should create a valid IPv4 address", input: "192.168.1.1", expectedString: "192.168.1.1", isV4: true},
		{name: "should create and canonicalize a valid IPv6 address", input: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", expectedString: "2001:db8:85a3::8a2e:370:7334", isV6: true},
		{name: "should create a valid compressed IPv6 address", input: "::1", expectedString: "::1", isV6: true},
		{name: "should fail for an empty string", input: "", expectError: true},
		{name: "should fail for an invalid address", input: "999.999.999.999", expectError: true},
		{name: "should fail for a random string", input: "not-an-ip", expectError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ip, err := wisp.NewIPAddress(tc.input)
			if tc.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expectedString, ip.String()) // Corrigido para usar expectedString
				s.Equal(tc.isV4, ip.IsV4())
				s.Equal(tc.isV6, ip.IsV6())
			}
		})
	}
}

func (s *IPAddressSuite) TestIPAddress_JSON_SQL() {
	ip, _ := wisp.NewIPAddress("127.0.0.1")

	s.Run("JSON Marshaling and Unmarshaling", func() {
		data, err := json.Marshal(ip)
		s.Require().NoError(err)
		s.JSONEq(`"127.0.0.1"`, string(data))

		var unmarshaledIP wisp.IPAddress
		err = json.Unmarshal(data, &unmarshaledIP)
		s.Require().NoError(err)
		s.Equal(ip, unmarshaledIP)

		err = json.Unmarshal([]byte(`"invalid"`), &unmarshaledIP)
		s.Require().Error(err)
	})

	s.Run("SQL Interface", func() {
		val, err := ip.Value()
		s.Require().NoError(err)
		s.Equal("127.0.0.1", val)

		var scannedIP wisp.IPAddress
		err = scannedIP.Scan("8.8.8.8")
		s.Require().NoError(err)
		s.Equal("8.8.8.8", scannedIP.String())

		err = scannedIP.Scan(nil)
		s.Require().NoError(err)
		s.True(scannedIP.IsZero())
	})
}
