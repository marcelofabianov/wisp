package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type TimeOfDaySuite struct {
	suite.Suite
}

func TestTimeOfDaySuite(t *testing.T) {
	suite.Run(t, new(TimeOfDaySuite))
}

func (s *TimeOfDaySuite) TestNewTimeOfDay() {
	s.Run("should create a valid time of day", func() {
		tod, err := wisp.NewTimeOfDay(14, 30)
		s.Require().NoError(err)
		s.Equal(14, tod.Hour())
		s.Equal(30, tod.Minute())
	})

	s.Run("should create valid times at boundaries", func() {
		tod, err := wisp.NewTimeOfDay(0, 0)
		s.Require().NoError(err)
		s.Equal(0, tod.Hour())

		tod, err = wisp.NewTimeOfDay(23, 59)
		s.Require().NoError(err)
		s.Equal(23, tod.Hour())
		s.Equal(59, tod.Minute())
	})

	s.Run("should fail for out of range values", func() {
		_, err := wisp.NewTimeOfDay(24, 0)
		s.Require().Error(err)

		_, err = wisp.NewTimeOfDay(-1, 0)
		s.Require().Error(err)

		_, err = wisp.NewTimeOfDay(10, 60)
		s.Require().Error(err)

		_, err = wisp.NewTimeOfDay(10, -1)
		s.Require().Error(err)
	})
}

func (s *TimeOfDaySuite) TestParseTimeOfDay() {
	s.Run("should parse a valid HH:MM string", func() {
		tod, err := wisp.ParseTimeOfDay("16:05")
		s.Require().NoError(err)
		s.Equal(16, tod.Hour())
		s.Equal(5, tod.Minute())
	})

	s.Run("should fail for invalid formats", func() {
		_, err := wisp.ParseTimeOfDay("16:05:30")
		s.Require().Error(err)

		_, err = wisp.ParseTimeOfDay("9:30")
		s.Require().Error(err)

		_, err = wisp.ParseTimeOfDay("abc")
		s.Require().Error(err)
	})
}

func (s *TimeOfDaySuite) TestTimeOfDay_Methods() {
	tod1, _ := wisp.NewTimeOfDay(9, 0)
	tod2, _ := wisp.NewTimeOfDay(15, 30)

	s.True(tod1.Before(tod2))
	s.False(tod2.Before(tod1))
	s.True(tod2.After(tod1))
	s.False(tod1.After(tod2))
	s.Equal("09:00", tod1.String())
}

func (s *TimeOfDaySuite) TestTimeOfDay_JSON() {
	tod, _ := wisp.NewTimeOfDay(22, 15)

	s.Run("should marshal and unmarshal correctly", func() {
		data, err := json.Marshal(tod)
		s.Require().NoError(err)
		s.Equal(`"22:15"`, string(data))

		var unmarshaledTOD wisp.TimeOfDay
		err = json.Unmarshal(data, &unmarshaledTOD)
		s.Require().NoError(err)
		s.Equal(tod, unmarshaledTOD)
	})

	s.Run("should fail for invalid JSON data", func() {
		var unmarshaledTOD wisp.TimeOfDay
		err := json.Unmarshal([]byte(`"invalid"`), &unmarshaledTOD)
		s.Require().Error(err)
	})
}

func (s *TimeOfDaySuite) TestTimeOfDay_SQL() {
	tod, _ := wisp.NewTimeOfDay(10, 25) // 10 * 60 + 25 = 625 minutes

	s.Run("Value", func() {
		val, err := tod.Value()
		s.Require().NoError(err)
		s.Equal(int64(625), val)
	})

	s.Run("Scan", func() {
		var scannedTOD wisp.TimeOfDay
		err := scannedTOD.Scan(int64(625))
		s.Require().NoError(err)
		s.Equal(tod, scannedTOD)

		err = scannedTOD.Scan(nil)
		s.Require().NoError(err)
		s.True(scannedTOD.IsZero())

		err = scannedTOD.Scan(int64(9999))
		s.Require().Error(err)
	})
}
