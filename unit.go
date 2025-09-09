package wisp

import "strings"

type Unit string

var validUnits = make(map[Unit]struct{})

func RegisterUnits(units ...Unit) {
	for _, u := range units {
		normalized := Unit(strings.ToUpper(strings.TrimSpace(string(u))))
		if normalized != "" {
			validUnits[normalized] = struct{}{}
		}
	}
}

func ClearRegisteredUnits() {
	validUnits = make(map[Unit]struct{})
}

func (u Unit) String() string {
	return string(u)
}

func (u Unit) IsValid() bool {
	_, ok := validUnits[u]
	return ok
}
