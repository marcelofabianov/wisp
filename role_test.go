package wisp_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type RoleSuite struct {
	suite.Suite
}

func TestRoleSuite(t *testing.T) {
	suite.Run(t, new(RoleSuite))
}

func (s *RoleSuite) SetupTest() {
	wisp.ClearRegisteredRoles()
}

func (s *RoleSuite) TestRegisterAndValidateRoles() {
	s.Run("should return false for an unregistered role", func() {
		s.False(wisp.Role("ADMIN").IsValid())
	})

	s.Run("should correctly register and validate new roles", func() {
		const (
			RoleAdmin   wisp.Role = "ADMIN"
			RoleStudent wisp.Role = "STUDENT"
		)

		wisp.RegisterRoles(RoleAdmin, RoleStudent)

		s.True(RoleAdmin.IsValid())
		s.True(RoleStudent.IsValid())
		s.False(wisp.Role("TEACHER").IsValid())
	})

	s.Run("should normalize roles during registration", func() {
		wisp.RegisterRoles(wisp.Role("  editor  "), wisp.Role("viewer"))

		s.True(wisp.Role("EDITOR").IsValid())
		s.True(wisp.Role("VIEWER").IsValid())
		s.False(wisp.Role("editor").IsValid(), "Validation should be case-sensitive after registration")
	})
}

func (s *RoleSuite) TestRoleStringer() {
	s.Equal("ADMIN", wisp.Role("ADMIN").String())
}
