package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/marcelofabianov/fault"
)

// AuditUser is a value object representing the identifier of a user or system performing an action.
// It is used in the `Audit` struct to track who created, updated, or deleted an entity.
// An AuditUser can be one of two things:
//  1. A valid email address, representing a human user.
//  2. The special literal string "system", representing an automated process.
//
// The value is normalized to lowercase.
//
// Example:
//   user, _ := NewAuditUser("jane.doe@example.com")
//   system := NewSystemAuditUser() // "system"
type AuditUser string

// EmptyAuditUser represents the zero value for the AuditUser type.
var EmptyAuditUser AuditUser

// SystemAuditUser is a special constant representing an automated system action.
const SystemAuditUser AuditUser = "system"

// NewSystemAuditUser returns an AuditUser representing the system.
func NewSystemAuditUser() AuditUser {
	return SystemAuditUser
}

// NewAuditUser creates a new AuditUser from a string.
// It validates that the input is either a valid email address or the string "system".
// The input is normalized to lowercase.
// Returns an error if the input is not valid.
func NewAuditUser(input string) (AuditUser, error) {
	user := AuditUser(strings.ToLower(strings.TrimSpace(input)))

	if user.IsZero() || user.IsSystem() {
		return user, nil
	}

	if _, err := NewEmail(string(user)); err != nil {
		return EmptyAuditUser, fault.New(
			"audit user must be a valid email or the literal 'system'",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_user", input),
		)
	}

	return user, nil
}

// String returns the audit user identifier as a string.
func (au AuditUser) String() string {
	return string(au)
}

// IsZero returns true if the AuditUser is the zero value.
func (au AuditUser) IsZero() bool {
	return au == EmptyAuditUser
}

// IsSystem returns true if the AuditUser represents the system.
func (au AuditUser) IsSystem() bool {
	return au == SystemAuditUser
}

// IsEmail returns true if the AuditUser represents a user with an email address.
func (au AuditUser) IsEmail() bool {
	return !au.IsZero() && !au.IsSystem()
}

// Email returns the email address if the AuditUser is an email, and a boolean indicating success.
// If the AuditUser is the system or zero, it returns an empty Email and false.
func (au AuditUser) Email() (Email, bool) {
	if !au.IsEmail() {
		return EmptyEmail, false
	}

	email, _ := NewEmail(string(au))
	return email, true
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the AuditUser to its string representation.
func (au AuditUser) MarshalJSON() ([]byte, error) {
	return json.Marshal(au.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into an AuditUser, with validation.
func (au *AuditUser) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "AuditUser must be a valid JSON string", fault.WithCode(fault.Invalid))
	}
	user, err := NewAuditUser(s)
	if err != nil {
		return err
	}
	*au = user
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the AuditUser as a string.
func (au AuditUser) Value() (driver.Value, error) {
	if au.IsZero() {
		return nil, nil
	}
	return au.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a string or byte slice from the database and converts it into an AuditUser, with validation.
func (au *AuditUser) Scan(src interface{}) error {
	if src == nil {
		*au = EmptyAuditUser
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for AuditUser", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}

	user, err := NewAuditUser(s)
	if err != nil {
		return err
	}
	*au = user
	return nil
}
