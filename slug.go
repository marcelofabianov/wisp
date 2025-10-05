package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"regexp"
	"strings"
	"unicode"

	"github.com/marcelofabianov/fault"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	invalidSlugCharsRegex = regexp.MustCompile(`[^a-z0-9-]`)
	multipleHyphensRegex  = regexp.MustCompile(`-{2,}`)
	symbolReplacer        = strings.NewReplacer(
		"%", " percent ",
		"&", " and ",
		"@", " at ",
		"$", " dollar ",
		"€", " euro ",
		"£", " pound ",
		"+", " plus ",
	)
)

// Slug represents a URL-friendly string identifier that follows web standards.
// It automatically normalizes input text to create safe, consistent URLs.
//
// A slug is created by:
//   - Removing diacritics (é -> e, ñ -> n)
//   - Converting to lowercase
//   - Replacing symbols with descriptive words (& -> and, % -> percent)
//   - Replacing non-alphanumeric characters with hyphens
//   - Collapsing multiple hyphens into single hyphens
//   - Trimming leading and trailing hyphens
//
// Examples:
//   - Input: "Hello, World!" -> Output: "hello-world"
//   - Input: "Café & Bar" -> Output: "cafe-and-bar"
//   - Input: "100% Natural" -> Output: "100-percent-natural"
//
// Slugs are commonly used for:
//   - URL paths: /blog/my-article-title
//   - File names: my-document-name.pdf
//   - API endpoints: /users/john-doe
type Slug string

// EmptySlug represents the zero value for Slug type.
var EmptySlug Slug

// NewSlug creates a new Slug from the given input string.
// It performs normalization to ensure the result is URL-safe and consistent.
//
// The normalization process:
//   1. Removes diacritics and accents (Unicode normalization)
//   2. Converts to lowercase
//   3. Replaces common symbols with descriptive words
//   4. Replaces non-alphanumeric characters with hyphens
//   5. Collapses multiple consecutive hyphens
//   6. Trims leading and trailing hyphens
//
// Returns an error if the input results in an empty string after normalization.
//
// Examples:
//   slug, err := NewSlug("Hello, World!")     // Returns: "hello-world"
//   slug, err := NewSlug("Café & Restaurant") // Returns: "cafe-and-restaurant"
//   slug, err := NewSlug("---")              // Returns: error (empty after normalization)
func NewSlug(input string) (Slug, error) {
	// 1. Transliterate to remove diacritics (e.g., "é" -> "e")
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalized, _, err := transform.String(t, input)
	if err != nil {
		return EmptySlug, fault.Wrap(err, "failed to normalize string for slug", fault.WithCode(fault.Internal))
	}

	// 2. Convert to lowercase
	normalized = strings.ToLower(normalized)

	// 3. Replace common symbols with words
	normalized = symbolReplacer.Replace(normalized)

	// 4. Replace all other non-alphanumeric characters with a hyphen
	normalized = invalidSlugCharsRegex.ReplaceAllString(normalized, "-")

	// 5. Collapse multiple hyphens into a single one
	normalized = multipleHyphensRegex.ReplaceAllString(normalized, "-")

	// 6. Trim leading and trailing hyphens
	normalized = strings.Trim(normalized, "-")

	if normalized == "" {
		return EmptySlug, fault.New(
			"slug input cannot be empty or result in an empty string after normalization",
			fault.WithCode(fault.Invalid),
			fault.WithContext("original_input", input),
		)
	}

	return Slug(normalized), nil
}

// String returns the slug as a string.
// The returned value is guaranteed to be URL-safe and normalized.
func (s Slug) String() string {
	return string(s)
}

// IsZero returns true if the slug is the zero value (EmptySlug).
func (s Slug) IsZero() bool {
	return s == EmptySlug
}

// Equals returns true if this slug is equal to the other slug.
// Since slugs are normalized, this performs a simple string comparison.
func (s Slug) Equals(other Slug) bool {
	return s == other
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the slug as a JSON string.
func (s Slug) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into a Slug, performing full normalization.
func (s *Slug) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fault.Wrap(err, "Slug must be a valid JSON string", fault.WithCode(fault.Invalid))
	}

	slug, err := NewSlug(str)
	if err != nil {
		return err
	}
	*s = slug
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the slug as a string or nil if zero value.
func (s Slug) Value() (driver.Value, error) {
	if s.IsZero() {
		return nil, nil
	}
	return s.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts string or []byte values and stores them as-is (assuming they're already normalized).
// For proper validation, use NewSlug when creating slugs from user input.
func (s *Slug) Scan(src interface{}) error {
	if src == nil {
		*s = EmptySlug
		return nil
	}

	var str string
	switch v := src.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fault.New("unsupported scan type for Slug", fault.WithCode(fault.Invalid))
	}

	if str == "" {
		*s = EmptySlug
		return nil
	}

	*s = Slug(str)
	return nil
}
