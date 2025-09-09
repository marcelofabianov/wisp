package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/mail"
	"strings"

	"github.com/marcelofabianov/fault"
)

const (
	MaxEmailLength = 254
)

type Email string

var EmptyEmail Email

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

func NewEmail(emailStr string) (Email, error) {
	return parseEmail(emailStr)
}

func MustNewEmail(emailStr string) Email {
	email, err := NewEmail(emailStr)
	if err != nil {
		panic(err)
	}
	return email
}

func (e Email) String() string {
	return string(e)
}

func (e Email) IsEmpty() bool {
	return e == EmptyEmail
}

func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

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

func (e Email) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

func (e *Email) UnmarshalText(text []byte) error {
	validatedEmail, err := parseEmail(string(text))
	if err != nil {
		return err
	}
	*e = validatedEmail
	return nil
}

func (e Email) Value() (driver.Value, error) {
	if e.IsEmpty() {
		return nil, nil
	}
	return e.String(), nil
}

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
