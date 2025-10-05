package wisp

import "strings"

// Unit is a value object representing a generic unit of measure (e.g., "BOX", "KG", "LITER").
// It provides a flexible way to define and validate custom units for use with the Quantity type.
//
// All valid units must be registered in a global allow-list before they can be used.
// This ensures that only predefined, domain-specific units are permitted.
//
// Example:
//   wisp.RegisterUnits("BOX", "PALLET")
//   unit := wisp.Unit("BOX")
//   isValid := unit.IsValid() // true
type Unit string

// validUnits holds the global set of registered units of measure.
var validUnits = make(map[Unit]struct{})

// RegisterUnits adds one or more units to the global registry of valid units.
// It normalizes the units to uppercase and trims whitespace.
// This function should be called at application startup to define all possible units.
func RegisterUnits(units ...Unit) {
	for _, u := range units {
		normalized := Unit(strings.ToUpper(strings.TrimSpace(string(u))))
		if normalized != "" {
			validUnits[normalized] = struct{}{}
		}
	}
}

// ClearRegisteredUnits removes all units from the global registry.
// This is primarily for testing purposes to ensure a clean state.
func ClearRegisteredUnits() {
	validUnits = make(map[Unit]struct{})
}

// String returns the unit as a string.
func (u Unit) String() string {
	return string(u)
}

// IsValid checks if the unit is in the global registry of valid units.
func (u Unit) IsValid() bool {
	_, ok := validUnits[u]
	return ok
}
