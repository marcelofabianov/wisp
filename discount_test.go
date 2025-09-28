package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type DiscountSuite struct {
	suite.Suite
}

func TestDiscountSuite(t *testing.T) {
	suite.Run(t, new(DiscountSuite))
}

func (s *DiscountSuite) TestNewDiscount() {
	s.Run("should create a valid fixed discount", func() {
		m, _ := wisp.NewMoney(1000, wisp.BRL)
		d, err := wisp.NewFixedDiscount(m)
		s.Require().NoError(err)
		s.False(d.IsZero())
	})

	s.Run("should fail for a negative fixed discount", func() {
		m, _ := wisp.NewMoney(-1, wisp.BRL)
		_, err := wisp.NewFixedDiscount(m)
		s.Require().Error(err)
	})

	s.Run("should create a valid percentage discount", func() {
		p, _ := wisp.NewPercentageFromFloat(0.15) // 15%
		d, err := wisp.NewPercentageDiscount(p)
		s.Require().NoError(err)
		s.False(d.IsZero())
	})

	s.Run("should fail for a percentage discount over 100%", func() {
		p, _ := wisp.NewPercentageFromFloat(1.1) // 110%
		_, err := wisp.NewPercentageDiscount(p)
		s.Require().Error(err)
	})
}

func (s *DiscountSuite) TestDiscount_ApplyTo() {
	originalPrice, _ := wisp.NewMoney(10000, wisp.BRL) // R$ 100,00

	s.Run("should apply a fixed discount", func() {
		discountValue, _ := wisp.NewMoney(1500, wisp.BRL) // R$ 15,00
		d, _ := wisp.NewFixedDiscount(discountValue)

		finalPrice, err := d.ApplyTo(originalPrice)
		s.Require().NoError(err)
		s.Equal(int64(8500), finalPrice.Amount())
	})

	s.Run("should apply a percentage discount", func() {
		p, _ := wisp.NewPercentageFromFloat(0.20) // 20%
		d, _ := wisp.NewPercentageDiscount(p)

		finalPrice, err := d.ApplyTo(originalPrice)
		s.Require().NoError(err)
		s.Equal(int64(8000), finalPrice.Amount()) // 10000 - (10000 * 0.20) = 8000
	})

	s.Run("should not result in a negative price", func() {
		discountValue, _ := wisp.NewMoney(12000, wisp.BRL) // R$ 120,00
		d, _ := wisp.NewFixedDiscount(discountValue)

		finalPrice, err := d.ApplyTo(originalPrice)
		s.Require().NoError(err)
		s.Equal(int64(0), finalPrice.Amount())
	})
}

func (s *DiscountSuite) TestDiscount_JSON() {
	s.Run("should marshal and unmarshal a fixed discount", func() {
		m, _ := wisp.NewMoney(500, wisp.BRL)
		d, _ := wisp.NewFixedDiscount(m)

		data, err := json.Marshal(d)
		s.Require().NoError(err)
		s.JSONEq(`{"type": "fixed", "value": {"amount": 500, "currency": "BRL"}}`, string(data))

		var unmarshaledD wisp.Discount
		err = json.Unmarshal(data, &unmarshaledD)
		s.Require().NoError(err)
		s.Equal(d.String(), unmarshaledD.String())
	})

	s.Run("should marshal and unmarshal a percentage discount", func() {
		p, _ := wisp.NewPercentageFromFloat(0.25)
		d, _ := wisp.NewPercentageDiscount(p)

		data, err := json.Marshal(d)
		s.Require().NoError(err)
		s.JSONEq(`{"type": "percentage", "value": 0.25}`, string(data))

		var unmarshaledD wisp.Discount
		err = json.Unmarshal(data, &unmarshaledD)
		s.Require().NoError(err)
		s.Equal(d.String(), unmarshaledD.String())
	})
}
