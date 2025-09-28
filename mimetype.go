package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/marcelofabianov/fault"
)

var registeredMIMETypes = make(map[MIMEType]struct{})

type MIMEType string

var EmptyMIMEType MIMEType

func RegisterMIMETypes(types ...string) {
	for _, t := range types {
		normalized := strings.ToLower(strings.TrimSpace(t))
		if normalized != "" {
			parts := strings.Split(normalized, "/")
			if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
				registeredMIMETypes[MIMEType(normalized)] = struct{}{}
			}
		}
	}
}

func ClearRegisteredMIMETypes() {
	registeredMIMETypes = make(map[MIMEType]struct{})
}

func NewMIMEType(input string) (MIMEType, error) {
	normalized := strings.ToLower(strings.TrimSpace(input))

	if normalized == "" {
		return EmptyMIMEType, fault.New("mime type input cannot be empty", fault.WithCode(fault.Invalid))
	}

	parts := strings.Split(normalized, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return EmptyMIMEType, fault.New(
			"mime type must follow the 'type/subtype' format",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", input),
		)
	}

	mt := MIMEType(normalized)
	if !mt.IsRegistered() {
		return EmptyMIMEType, fault.New(
			"mime type is not registered in the allowed list",
			fault.WithCode(fault.Invalid),
			fault.WithContext("mime_type", normalized),
		)
	}

	return mt, nil
}

func (mt MIMEType) IsRegistered() bool {
	_, ok := registeredMIMETypes[mt]
	return ok
}

func (mt MIMEType) Type() string {
	if mt.IsZero() {
		return ""
	}
	parts := strings.Split(string(mt), "/")
	return parts[0]
}

func (mt MIMEType) SubType() string {
	if mt.IsZero() {
		return ""
	}
	parts := strings.Split(string(mt), "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func (mt MIMEType) String() string {
	return string(mt)
}

func (mt MIMEType) IsZero() bool {
	return mt == EmptyMIMEType
}

func (mt MIMEType) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.String())
}

func (mt *MIMEType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "MIMEType must be a valid JSON string", fault.WithCode(fault.Invalid))
	}
	newMT, err := NewMIMEType(s)
	if err != nil {
		return err
	}
	*mt = newMT
	return nil
}

func (mt MIMEType) Value() (driver.Value, error) {
	if mt.IsZero() {
		return nil, nil
	}
	return mt.String(), nil
}

func (mt *MIMEType) Scan(src interface{}) error {
	if src == nil {
		*mt = EmptyMIMEType
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for MIMEType", fault.WithCode(fault.Invalid))
	}

	newMT, err := NewMIMEType(s)
	if err != nil {
		return err
	}
	*mt = newMT
	return nil
}
