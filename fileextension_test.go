package wisp_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type FileExtensionSuite struct {
	suite.Suite
}

func TestFileExtensionSuite(t *testing.T) {
	suite.Run(t, new(FileExtensionSuite))
}

func (s *FileExtensionSuite) SetupTest() {
	wisp.ClearRegisteredFileExtensions()
}

func (s *FileExtensionSuite) TestNewFileExtensionWithRegistry() {
	wisp.RegisterFileExtensions("pdf", ".JPG", "  docx  ")

	s.Run("should create a valid extension that is registered", func() {
		ext, err := wisp.NewFileExtension("pdf")
		s.Require().NoError(err)
		s.Equal(wisp.FileExtension("pdf"), ext)
	})

	s.Run("should create and normalize a registered extension", func() {
		ext, err := wisp.NewFileExtension(".JpG")
		s.Require().NoError(err)
		s.Equal(wisp.FileExtension("jpg"), ext)
	})

	s.Run("should fail for an extension that is not registered", func() {
		_, err := wisp.NewFileExtension("xml")
		s.Require().Error(err)
		s.Contains(err.Error(), "not registered in the allowed list")
	})

	s.Run("should fail for an empty string", func() {
		_, err := wisp.NewFileExtension("")
		s.Require().Error(err)
	})
}

func (s *FileExtensionSuite) TestIsRegistered() {
	wisp.RegisterFileExtensions("txt", "log")

	s.True(wisp.FileExtension("txt").IsRegistered())
	s.False(wisp.FileExtension("csv").IsRegistered())
}
