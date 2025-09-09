package wisp_test

import (
	"testing"

	"github.com/marcelofabianov/wisp"
	"github.com/stretchr/testify/suite"
)

type UnitSuite struct {
	suite.Suite
}

func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}

func (s *UnitSuite) SetupTest() {
	wisp.ClearRegisteredUnits()
}

func (s *UnitSuite) TestRegisterAndValidateUnits() {
	s.Run("should return false for an unregistered unit", func() {
		s.False(wisp.Unit("KG").IsValid())
	})

	s.Run("should correctly register and validate new units", func() {
		const (
			Kilogram wisp.Unit = "KG"
			Box      wisp.Unit = "BOX"
		)

		wisp.RegisterUnits(Kilogram, Box)

		s.True(Kilogram.IsValid())
		s.True(Box.IsValid())
		s.False(wisp.Unit("LITER").IsValid())
	})

	s.Run("should normalize units during registration", func() {
		wisp.RegisterUnits(wisp.Unit("  un  "), wisp.Unit("lT"))

		s.True(wisp.Unit("UN").IsValid())
		s.True(wisp.Unit("LT").IsValid())
		s.False(wisp.Unit("un").IsValid(), "Validation should be case-sensitive after registration")
	})

	s.Run("should not register empty or blank units", func() {
		wisp.RegisterUnits("", "   ")
		s.False(wisp.Unit("").IsValid())
	})
}

func (s *UnitSuite) TestUnitStringer() {
	s.Equal("KG", wisp.Unit("KG").String())
}
