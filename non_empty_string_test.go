package wisp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type NonEmptyStringSuite struct {
	suite.Suite
}

func TestNonEmptyStringSuite(t *testing.T) {
	suite.Run(t, new(NonEmptyStringSuite))
}

func (s *NonEmptyStringSuite) TestNewNonEmptyString() {
	s.Run("should create a valid non-empty string", func() {
		nes, err := wisp.NewNonEmptyString("  Hello World  ")
		s.Require().NoError(err)
		s.Equal("Hello World", nes.String())
		s.False(nes.IsZero())
	})

	s.Run("should fail for an empty string", func() {
		_, err := wisp.NewNonEmptyString("")
		s.Require().Error(err)
	})

	s.Run("should fail for a string with only whitespace", func() {
		_, err := wisp.NewNonEmptyString("   \t\n   ")
		s.Require().Error(err)
	})
}

func (s *NonEmptyStringSuite) TestNonEmptyString_JSON() {
	nes, _ := wisp.NewNonEmptyString("Test")

	data, err := json.Marshal(nes)
	s.Require().NoError(err)
	s.Equal(`"Test"`, string(data))

	var unmarshaledNES wisp.NonEmptyString
	err = json.Unmarshal(data, &unmarshaledNES)
	s.Require().NoError(err)
	s.Equal(nes, unmarshaledNES)

	err = json.Unmarshal([]byte(`"   "`), &unmarshaledNES)
	s.Require().Error(err)
}

func (s *NonEmptyStringSuite) TestNonEmptyString_SQL() {
	nes, _ := wisp.NewNonEmptyString("Test")

	val, err := nes.Value()
	s.Require().NoError(err)
	s.Equal("Test", val)

	var scannedNES wisp.NonEmptyString
	err = scannedNES.Scan(" Scanned ")
	s.Require().NoError(err)
	s.Equal("Scanned", scannedNES.String())

	err = scannedNES.Scan(" ")
	s.Require().Error(err)

	err = scannedNES.Scan(nil)
	s.Require().NoError(err)
	s.True(scannedNES.IsZero())
}
