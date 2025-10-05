package wisp

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/marcelofabianov/fault"
)

// Preferences is a value object for storing key-value data, such as user settings or configuration.
// It is an immutable wrapper around a `map[string]any`, ensuring that modifications
// do not have side effects. Operations like `Set` return a new `Preferences` instance.
//
// This is useful for managing user-specific settings like theme, language, or notification preferences.
//
// Example:
//   prefs, _ := NewPreferences(map[string]any{"theme": "dark"})
//   newPrefs := prefs.Set("language", "en")
//   theme, _ := newPrefs.Get("theme") // "dark"
type Preferences struct {
	data map[string]any
}

// EmptyPreferences represents the zero value for Preferences (an empty map).
var EmptyPreferences = Preferences{data: make(map[string]any)}

// NewPreferences creates a new Preferences object from a map.
// It creates a defensive copy of the input map to maintain immutability.
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

// ParsePreferences creates a new Preferences object from a JSON byte slice.
// It returns an error if the JSON is invalid.
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

// Get retrieves a value from the preferences by its key.
// The second return value is false if the key does not exist.
func (p Preferences) Get(key string) (any, bool) {
	val, ok := p.data[key]
	return val, ok
}

// Set adds or updates a key-value pair, returning a new Preferences instance.
// This operation is immutable and does not modify the original Preferences object.
func (p Preferences) Set(key string, value any) Preferences {
	newData := make(map[string]any, len(p.data)+1)
	for k, v := range p.data {
		newData[k] = v
	}
	newData[key] = value

	return Preferences{data: newData}
}

// IsZero returns true if the preferences map is empty.
func (p Preferences) IsZero() bool {
	return len(p.data) == 0
}

// Data returns a copy of the underlying data map.
func (p Preferences) Data() map[string]any {
	copyData := make(map[string]any, len(p.data))
	for k, v := range p.data {
		copyData[k] = v
	}
	return copyData
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the preferences map to a JSON object.
func (p Preferences) MarshalJSON() ([]byte, error) {
	if p.IsZero() {
		return []byte("{}"), nil
	}
	return json.Marshal(p.data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object into a Preferences instance.
func (p *Preferences) UnmarshalJSON(data []byte) error {
	prefs, err := ParsePreferences(data)
	if err != nil {
		return err
	}
	*p = prefs
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the Preferences as a JSON byte array.
func (p Preferences) Value() (driver.Value, error) {
	if p.IsZero() {
		return nil, nil
	}
	return p.MarshalJSON()
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a JSON byte array or string from the database and converts it into a Preferences object.
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
