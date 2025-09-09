package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	wisp "github.com/marcelofabianov/wisp"
)

const (
	UnitKG wisp.Unit = "KG"
	UnitUN wisp.Unit = "UN"
	UnitL  wisp.Unit = "L"
)

type QuantitySuite struct {
	suite.Suite
}

func TestQuantitySuite(t *testing.T) {
	suite.Run(t, new(QuantitySuite))
}

func (s *QuantitySuite) SetupTest() {
	wisp.ClearRegisteredUnits()
	wisp.RegisterUnits(UnitKG, UnitUN, UnitL)
	wisp.SetDefaultPrecision(3) // Default to 3 decimal places for tests
}

func (s *QuantitySuite) TestNewQuantity() {
	s.Run("should create quantity with default precision", func() {
		q, err := wisp.NewQuantity(1.575, UnitKG)
		s.Require().NoError(err)
		s.Equal(int64(1575), q.Value())
		s.Equal(3, q.Precision())
		s.Equal(UnitKG, q.Unit())
	})

	s.Run("should create quantity with explicit precision", func() {
		q, err := wisp.NewQuantityWithPrecision(10, UnitUN, 0)
		s.Require().NoError(err)
		s.Equal(int64(10), q.Value())
		s.Equal(0, q.Precision())
	})

	s.Run("should fail for unregistered unit", func() {
		_, err := wisp.NewQuantity(1, "BOX")
		s.Require().Error(err)
	})
}

func (s *QuantitySuite) TestQuantity_Add() {
	q1, _ := wisp.NewQuantityWithPrecision(10.5, UnitKG, 2)
	q2, _ := wisp.NewQuantityWithPrecision(2.5, UnitKG, 2)
	q3, _ := wisp.NewQuantityWithPrecision(5, UnitL, 2)
	q4, _ := wisp.NewQuantityWithPrecision(2.5, UnitKG, 3) // Different precision

	s.Run("should add quantities with same unit and precision", func() {
		result, err := q1.Add(q2)
		s.Require().NoError(err)
		s.Equal(int64(1300), result.Value()) // 1050 + 250
		s.InDelta(13.0, result.Float64(), 0.001)
	})

	s.Run("should fail to add quantities with different units", func() {
		_, err := q1.Add(q3)
		s.Require().Error(err)
	})

	s.Run("should fail to add quantities with different precisions", func() {
		_, err := q1.Add(q4)
		s.Require().Error(err)
	})
}

func (s *QuantitySuite) TestQuantity_MultiplyByMoney() {
	price, _ := wisp.NewMoney(1031, wisp.BRL) // R$ 10.31
	qty, _ := wisp.NewQuantityWithPrecision(1.57, UnitKG, 2)

	result, err := qty.MultiplyByMoney(price)
	s.Require().NoError(err)

	// 1.57 * 1031 = 1618.67 -> rounds to 1619
	expected, _ := wisp.NewMoney(1619, wisp.BRL)
	s.True(expected.Equals(result))
}

func (s *QuantitySuite) TestQuantity_JSONMarshaling() {
	s.Run("should marshal and unmarshal correctly", func() {
		q, _ := wisp.NewQuantityWithPrecision(12.75, UnitL, 2)
		data, err := json.Marshal(q)
		s.Require().NoError(err)
		s.JSONEq(`{"value": 12.75, "unit": "L"}`, string(data))

		var unmarshaledQ wisp.Quantity
		err = json.Unmarshal(data, &unmarshaledQ)
		s.Require().NoError(err)
		s.Equal(q.Value(), unmarshaledQ.Value())
		s.Equal(q.Unit(), unmarshaledQ.Unit())
		s.Equal(q.Precision(), unmarshaledQ.Precision())
	})

	s.Run("should fail unmarshal with unregistered unit", func() {
		invalidJSON := `{"value": 10.0, "unit": "XXX"}`
		var q wisp.Quantity
		err := json.Unmarshal([]byte(invalidJSON), &q)
		s.Require().Error(err)
	})
}
