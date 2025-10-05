package wisp

import (
	"strings"

	"github.com/marcelofabianov/fault"
)

// Role is a value object representing a user role within the system (e.g., "ADMIN", "USER").
// It ensures that only explicitly defined and registered roles are used, providing type safety
// for authorization and access control logic.
//
// All valid roles must be registered in a global allow-list before they can be used.
// Roles are stored in a normalized (uppercase) format.
//
// Example:
//   wisp.RegisterRoles("ADMIN", "USER", "GUEST")
//   r, err := wisp.NewRole("admin")
//   isAdmin := r == "ADMIN"
type Role string

// validRoles holds the global set of registered roles.
var validRoles = make(map[Role]struct{})

// EmptyRole represents the zero value for the Role type.
var EmptyRole Role

// RegisterRoles adds one or more roles to the global registry of valid roles.
// It normalizes them to uppercase and trims whitespace.
// This function should be called at application startup to define all possible user roles.
func RegisterRoles(roles ...Role) {
	for _, r := range roles {
		normalized := Role(strings.ToUpper(strings.TrimSpace(string(r))))
		if normalized != "" {
			validRoles[normalized] = struct{}{}
		}
	}
}

// NewRole creates a new Role from a string.
// It normalizes the input to uppercase and validates it against the global registry.
// Returns an error if the role is not registered.
func NewRole(value string) (Role, error) {
	normalized := Role(strings.ToUpper(strings.TrimSpace(value)))
	if normalized == EmptyRole {
		return EmptyRole, nil
	}

	if !normalized.IsValid() {
		return EmptyRole, fault.New(
			"role is not registered as a valid role",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_role", value),
		)
	}
	return normalized, nil
}

// ClearRegisteredRoles removes all roles from the global registry.
// This is primarily for testing purposes to ensure a clean state.
func ClearRegisteredRoles() {
	validRoles = make(map[Role]struct{})
}

// String returns the role as a string.
func (r Role) String() string {
	return string(r)
}

// IsValid checks if the role is in the global registry of valid roles.
func (r Role) IsValid() bool {
	_, ok := validRoles[r]
	return ok
}

// IsZero returns true if the Role is the zero value.
func (r Role) IsZero() bool {
	return r == EmptyRole
}
