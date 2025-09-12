package wisp

import "strings"

type Role string

var validRoles = make(map[Role]struct{})

func RegisterRoles(roles ...Role) {
	for _, r := range roles {
		normalized := Role(strings.ToUpper(strings.TrimSpace(string(r))))
		if normalized != "" {
			validRoles[normalized] = struct{}{}
		}
	}
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
