package wisp

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/marcelofabianov/fault"
)

type UUID uuid.UUID

var Nil UUID

func NewUUID() (UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Nil, fault.Wrap(err,
			"failed to generate v7 UUID",
			fault.WithCode(fault.Internal),
			fault.WithContext("operation", "wisp.NewUUID"),
		)
	}
	return UUID(id), nil
}

func MustNewUUID() UUID {
	id, err := NewUUID()
	if err != nil {
		panic(err)
	}
	return id
}

func ParseUUID(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return Nil, fault.Wrap(err,
			"failed to parse UUID string",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input", s),
			fault.WithContext("operation", "wisp.ParseUUID"),
		)
	}
	return UUID(id), nil
}

func MustParseUUID(s string) UUID {
	id, err := ParseUUID(s)
	if err != nil {
		panic(err)
	}
	return id
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

func (u UUID) IsNil() bool {
	return u == Nil
}

func (u UUID) MarshalText() ([]byte, error) {
	return uuid.UUID(u).MarshalText()
}

func (u *UUID) UnmarshalText(text []byte) error {
	var underlyingUUID uuid.UUID
	if err := underlyingUUID.UnmarshalText(text); err != nil {
		return fault.Wrap(err,
			"invalid text representation for UUID",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_text", string(text)),
			fault.WithContext("operation", "wisp.UUID.UnmarshalText"),
		)
	}
	*u = UUID(underlyingUUID)
	return nil
}

func (u UUID) Value() (driver.Value, error) {
	if u == Nil {
		return nil, nil
	}
	return u.String(), nil
}

func (u *UUID) Scan(src interface{}) error {
	var underlyingUUID uuid.UUID
	if err := underlyingUUID.Scan(src); err != nil {
		return fault.Wrap(err,
			"failed to scan database value into UUID",
			fault.WithCode(fault.Invalid),
			fault.WithContext("source_type", fmt.Sprintf("%T", src)),
			fault.WithContext("operation", "wisp.UUID.Scan"),
		)
	}
	*u = UUID(underlyingUUID)
	return nil
}
