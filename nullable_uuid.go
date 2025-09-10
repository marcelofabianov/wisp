package wisp

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/marcelofabianov/fault"
)

type NullableUUID struct {
	UUID  UUID
	Valid bool
}

func NewNullableUUID(id UUID) NullableUUID {
	return NullableUUID{
		UUID:  id,
		Valid: !id.IsNil(),
	}
}

func (nu NullableUUID) IsZero() bool {
	return !nu.Valid
}

func (nu NullableUUID) MarshalJSON() ([]byte, error) {
	if !nu.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nu.UUID)
}

func (nu *NullableUUID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nu.Valid = false
		return nil
	}

	var u UUID
	if err := json.Unmarshal(data, &u); err != nil {
		return fault.Wrap(err, "NullableUUID must be a valid JSON UUID string or null", fault.WithCode(fault.Invalid))
	}

	nu.UUID = u
	nu.Valid = true
	return nil
}

func (nu NullableUUID) Value() (driver.Value, error) {
	if !nu.Valid {
		return nil, nil
	}
	return nu.UUID.Value()
}

func (nu *NullableUUID) Scan(src interface{}) error {
	if src == nil {
		nu.UUID, nu.Valid = Nil, false
		return nil
	}

	var u UUID
	if err := u.Scan(src); err != nil {
		return err
	}

	nu.UUID = u
	nu.Valid = true
	return nil
}
