package wisp

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/marcelofabianov/fault"
)

type Preferences struct {
	data map[string]any
}

var EmptyPreferences = Preferences{data: make(map[string]any)}

func NewPreferences(data map[string]any) (Preferences, error) {
	if data == nil {
		return EmptyPreferences, nil
	}

	newData := make(map[string]any, len(data))
	for k, v := range data {
		newData[k] = v
	}

	return Preferences{data: newData}, nil
}

func ParsePreferences(jsonData []byte) (Preferences, error) {
	if len(jsonData) == 0 || string(jsonData) == "null" {
		return EmptyPreferences, nil
	}

	var data map[string]any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return EmptyPreferences, fault.Wrap(err, "invalid JSON format for Preferences", fault.WithCode(fault.Invalid))
	}

	return NewPreferences(data)
}

func (p Preferences) Get(key string) (any, bool) {
	val, ok := p.data[key]
	return val, ok
}

func (p Preferences) Set(key string, value any) Preferences {
	newData := make(map[string]any, len(p.data)+1)
	for k, v := range p.data {
		newData[k] = v
	}
	newData[key] = value

	return Preferences{data: newData}
}

func (p Preferences) IsZero() bool {
	return len(p.data) == 0
}

func (p Preferences) Data() map[string]any {
	copyData := make(map[string]any, len(p.data))
	for k, v := range p.data {
		copyData[k] = v
	}
	return copyData
}

func (p Preferences) MarshalJSON() ([]byte, error) {
	if p.IsZero() {
		return []byte("{}"), nil
	}
	return json.Marshal(p.data)
}

func (p *Preferences) UnmarshalJSON(data []byte) error {
	prefs, err := ParsePreferences(data)
	if err != nil {
		return err
	}
	*p = prefs
	return nil
}

func (p Preferences) Value() (driver.Value, error) {
	if p.IsZero() {
		return nil, nil
	}
	return p.MarshalJSON()
}

func (p *Preferences) Scan(src interface{}) error {
	if src == nil {
		*p = EmptyPreferences
		return nil
	}

	var bytes []byte
	switch v := src.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fault.New("unsupported scan type for Preferences", fault.WithCode(fault.Invalid))
	}

	return p.UnmarshalJSON(bytes)
}
