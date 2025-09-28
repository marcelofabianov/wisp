package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/marcelofabianov/fault"
)

var validExtensionRegex = regexp.MustCompile(`^[a-z0-9]+$`)
var registeredExtensions = make(map[FileExtension]struct{})

type FileExtension string

var EmptyFileExtension FileExtension

func RegisterFileExtensions(extensions ...string) {
	for _, extStr := range extensions {
		normalized := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(extStr), "."))
		if normalized != "" && validExtensionRegex.MatchString(normalized) {
			registeredExtensions[FileExtension(normalized)] = struct{}{}
		}
	}
}

func ClearRegisteredFileExtensions() {
	registeredExtensions = make(map[FileExtension]struct{})
}

func NewFileExtension(input string) (FileExtension, error) {
	trimmed := strings.TrimSpace(input)
	normalized := strings.ToLower(strings.TrimPrefix(trimmed, "."))

	if normalized == "" {
		return EmptyFileExtension, fault.New("file extension cannot be empty", fault.WithCode(fault.Invalid))
	}

	ext := FileExtension(normalized)
	if !ext.IsRegistered() {
		return EmptyFileExtension, fault.New(
			"file extension is not registered in the allowed list",
			fault.WithCode(fault.Invalid),
			fault.WithContext("extension", normalized),
		)
	}

	return ext, nil
}

func (fe FileExtension) IsRegistered() bool {
	_, ok := registeredExtensions[fe]
	return ok
}

func (fe FileExtension) String() string {
	return string(fe)
}

func (fe FileExtension) StringWithDot() string {
	if fe.IsZero() {
		return ""
	}
	return "." + string(fe)
}

func (fe FileExtension) IsZero() bool {
	return fe == EmptyFileExtension
}

func (fe FileExtension) MarshalJSON() ([]byte, error) {
	return json.Marshal(fe.String())
}

func (fe *FileExtension) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "FileExtension must be a valid JSON string", fault.WithCode(fault.Invalid))
	}
	ext, err := NewFileExtension(s)
	if err != nil {
		return err
	}
	*fe = ext
	return nil
}

func (fe FileExtension) Value() (driver.Value, error) {
	if fe.IsZero() {
		return nil, nil
	}
	return fe.String(), nil
}

func (fe *FileExtension) Scan(src interface{}) error {
	if src == nil {
		*fe = EmptyFileExtension
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for FileExtension", fault.WithCode(fault.Invalid))
	}

	ext, err := NewFileExtension(s)
	if err != nil {
		return err
	}
	*fe = ext
	return nil
}
