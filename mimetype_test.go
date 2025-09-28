package wisp_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type MIMETypeSuite struct {
	suite.Suite
}

func TestMIMETypeSuite(t *testing.T) {
	suite.Run(t, new(MIMETypeSuite))
}

func (s *MIMETypeSuite) SetupTest() {
	wisp.ClearRegisteredMIMETypes()
}

func (s *MIMETypeSuite) TestNewMIMEType() {
	wisp.RegisterMIMETypes("image/jpeg", "application/pdf")

	s.Run("should create a valid and registered MIME type", func() {
		mt, err := wisp.NewMIMEType("image/jpeg")
		s.Require().NoError(err)
		s.Equal(wisp.MIMEType("image/jpeg"), mt)
	})

	s.Run("should normalize and create a valid type", func() {
		mt, err := wisp.NewMIMEType("  APPLICATION/PDF  ")
		s.Require().NoError(err)
		s.Equal(wisp.MIMEType("application/pdf"), mt)
	})

	s.Run("should fail for an unregistered type", func() {
		_, err := wisp.NewMIMEType("text/plain")
		s.Require().Error(err)
		s.Contains(err.Error(), "not registered")
	})

	s.Run("should fail for an invalid format", func() {
		_, err := wisp.NewMIMEType("application")
		s.Require().Error(err)
		s.Contains(err.Error(), "type/subtype")

		_, err = wisp.NewMIMEType("application/")
		s.Require().Error(err)
		s.Contains(err.Error(), "type/subtype")

		_, err = wisp.NewMIMEType("/pdf")
		s.Require().Error(err)
		s.Contains(err.Error(), "type/subtype")
	})
}

func (s *MIMETypeSuite) TestMIMEType_Methods() {
	wisp.RegisterMIMETypes("application/vnd.api+json")
	mt, _ := wisp.NewMIMEType("application/vnd.api+json")

	s.Equal("application", mt.Type())
	s.Equal("vnd.api+json", mt.SubType())
	s.True(mt.IsRegistered())
}
