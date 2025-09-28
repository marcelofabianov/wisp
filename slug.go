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

type Slug string

var EmptySlug Slug

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

func (s Slug) String() string {
	return string(s)
}

func (s Slug) IsZero() bool {
	return s == EmptySlug
}

func (s Slug) Equals(other Slug) bool {
	return s == other
}

func (s Slug) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Slug) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "Slug must be a valid JSON string", fault.WithCode(fault.Invalid))
	}

	slug, err := NewSlug(str)
	if err != nil {
		return err
	}
	*s = slug
	return nil
}

func (s Slug) Value() (driver.Value, error) {
	if s.IsZero() {
		return nil, nil
	}
	return s.String(), nil
}

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
