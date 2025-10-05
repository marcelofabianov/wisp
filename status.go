package wisp

import (
	"strings"

	"github.com/marcelofabianov/fault"
)

// Status is a value object representing a state in a workflow or state machine (e.g., "PENDING", "ACTIVE", "INACTIVE").
// It ensures that only explicitly defined and registered statuses are used, providing type safety for state management.
//
// All valid statuses must be registered in a global allow-list before they can be used.
// Statuses are stored in a normalized (uppercase) format.
//
// Example:
//   wisp.RegisterStatuses("PENDING", "ACTIVE", "INACTIVE")
//   s, err := wisp.NewStatus("active")
//   isPending := s == "PENDING"
type Status string

// validStatuses holds the global set of registered statuses.
var validStatuses = make(map[Status]struct{})

// EmptyStatus represents the zero value for the Status type.
var EmptyStatus Status

// RegisterStatuses adds one or more statuses to the global registry of valid statuses.
// It normalizes them to uppercase and trims whitespace.
// This function should be called at application startup to define all possible statuses.
func RegisterStatuses(statuses ...Status) {
	for _, s := range statuses {
		normalized := Status(strings.ToUpper(strings.TrimSpace(string(s))))
		if normalized != "" {
			validStatuses[normalized] = struct{}{}
		}
	}
}

// NewStatus creates a new Status from a string.
// It normalizes the input to uppercase and validates it against the global registry.
// Returns an error if the status is not registered.
func NewStatus(value string) (Status, error) {
	normalized := Status(strings.ToUpper(strings.TrimSpace(value)))
	if normalized == EmptyStatus {
		return EmptyStatus, nil
	}

	if !normalized.IsValid() {
		return EmptyStatus, fault.New(
			"status is not registered as a valid status",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_status", value),
		)
	}
	return normalized, nil
}

// ClearRegisteredStatuses removes all statuses from the global registry.
// This is primarily for testing purposes to ensure a clean state.
func ClearRegisteredStatuses() {
	validStatuses = make(map[Status]struct{})
}

// String returns the status as a string.
func (s Status) String() string {
	return string(s)
}

// IsValid checks if the status is in the global registry of valid statuses.
func (s Status) IsValid() bool {
	_, ok := validStatuses[s]
	return ok
}

// IsZero returns true if the Status is the zero value.
func (s Status) IsZero() bool {
	return s == EmptyStatus
}
