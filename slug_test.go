package wisp_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type SlugSuite struct {
	suite.Suite
}

func TestSlugSuite(t *testing.T) {
	suite.Run(t, new(SlugSuite))
}

func (s *SlugSuite) TestNewSlug() {
	testCases := []struct {
		name     string
		input    string
		expected wisp.Slug
		hasError bool
	}{
		{name: "simple case", input: "Hello World", expected: "hello-world"},
		{name: "with diacritics", input: "Geração de Leads É Importante", expected: "geracao-de-leads-e-importante"},
		{name: "with special characters", input: "Go! Is it 100% awesome?", expected: "go-is-it-100-percent-awesome"},
		{name: "with multiple spaces and hyphens", input: "  a--b   c- d ", expected: "a-b-c-d"},
		{name: "with leading and trailing hyphens", input: "-My--Slug-", expected: "my-slug"},
		{name: "should fail for empty input", input: "", hasError: true},
		{name: "should fail for input that results in empty slug", input: "!#*()_[]", hasError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			slug, err := wisp.NewSlug(tc.input)
			if tc.hasError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expected, slug)
			}
		})
	}
}

func (s *SlugSuite) TestSlug_Equals() {
	s1, _ := wisp.NewSlug("hello-world")
	s2, _ := wisp.NewSlug("Hello World")
	s3, _ := wisp.NewSlug("another-slug")

	s.True(s1.Equals(s2))
	s.False(s1.Equals(s3))
}
