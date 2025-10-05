package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/marcelofabianov/fault"
)

// UF represents a Brazilian state code (Unidade Federativa).
// It is a value object that ensures the code is a valid, two-letter, uppercase abbreviation
// corresponding to one of the Brazilian states or the Federal District.
//
// Examples:
//   - Input: "sp" or " SP "
//   - Stored as: "SP"
type UF string

// EmptyUF represents the zero value for the UF type.
var EmptyUF UF

// validUFs holds the set of all valid Brazilian state codes.
var validUFs = map[UF]struct{}{
	"AC": {}, "AL": {}, "AP": {}, "AM": {}, "BA": {}, "CE": {}, "DF": {}, "ES": {}, "GO": {},
	"MA": {}, "MT": {}, "MS": {}, "MG": {}, "PA": {}, "PB": {}, "PR": {}, "PE": {}, "PI": {},
	"RJ": {}, "RN": {}, "RS": {}, "RO": {}, "RR": {}, "SC": {}, "SP": {}, "SE": {}, "TO": {},
}

// NewUF creates a new UF from a string.
// It normalizes the input to uppercase and validates it against the list of official Brazilian state codes.
// Returns an error if the code is not a valid UF.
func NewUF(input string) (UF, error) {
	uf := UF(strings.ToUpper(strings.TrimSpace(input)))

	if uf.IsZero() {
		return EmptyUF, nil
	}

	if !uf.IsValid() {
		return EmptyUF, fault.New(
			"invalid UF code",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_code", input),
		)
	}
	return uf, nil
}

// String returns the UF code as a string.
func (u UF) String() string {
	return string(u)
}

// IsValid checks if the UF is in the list of official Brazilian state codes.
func (u UF) IsValid() bool {
	_, ok := validUFs[u]
	return ok
}

// IsZero returns true if the UF is the zero value.
func (u UF) IsZero() bool {
	return u == EmptyUF
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the UF to its string representation.
func (u UF) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into a UF, with validation.
func (u *UF) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*u = EmptyUF
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "UF must be a valid JSON string", fault.WithCode(fault.Invalid))
	}

	uf, err := NewUF(s)
	if err != nil {
		return err
	}
	*u = uf
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the UF as a string.
func (u UF) Value() (driver.Value, error) {
	if u.IsZero() {
		return nil, nil
	}
	return u.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a string or byte slice from the database and converts it into a UF, with validation.
func (u *UF) Scan(src interface{}) error {
	if src == nil {
		*u = EmptyUF
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for UF", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}

	uf, err := NewUF(s)
	if err != nil {
		return err
	}
	*u = uf
	return nil
}
