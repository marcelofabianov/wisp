package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/marcelofabianov/fault"
)

type AuditUser string

var EmptyAuditUser AuditUser

const SystemAuditUser AuditUser = "system"

func NewSystemAuditUser() AuditUser {
	return SystemAuditUser
}

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

func (au AuditUser) String() string {
	return string(au)
}

func (au AuditUser) IsZero() bool {
	return au == EmptyAuditUser
}

func (au AuditUser) IsSystem() bool {
	return au == SystemAuditUser
}

func (au AuditUser) IsEmail() bool {
	return !au.IsZero() && !au.IsSystem()
}

func (au AuditUser) Email() (Email, bool) {
	if !au.IsEmail() {
		return EmptyEmail, false
	}

	email, _ := NewEmail(string(au))
	return email, true
}

func (au AuditUser) MarshalJSON() ([]byte, error) {
	return json.Marshal(au.String())
}

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

func (au AuditUser) Value() (driver.Value, error) {
	if au.IsZero() {
		return nil, nil
	}
	return au.String(), nil
}

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
