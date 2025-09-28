package wisp_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type TimezoneSuite struct {
	suite.Suite
}

func TestTimezoneSuite(t *testing.T) {
	suite.Run(t, new(TimezoneSuite))
}

func (s *TimezoneSuite) SetupTest() {
	wisp.ClearRegisteredTimezones()
}

func (s *TimezoneSuite) TestRegisterTimezones() {
	s.Run("should succeed in registering valid IANA names", func() {
		wisp.ClearRegisteredTimezones()
		err := wisp.RegisterTimezones("America/Sao_Paulo", "UTC")
		s.Require().NoError(err)
		s.True(wisp.IsTimezoneRegistered("UTC"))
		s.True(wisp.IsTimezoneRegistered("America/Sao_Paulo"))
	})

	s.Run("should fail and not register any zones if one is invalid", func() {
		wisp.ClearRegisteredTimezones()
		err := wisp.RegisterTimezones("UTC", "Invalid/Zone")
		s.Require().Error(err)

		s.Contains(err.Error(), "Invalid/Zone")

		s.False(wisp.IsTimezoneRegistered("UTC"))
	})
}

func (s *TimezoneSuite) TestNewTimezone() {
	err := wisp.RegisterTimezones("America/Sao_Paulo", "UTC")
	s.Require().NoError(err)

	s.Run("should create a valid timezone that is registered", func() {
		tz, err := wisp.NewTimezone("America/Sao_Paulo")
		s.Require().NoError(err)
		s.Equal("America/Sao_Paulo", tz.String())
	})

	s.Run("should fail for a valid IANA zone that is not registered", func() {
		_, err := wisp.NewTimezone("Europe/London")
		s.Require().Error(err)
		s.Contains(err.Error(), "not registered in the allowed list")
	})

	s.Run("should fail for an empty string", func() {
		_, err := wisp.NewTimezone("")
		s.Require().Error(err)
	})
}

func (s *TimezoneSuite) TestTimezone_Methods() {
	err := wisp.RegisterTimezones("America/Sao_Paulo", "UTC", "Europe/London")
	s.Require().NoError(err)

	utcTZ, _ := wisp.NewTimezone("UTC")
	saoPauloTZ, _ := wisp.NewTimezone("America/Sao_Paulo")
	saoPauloTZClone, _ := wisp.NewTimezone("America/Sao_Paulo")

	s.Run("Convert", func() {
		utcTime := time.Date(2025, time.September, 28, 12, 0, 0, 0, time.UTC)
		convertedTime := saoPauloTZ.Convert(utcTime)
		s.Equal("America/Sao_Paulo", convertedTime.Location().String())
		s.Equal(9, convertedTime.Hour())
	})

	s.Run("Equals", func() {
		s.True(saoPauloTZ.Equals(saoPauloTZClone))
		s.False(saoPauloTZ.Equals(utcTZ))
	})
}

func (s *TimezoneSuite) TestTimezone_JSON_SQL() {
	err := wisp.RegisterTimezones("Europe/London", "America/New_York")
	s.Require().NoError(err)
	tz, _ := wisp.NewTimezone("Europe/London")

	s.Run("JSON Marshaling and Unmarshaling", func() {
		data, err := json.Marshal(tz)
		s.Require().NoError(err)
		s.JSONEq(`"Europe/London"`, string(data))

		var unmarshaledTz wisp.Timezone
		err = json.Unmarshal(data, &unmarshaledTz)
		s.Require().NoError(err)
		s.True(tz.Equals(unmarshaledTz))

		unregisteredJSON := `"Asia/Tokyo"`
		err = json.Unmarshal([]byte(unregisteredJSON), &unmarshaledTz)
		s.Require().Error(err)
	})

	s.Run("SQL Interface", func() {
		val, err := tz.Value()
		s.Require().NoError(err)
		s.Equal("Europe/London", val)

		var scannedTz wisp.Timezone
		err = scannedTz.Scan("America/New_York")
		s.Require().NoError(err)
		s.Equal("America/New_York", scannedTz.String())
	})
}
