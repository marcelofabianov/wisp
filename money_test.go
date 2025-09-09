package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/suite"

	wisp "github.com/marcelofabianov/wisp"
)

type MoneySuite struct {
	suite.Suite
}

func TestMoneySuite(t *testing.T) {
	suite.Run(t, new(MoneySuite))
}

func (s *MoneySuite) TestNewMoney() {
	s.Run("should create a new money object successfully", func() {
		m, err := wisp.NewMoney(1050, wisp.BRL)
		s.Require().NoError(err)
		s.False(m.IsZero())
		s.Equal(int64(1050), m.Amount())
		s.Equal(wisp.BRL, m.Currency())
	})

	s.Run("should fail with an invalid currency", func() {
		m, err := wisp.NewMoney(1000, wisp.Currency("XYZ"))
		s.Require().Error(err)
		s.True(m.IsZero())
		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok)
		s.Equal(fault.Invalid, faultErr.Code)
	})

	s.Run("should fail with a zero currency", func() {
		m, err := wisp.NewMoney(1000, wisp.EmptyCurrency)
		s.Require().Error(err)
		s.True(m.IsZero())
	})
}

func (s *MoneySuite) TestMoney_Comparison() {
	brl100, _ := wisp.NewMoney(10000, wisp.BRL)
	brl200, _ := wisp.NewMoney(20000, wisp.BRL)
	usd100, _ := wisp.NewMoney(10000, wisp.USD)

	s.Run("Equals", func() {
		brl100Clone, _ := wisp.NewMoney(10000, wisp.BRL)
		s.True(brl100.Equals(brl100Clone))
		s.False(brl100.Equals(brl200))
		s.False(brl100.Equals(usd100))
		s.False(brl100.Equals(wisp.ZeroMoney))
	})

	s.Run("GreaterThan and LessThan", func() {
		gt, err := brl200.GreaterThan(brl100)
		s.Require().NoError(err)
		s.True(gt)

		lt, err := brl100.LessThan(brl200)
		s.Require().NoError(err)
		s.True(lt)

		_, err = brl100.GreaterThan(usd100)
		s.Require().Error(err)
		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok)
		s.Equal(fault.DomainViolation, faultErr.Code)
	})
}

func (s *MoneySuite) TestMoney_Arithmetic() {
	brl100, _ := wisp.NewMoney(10000, wisp.BRL)
	brl50, _ := wisp.NewMoney(5000, wisp.BRL)
	usd50, _ := wisp.NewMoney(5000, wisp.USD)

	s.Run("Add", func() {
		result, err := brl100.Add(brl50)
		s.Require().NoError(err)
		s.Equal(int64(15000), result.Amount())

		_, err = brl100.Add(usd50)
		s.Require().Error(err)
	})

	s.Run("Subtract", func() {
		result, err := brl100.Subtract(brl50)
		s.Require().NoError(err)
		s.Equal(int64(5000), result.Amount())

		_, err = brl100.Subtract(usd50)
		s.Require().Error(err)
	})

	s.Run("Multiply", func() {
		result := brl50.Multiply(3)
		s.Equal(int64(15000), result.Amount())
		s.Equal(wisp.BRL, result.Currency())
	})
}

func (s *MoneySuite) TestMoney_Split() {
	s.Run("should split an even amount correctly", func() {
		m, _ := wisp.NewMoney(10000, wisp.BRL)
		parts, err := m.Split(4)
		s.Require().NoError(err)
		s.Len(parts, 4)
		s.Equal(int64(2500), parts[0].Amount())
		s.Equal(int64(2500), parts[1].Amount())
		s.Equal(int64(2500), parts[2].Amount())
		s.Equal(int64(2500), parts[3].Amount())
	})

	s.Run("should split an uneven amount distributing the remainder", func() {
		m, _ := wisp.NewMoney(10000, wisp.BRL) // 100.00
		parts, err := m.Split(3)
		s.Require().NoError(err)
		s.Len(parts, 3)

		// 10000 / 3 = 3333 with remainder 1
		s.Equal(int64(3334), parts[0].Amount()) // 33.34
		s.Equal(int64(3333), parts[1].Amount()) // 33.33
		s.Equal(int64(3333), parts[2].Amount()) // 33.33

		// Verify total amount is conserved
		total := parts[0].Amount() + parts[1].Amount() + parts[2].Amount()
		s.Equal(m.Amount(), total)
	})

	s.Run("should fail for invalid split count", func() {
		m, _ := wisp.NewMoney(10000, wisp.BRL)
		_, err := m.Split(0)
		s.Require().Error(err)
		_, err = m.Split(-1)
		s.Require().Error(err)
	})
}

func (s *MoneySuite) TestMoney_JSONMarshaling() {
	s.Run("should marshal and unmarshal correctly", func() {
		m, _ := wisp.NewMoney(12345, wisp.EUR)
		expectedJSON := `{"amount":12345,"currency":"EUR"}`

		data, err := json.Marshal(m)
		s.Require().NoError(err)
		s.JSONEq(expectedJSON, string(data))

		var unmarshaledMoney wisp.Money
		err = json.Unmarshal(data, &unmarshaledMoney)
		s.Require().NoError(err)
		s.True(m.Equals(unmarshaledMoney))
	})

	s.Run("should fail to unmarshal with invalid currency", func() {
		invalidJSON := `{"amount":1000,"currency":"XYZ"}`
		var m wisp.Money
		err := json.Unmarshal([]byte(invalidJSON), &m)
		s.Require().Error(err)
		faultErr, ok := err.(*fault.Error)
		s.Require().True(ok)
		s.Equal(fault.Invalid, faultErr.Code)
	})

	s.Run("should fail to unmarshal with missing currency", func() {
		invalidJSON := `{"amount":1000}`
		var m wisp.Money
		err := json.Unmarshal([]byte(invalidJSON), &m)
		s.Require().Error(err)
	})
}

func (s *MoneySuite) TestMoney_Representation() {
	m, _ := wisp.NewMoney(9990, wisp.USD)

	s.Run("Float64", func() {
		s.InDelta(99.90, m.Float64(), 0.001)
	})

	s.Run("String", func() {
		s.Equal("USD 99.90", m.String())
	})
}
