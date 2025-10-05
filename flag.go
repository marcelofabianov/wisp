package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

// Flag is a generic value object for representing a binary state.
// It is useful for any scenario that requires a choice between two specific values,
// such as `true`/`false`, `"ACTIVE"`/`"INACTIVE"`, or `1`/`0`.
// The type parameter `T` can be any comparable type (string, int, bool, etc.).
//
// It ensures that the flag's value is always one of the two allowed states (primary or secondary).
//
// Example:
//   // Using bool
//   activeFlag, _ := NewFlag(true, true, false)
//   isActive := activeFlag.IsPrimary() // true
//
//   // Using string
//   statusFlag, _ := NewFlag("ACTIVE", "ACTIVE", "INACTIVE")
//   isInactive := statusFlag.IsSecondary() // false
type Flag[T comparable] struct {
	value          T
	primaryValue   T
	secondaryValue T
}

// NewFlag creates a new Flag.
// It requires the current value, a primary value, and a secondary value.
// It returns an error if the current value is not one of the primary or secondary values.
func NewFlag[T comparable](currentValue, primaryValue, secondaryValue T) (Flag[T], error) {
	if currentValue != primaryValue && currentValue != secondaryValue {
		return Flag[T]{}, fault.New(
			"current value is not one of the allowed flag values",
			fault.WithCode(fault.Invalid),
			fault.WithContext("current_value", currentValue),
			fault.WithContext("allowed_values", []T{primaryValue, secondaryValue}),
		)
	}

	return Flag[T]{
		value:          currentValue,
		primaryValue:   primaryValue,
		secondaryValue: secondaryValue,
	}, nil
}

// Get returns the underlying value of the flag.
func (f Flag[T]) Get() T {
	return f.value
}

// IsPrimary returns true if the flag's value is equal to its primary value.
func (f Flag[T]) IsPrimary() bool {
	return f.value == f.primaryValue
}

// IsSecondary returns true if the flag's value is equal to its secondary value.
func (f Flag[T]) IsSecondary() bool {
	return f.value == f.secondaryValue
}

// Is checks if the flag's value is equal to a given value.
func (f Flag[T]) Is(value T) bool {
	return f.value == value
}

// String returns the string representation of the flag's underlying value.
func (f Flag[T]) String() string {
	return fmt.Sprintf("%v", f.value)
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the flag's underlying value to JSON.
func (f Flag[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.value)
}

// Value implements the driver.Valuer interface for database storage.
// It returns the flag's underlying value.
func (f Flag[T]) Value() (driver.Value, error) {
	return f.value, nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a value of type T from the database and sets it as the flag's value.
// Note: This does not re-validate against primary/secondary values, assuming database integrity.
func (f *Flag[T]) Scan(src interface{}) error {
	if src == nil {
		var zero T
		f.value = zero
		return nil
	}

	val, ok := src.(T)
	if !ok {
		return fault.New("unsupported scan type for Flag", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}
	f.value = val

	return nil
}
