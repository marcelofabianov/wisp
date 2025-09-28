package wisp_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type FlagSuite struct {
	suite.Suite
}

func TestFlagSuite(t *testing.T) {
	suite.Run(t, new(FlagSuite))
}

func (s *FlagSuite) TestFlag_String() {
	s.Run("should create and validate a string flag", func() {
		statusPublished := "published"
		statusDraft := "draft"

		flag1, err := wisp.NewFlag(statusPublished, statusPublished, statusDraft)
		s.Require().NoError(err)
		s.True(flag1.IsPrimary())
		s.False(flag1.IsSecondary())
		s.True(flag1.Is(statusPublished))
		s.False(flag1.Is(statusDraft))
		s.Equal(statusPublished, flag1.Get())

		flag2, err := wisp.NewFlag(statusDraft, statusPublished, statusDraft)
		s.Require().NoError(err)
		s.False(flag2.IsPrimary())
		s.True(flag2.IsSecondary())
		s.False(flag2.Is("published"))
		s.True(flag2.Is("draft"))

		_, err = wisp.NewFlag("archived", statusPublished, statusDraft)
		s.Require().Error(err)
	})
}

func (s *FlagSuite) TestFlag_Int() {
	s.Run("should create and validate an int flag", func() {
		active := 1
		inactive := 0

		flag, err := wisp.NewFlag(active, active, inactive)
		s.Require().NoError(err)
		s.True(flag.IsPrimary())
		s.True(flag.Is(1))
		s.False(flag.Is(0))
		s.Equal(active, flag.Get())
	})
}
