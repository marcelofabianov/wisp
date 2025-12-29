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
	const (
		RoleAdmin wisp.Role = "ADMIN"
		RoleUser  wisp.Role = "USER"
	)
	wisp.RegisterRoles(RoleAdmin, RoleUser)

	s.True(RoleAdmin.IsValid())
	s.True(wisp.Role("USER").IsValid())
	s.False(wisp.Role("GUEST").IsValid())
	s.False(wisp.EmptyRole.IsValid())
}

func (s *RoleSuite) TestNewRole() {
	wisp.RegisterRoles("ADMIN", "editor")

	testCases := []struct {
		name        string
		input       string
		expected    wisp.Role
		expectError bool
	}{
		{name: "should create a valid role", input: "ADMIN", expected: "ADMIN"},
		{name: "should create a valid role preserving case", input: "  editor  ", expected: "editor"},
		{name: "should create an empty role from an empty string", input: "", expected: wisp.EmptyRole},
		{name: "should create an empty role from a blank string", input: "   ", expected: wisp.EmptyRole},
		{name: "should fail for an unregistered role", input: "GUEST", expectError: true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			role, err := wisp.NewRole(tc.input)
			if tc.expectError {
				s.Require().Error(err)
				s.True(role.IsZero())
			} else {
				s.Require().NoError(err)
				s.Equal(tc.expected, role)
			}
		})
	}
}

func (s *RoleSuite) TestRole_IsZero() {
	wisp.RegisterRoles("ADMIN")

	adminRole, err := wisp.NewRole("ADMIN")
	s.Require().NoError(err, "NewRole should succeed for a registered role")

	s.False(adminRole.IsZero())
	s.True(wisp.EmptyRole.IsZero())
}
