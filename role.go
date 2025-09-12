package wisp

import (
	"strings"

	"github.com/marcelofabianov/fault"
)

type Role string

var validRoles = make(map[Role]struct{})
var EmptyRole Role

func RegisterRoles(roles ...Role) {
	for _, r := range roles {
		normalized := Role(strings.ToUpper(strings.TrimSpace(string(r))))
		if normalized != "" {
			validRoles[normalized] = struct{}{}
		}
	}
}

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

func ClearRegisteredRoles() {
	validRoles = make(map[Role]struct{})
}

func (r Role) String() string {
	return string(r)
}

func (r Role) IsValid() bool {
	_, ok := validRoles[r]
	return ok
}

func (r Role) IsZero() bool {
	return r == EmptyRole
}
