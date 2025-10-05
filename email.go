package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/mail"
	"strings"

	"github.com/marcelofabianov/fault"
)

// MaxEmailLength is the standard maximum length for an email address (254 characters).
const (
	MaxEmailLength = 254
)

// Email is a value object representing a validated and normalized email address.
// It ensures that the email has a valid format according to RFC 5322, is within the standard length limit,
// and is stored in a consistent, lowercase format.
//
// The zero value is EmptyEmail.
//
// Examples:
//   e, err := NewEmail("  Test@Example.COM ")
//   fmt.Println(e) // "test@example.com"
type Email string

// EmptyEmail represents the zero value for the Email type.
var EmptyEmail Email

// parseEmail contains the core logic for validating and normalizing an email string.
func parseEmail(emailStr string) (Email, error) {
	trimmedEmail := strings.TrimSpace(emailStr)

	if trimmedEmail == "" {
		return EmptyEmail, fault.New(
			"email address cannot be empty",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input", emailStr),
		)
	}

	if len(trimmedEmail) > MaxEmailLength {
		return EmptyEmail, fault.New(
			"email address exceeds maximum length",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input", emailStr),
			fault.WithContext("length", len(trimmedEmail)),
			fault.WithContext("max_length", MaxEmailLength),
		)
	}

	parsedAddr, err := mail.ParseAddress(trimmedEmail)
	if err != nil {
		return EmptyEmail, fault.Wrap(err,
			"email address has an invalid format",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input", emailStr),
		)
	}

	normalizedEmail := strings.ToLower(parsedAddr.Address)

	return Email(normalizedEmail), nil
}

// NewEmail creates a new Email from a string.
// It trims whitespace, validates the format and length, and normalizes the email to lowercase.
// Returns an error if the email is empty, too long, or has an invalid format.
func NewEmail(emailStr string) (Email, error) {
	return parseEmail(emailStr)
}

// MustNewEmail is like NewEmail but panics if the email is invalid.
// This is useful for initializing constant-like email values in tests or initial setup.
func MustNewEmail(emailStr string) Email {
	email, err := NewEmail(emailStr)
	if err != nil {
		panic(err)
	}
	return email
}

// String returns the normalized email address as a string.
func (e Email) String() string {
	return string(e)
}

// IsEmpty returns true if the Email is the zero value.
func (e Email) IsEmpty() bool {
	return e == EmptyEmail
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Email as a JSON string.
func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into an Email, with validation.
func (e *Email) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err,
			"failed to unmarshal email from JSON",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_json", string(data)),
		)
	}

	validatedEmail, err := parseEmail(s)
	if err != nil {
		return err
	}

	*e = validatedEmail
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e Email) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *Email) UnmarshalText(text []byte) error {
	validatedEmail, err := parseEmail(string(text))
	if err != nil {
		return err
	}
	*e = validatedEmail
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the Email as a string, or nil if it is empty.
func (e Email) Value() (driver.Value, error) {
	if e.IsEmpty() {
		return nil, nil
	}
	return e.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a string or byte slice from the database and converts it into an Email, with validation.
func (e *Email) Scan(src interface{}) error {
	if src == nil {
		*e = EmptyEmail
		return nil
	}

	var emailStr string
	switch sval := src.(type) {
	case string:
		emailStr = sval
	case []byte:
		emailStr = string(sval)
	default:
		return fault.New(
			"incompatible type for Email scan",
			fault.WithCode(fault.Invalid),
			fault.WithContext("received_type", fmt.Sprintf("%T", src)),
		)
	}

	validatedEmail, err := parseEmail(emailStr)
	if err != nil {
		return err
	}

	*e = validatedEmail
	return nil
}
