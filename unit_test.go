package atomic_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/atomic"
)

type UnitSuite struct {
	suite.Suite
}

func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}

func (s *UnitSuite) SetupTest() {
	atomic.ClearRegisteredUnits()
}

func (s *UnitSuite) TestRegisterAndValidateUnits() {
	s.Run("should return false for an unregistered unit", func() {
		s.False(atomic.Unit("KG").IsValid())
	})

	s.Run("should correctly register and validate new units", func() {
		const (
			Kilogram atomic.Unit = "KG"
			Box      atomic.Unit = "BOX"
		)

		atomic.RegisterUnits(Kilogram, Box)

		s.True(Kilogram.IsValid())
		s.True(Box.IsValid())
		s.False(atomic.Unit("LITER").IsValid())
	})

	s.Run("should normalize units during registration", func() {
		atomic.RegisterUnits(atomic.Unit("  un  "), atomic.Unit("lT"))

		s.True(atomic.Unit("UN").IsValid())
		s.True(atomic.Unit("LT").IsValid())
		s.False(atomic.Unit("un").IsValid(), "Validation should be case-sensitive after registration")
	})

	s.Run("should not register empty or blank units", func() {
		atomic.RegisterUnits("", "   ")
		s.False(atomic.Unit("").IsValid())
	})
}

func (s *UnitSuite) TestUnitStringer() {
	s.Equal("KG", atomic.Unit("KG").String())
}
