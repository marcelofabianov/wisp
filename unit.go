package wisp

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/marcelofabianov/fault"
)

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

// Value implements the driver.Valuer interface for database storage.
// It returns the unit as a string or nil if it's empty.
func (u Unit) Value() (driver.Value, error) {
	if u == "" {
		return nil, nil
	}
	return u.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts string or []byte values and validates them as a Unit.
func (u *Unit) Scan(src interface{}) error {
	if src == nil {
		*u = ""
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New(
			"unsupported scan type for Unit",
			fault.WithCode(fault.Invalid),
			fault.WithContext("received_type", fmt.Sprintf("%T", src)),
		)
	}

	*u = Unit(strings.ToUpper(strings.TrimSpace(s)))
	return nil
}
