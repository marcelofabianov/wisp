package atomic_test

import (
	"testing"
	"time"

	"github.com/marcelofabianov/atomic"
	"github.com/stretchr/testify/suite"
)

type BirthDateSuite struct {
	suite.Suite
}

func TestBirthDateSuite(t *testing.T) {
	suite.Run(t, new(BirthDateSuite))
}

func (s *BirthDateSuite) TearDownTest() {
	// Reset legal age to default after tests that change it
	atomic.SetLegalAge(18)
}

func (s *BirthDateSuite) TestNewBirthDate() {
	s.Run("should create a valid birth date", func() {
		bd, err := atomic.NewBirthDate(1990, time.October, 20)
		s.Require().NoError(err)
		s.False(bd.IsZero())
	})

	s.Run("should fail for a date in the future", func() {
		tomorrow := atomic.Today().AddDays(1)
		_, err := atomic.NewBirthDate(tomorrow.Year(), tomorrow.Month(), tomorrow.Day())
		s.Require().Error(err)
	})

	s.Run("should parse a valid birth date string", func() {
		bd, err := atomic.ParseBirthDate("1990-10-20")
		s.Require().NoError(err)
		expected, _ := atomic.NewBirthDate(1990, time.October, 20)
		s.Equal(expected.Date(), bd.Date())
	})
}

func (s *BirthDateSuite) TestBirthDate_AgeAndIsOfAge() {
	bd, _ := atomic.NewBirthDate(2005, time.December, 15) // 19 years old on Sep 9, 2025
	today, _ := atomic.NewDate(2025, time.September, 9)

	s.Run("should calculate age correctly", func() {
		s.Equal(19, bd.Age(today))

		// Before anniversary in the same year
		todayBeforeAnniversary, _ := atomic.NewDate(2025, time.November, 10)
		s.Equal(19, bd.Age(todayBeforeAnniversary))

		// After anniversary in the same year
		todayAfterAnniversary, _ := atomic.NewDate(2025, time.December, 20)
		s.Equal(20, bd.Age(todayAfterAnniversary))
	})

	s.Run("should check IsOfAge with default legal age (18)", func() {
		s.True(bd.IsOfAge(today))

		bdMinor, _ := atomic.NewBirthDate(2008, time.January, 1)
		s.False(bdMinor.IsOfAge(today))
	})

	s.Run("should check IsOfAge with custom legal age (21)", func() {
		atomic.SetLegalAge(21)
		s.False(bd.IsOfAge(today))

		bdMajor, _ := atomic.NewBirthDate(2002, time.May, 5)
		s.True(bdMajor.IsOfAge(today))
	})
}

func (s *BirthDateSuite) TestBirthDate_Anniversary() {
	bd, _ := atomic.NewBirthDate(1990, time.October, 20)

	s.Run("when anniversary has not passed", func() {
		today, _ := atomic.NewDate(2025, time.September, 9)
		anniversary, _ := atomic.NewDate(2025, time.October, 20)

		s.Equal(anniversary, bd.AnniversaryThisYear(today))
		s.False(bd.HasAnniversaryPassed(today))
	})

	s.Run("when anniversary has passed", func() {
		today, _ := atomic.NewDate(2025, time.November, 1)
		anniversary, _ := atomic.NewDate(2025, time.October, 20)

		s.Equal(anniversary, bd.AnniversaryThisYear(today))
		s.True(bd.HasAnniversaryPassed(today))
	})

	s.Run("should handle leap year birthdays", func() {
		leapBd, _ := atomic.NewBirthDate(2000, time.February, 29)

		// In a non-leap year, time.Date rolls Feb 29 to Mar 1
		nonLeapYear, _ := atomic.NewDate(2025, time.January, 1)
		expectedAnniversary, _ := atomic.NewDate(2025, time.March, 1)
		s.Equal(expectedAnniversary, leapBd.AnniversaryThisYear(nonLeapYear))

		// In a leap year
		leapYear, _ := atomic.NewDate(2024, time.January, 1)
		expectedLeapAnniversary, _ := atomic.NewDate(2024, time.February, 29)
		s.Equal(expectedLeapAnniversary, leapBd.AnniversaryThisYear(leapYear))
	})
}
