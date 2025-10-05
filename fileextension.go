package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/marcelofabianov/fault"
)

// validExtensionRegex defines the allowed characters in a file extension (alphanumeric).
var validExtensionRegex = regexp.MustCompile(`^[a-z0-9]+$`)

// registeredExtensions holds the global set of allowed file extensions.
var registeredExtensions = make(map[FileExtension]struct{})

// FileExtension is a value object representing a file extension (e.g., "jpg", "pdf").
// It ensures that only explicitly registered extensions are used, preventing the use of arbitrary
// or unsafe file types. Extensions are stored in a normalized (lowercase, no dot) format.
//
// Before use, extensions must be added to a global registry via `RegisterFileExtensions`.
//
// Example:
//   wisp.RegisterFileExtensions(".JPG", "pdf")
//   ext, err := wisp.NewFileExtension("jpg")
//   fmt.Println(ext.StringWithDot()) // ".jpg"
type FileExtension string

// EmptyFileExtension represents the zero value for the FileExtension type.
var EmptyFileExtension FileExtension

// RegisterFileExtensions adds one or more file extensions to the global registry.
// It normalizes them to lowercase and removes any leading dot.
// This function should be called at application startup to define the allowed file types.
func RegisterFileExtensions(extensions ...string) {
	for _, extStr := range extensions {
		normalized := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(extStr), "."))
		if normalized != "" && validExtensionRegex.MatchString(normalized) {
			registeredExtensions[FileExtension(normalized)] = struct{}{}
		}
	}
}

// ClearRegisteredFileExtensions removes all extensions from the global registry.
// This is primarily for testing purposes to ensure a clean state.
func ClearRegisteredFileExtensions() {
	registeredExtensions = make(map[FileExtension]struct{})
}

// NewFileExtension creates a new FileExtension from a string.
// It normalizes the input (lowercase, no dot) and validates it against the global registry.
// Returns an error if the extension is empty or not registered.
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

// IsRegistered checks if the file extension is in the global registry.
func (fe FileExtension) IsRegistered() bool {
	_, ok := registeredExtensions[fe]
	return ok
}

// String returns the normalized file extension without a leading dot.
func (fe FileExtension) String() string {
	return string(fe)
}

// StringWithDot returns the file extension with a leading dot (e.g., ".jpg").
func (fe FileExtension) StringWithDot() string {
	if fe.IsZero() {
		return ""
	}
	return "." + string(fe)
}

// IsZero returns true if the FileExtension is the zero value.
func (fe FileExtension) IsZero() bool {
	return fe == EmptyFileExtension
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the FileExtension to its string representation (without the dot).
func (fe FileExtension) MarshalJSON() ([]byte, error) {
	return json.Marshal(fe.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into a FileExtension, with validation against the registry.
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

// Value implements the driver.Valuer interface for database storage.
// It returns the FileExtension as a string.
func (fe FileExtension) Value() (driver.Value, error) {
	if fe.IsZero() {
		return nil, nil
	}
	return fe.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a string or byte slice from the database and converts it into a FileExtension, with validation.
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
