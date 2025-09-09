package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/suite"

	wisp "github.com/marcelofabianov/wisp"
)

type PercentageSuite struct {
	suite.Suite
}

func TestPercentageSuite(t *testing.T) {
	suite.Run(t, new(PercentageSuite))
}

func (s *PercentageSuite) TestNewPercentageFromFloat() {
	s.Run("should create a valid percentage", func() {
		p, err := wisp.NewPercentageFromFloat(0.155) // 15.5%
		s.Require().NoError(err)
		s.InDelta(0.155, p.Float64(), 0.0001)
		s.Equal("15.50%", p.String())
	})

	s.Run("should create a zero percentage", func() {
		p, err := wisp.NewPercentageFromFloat(0)
		s.Require().NoError(err)
		s.True(p.IsZero())
	})

	s.Run("should fail for a negative percentage", func() {
		_, err := wisp.NewPercentageFromFloat(-0.1)
		s.Require().Error(err)
		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok)
		s.Equal(fault.Invalid, faultErr.Code)
	})

	s.Run("should handle rounding correctly", func() {
		// 0.12345 -> 1234.5 -> rounds to 1234 (even)
		p, err := wisp.NewPercentageFromFloat(0.12345)
		s.Require().NoError(err)
		s.Equal(wisp.Percentage(1234), p)

		// 0.12355 -> 1235.5 -> rounds to 1236 (even)
		p, err = wisp.NewPercentageFromFloat(0.12355)
		s.Require().NoError(err)
		s.Equal(wisp.Percentage(1236), p)
	})
}

func (s *PercentageSuite) TestPercentage_ApplyTo() {
	money, _ := wisp.NewMoney(10000, wisp.BRL) // R$ 100.00

	testCases := []struct {
		name         string
		percentage   float64
		expectedAmnt int64
	}{
		{name: "should calculate simple percentage", percentage: 0.10, expectedAmnt: 1000},              // 10% of 100.00 is 10.00
		{name: "should calculate percentage with fractions", percentage: 0.155, expectedAmnt: 1550},     // 15.5% of 100.00 is 15.50
		{name: "should handle rounding down", percentage: 0.00114, expectedAmnt: 11},                    // 10000 * 0.00114 = 11.4 -> 11
		{name: "should handle rounding up", percentage: 0.00116, expectedAmnt: 12},                      // 10000 * 0.00116 = 11.6 -> 12
		{name: "should handle banker's rounding (half to even)", percentage: 0.00115, expectedAmnt: 12}, // 10000 * 0.00115 = 11.5 -> 12 (even)
		{name: "should handle banker's rounding (half to even)", percentage: 0.00125, expectedAmnt: 12}, // 10000 * 0.00125 = 12.5 -> 12 (even)
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			p, _ := wisp.NewPercentageFromFloat(tc.percentage)
			result := p.ApplyTo(money)
			s.Equal(tc.expectedAmnt, result.Amount())
			s.Equal(wisp.BRL, result.Currency())
		})
	}
}

func (s *PercentageSuite) TestPercentage_JSONMarshaling() {
	s.Run("should marshal and unmarshal correctly", func() {
		p, _ := wisp.NewPercentageFromFloat(0.50) // 50%
		data, err := json.Marshal(p)
		s.Require().NoError(err)
		s.Equal(`0.5`, string(data)) // Note: float marshaling might vary slightly

		var unmarshaledP wisp.Percentage
		err = json.Unmarshal(data, &unmarshaledP)
		s.Require().NoError(err)
		s.Equal(p, unmarshaledP)
	})

	s.Run("should fail to unmarshal an invalid value", func() {
		var p wisp.Percentage
		err := json.Unmarshal([]byte(`"not-a-number"`), &p)
		s.Require().Error(err)
	})
}

func (s *PercentageSuite) TestPercentage_DatabaseInterface() {
	p, _ := wisp.NewPercentageFromFloat(0.25) // 25% -> 2500

	s.Run("Value", func() {
		val, err := p.Value()
		s.Require().NoError(err)
		s.Equal(int64(2500), val)
	})

	s.Run("Scan", func() {
		var scannedP wisp.Percentage
		err := scannedP.Scan(int64(2500))
		s.Require().NoError(err)
		s.Equal(p, scannedP)

		err = scannedP.Scan(nil)
		s.Require().NoError(err)
		s.True(scannedP.IsZero())

		err = scannedP.Scan("invalid")
		s.Require().Error(err)
	})
}
