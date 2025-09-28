package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type WeightSuite struct {
	suite.Suite
}

func TestWeightSuite(t *testing.T) {
	suite.Run(t, new(WeightSuite))
}

func (s *WeightSuite) TestNewWeight() {
	s.Run("should create weight from kilograms", func() {
		w, err := wisp.NewWeight(1.5, wisp.Kilogram)
		s.Require().NoError(err)
		val, _ := w.In(wisp.Gram)
		s.InDelta(1500, val, 0.001)
	})

	s.Run("should create weight from pounds", func() {
		w, err := wisp.NewWeight(2.20462, wisp.Pound)
		s.Require().NoError(err)
		val, _ := w.In(wisp.Kilogram)
		s.InDelta(1.0, val, 0.001)
	})

	s.Run("should fail for negative weight", func() {
		_, err := wisp.NewWeight(-10, wisp.Gram)
		s.Require().Error(err)
	})

	s.Run("should fail for unsupported unit", func() {
		_, err := wisp.NewWeight(10, "ton")
		s.Require().Error(err)
	})
}

func (s *WeightSuite) TestWeight_Conversions() {
	w, _ := wisp.NewWeight(1, wisp.Kilogram)

	valG, _ := w.In(wisp.Gram)
	s.InDelta(1000, valG, 0.001)

	valLb, _ := w.In(wisp.Pound)
	s.InDelta(2.20462, valLb, 0.001)

	valOz, _ := w.In(wisp.Ounce)
	s.InDelta(35.274, valOz, 0.001)
}

func (s *WeightSuite) TestWeight_Arithmetic() {
	w1, _ := wisp.NewWeight(100, wisp.Gram)
	w2, _ := wisp.NewWeight(1.5, wisp.Kilogram)

	s.Run("Add", func() {
		sum := w1.Add(w2)
		sumInKg, _ := sum.In(wisp.Kilogram)
		s.InDelta(1.6, sumInKg, 0.001)
	})

	s.Run("Subtract", func() {
		// Test 1: Simple subtraction
		diff1 := w2.Subtract(w1) // 1.5kg - 100g
		diff1InKg, _ := diff1.In(wisp.Kilogram)
		s.InDelta(1.4, diff1InKg, 0.001)
		s.False(diff1.IsNegative())

		// Test 2: Subtraction resulting in a negative weight
		diff2 := w1.Subtract(w2) // 100g - 1.5kg
		diff2InKg, _ := diff2.In(wisp.Kilogram)
		s.InDelta(-1.4, diff2InKg, 0.001)
		s.True(diff2.IsNegative())

		// Test 3: Subtraction resulting in zero
		diff3 := w1.Subtract(w1)
		s.True(diff3.Equals(wisp.ZeroWeight))
		s.False(diff3.IsNegative())
	})
}

func (s *WeightSuite) TestWeight_JSON_SQL() {
	w, _ := wisp.NewWeight(2.5, wisp.Kilogram)

	s.Run("JSON Marshaling", func() {
		data, err := json.Marshal(w)
		s.Require().NoError(err)

		dto := &struct{ Value float64 }{}
		json.Unmarshal(data, dto)

		s.InDelta(2.5, dto.Value, 0.001)
		s.JSONEq(`{"value": 2.5, "unit": "kg"}`, string(data))
	})

	s.Run("JSON Unmarshaling", func() {
		jsonData := `{"value": 1.2, "unit": "kg"}`
		var unmarshaledW wisp.Weight
		err := json.Unmarshal([]byte(jsonData), &unmarshaledW)
		s.Require().NoError(err)

		val, _ := unmarshaledW.In(wisp.Gram)
		s.InDelta(1200, val, 0.001)
	})

	s.Run("SQL Interface", func() {
		val, err := w.Value()
		s.Require().NoError(err)
		s.Equal(int64(2500000), val) // 2.5kg = 2500g = 2,500,000mg

		var scannedW wisp.Weight
		err = scannedW.Scan(int64(500000)) // 500,000mg = 500g
		s.Require().NoError(err)

		g, _ := scannedW.In(wisp.Gram)
		s.InDelta(500, g, 0.001)
	})
}
