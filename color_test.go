package wisp_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type ColorSuite struct {
	suite.Suite
}

func TestColorSuite(t *testing.T) {
	suite.Run(t, new(ColorSuite))
}

func (s *ColorSuite) TestParseColor() {
	s.Run("should parse a 6-digit hex color", func() {
		c, err := wisp.ParseColor("#ff0000")
		s.Require().NoError(err)
		r, g, b, a := c.RGBA()
		s.Equal(uint8(255), r)
		s.Equal(uint8(0), g)
		s.Equal(uint8(0), b)
		s.Equal(uint8(255), a)
	})

	s.Run("should parse a 3-digit hex color", func() {
		c, err := wisp.ParseColor("#f0c")
		s.Require().NoError(err)
		s.Equal("#ff00cc", c.Hex())
	})

	s.Run("should parse and normalize an uppercase color with spaces", func() {
		c, err := wisp.ParseColor("  #00FF00  ")
		s.Require().NoError(err)
		s.Equal("#00ff00", c.Hex())
	})

	s.Run("should fail for invalid formats", func() {
		_, err := wisp.ParseColor("ff0000") // Missing '#'
		s.Require().Error(err)

		_, err = wisp.ParseColor("#12345") // Invalid length
		s.Require().Error(err)

		_, err = wisp.ParseColor("#gg0000") // Invalid hex characters
		s.Require().Error(err)
	})
}

func (s *ColorSuite) TestColor_Methods() {
	c, _ := wisp.ParseColor("#336699")

	s.Run("Hex and String", func() {
		s.Equal("#336699", c.Hex())
		s.Equal("#336699", c.String())
	})

	s.Run("RGBA", func() {
		r, g, b, a := c.RGBA()
		s.Equal(uint8(0x33), r)
		s.Equal(uint8(0x66), g)
		s.Equal(uint8(0x99), b)
		s.Equal(uint8(0xff), a)
	})

	s.Run("IsZero", func() {
		s.False(c.IsZero())
		s.True(wisp.ZeroColor.IsZero())
	})
}
