package atomic

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/marcelofabianov/fault"
)

type UF string

var EmptyUF UF

var validUFs = map[UF]struct{}{
	"AC": {}, "AL": {}, "AP": {}, "AM": {}, "BA": {}, "CE": {}, "DF": {}, "ES": {}, "GO": {},
	"MA": {}, "MT": {}, "MS": {}, "MG": {}, "PA": {}, "PB": {}, "PR": {}, "PE": {}, "PI": {},
	"RJ": {}, "RN": {}, "RS": {}, "RO": {}, "RR": {}, "SC": {}, "SP": {}, "SE": {}, "TO": {},
}

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

func (u UF) String() string {
	return string(u)
}

func (u UF) IsValid() bool {
	_, ok := validUFs[u]
	return ok
}

func (u UF) IsZero() bool {
	return u == EmptyUF
}

func (u UF) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

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

func (u UF) Value() (driver.Value, error) {
	if u.IsZero() {
		return nil, nil
	}
	return u.String(), nil
}

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
