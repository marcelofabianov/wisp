package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/marcelofabianov/fault"
)

// Phone is a value object representing a Brazilian phone number.
// It stores the number in a normalized E.164-like format (e.g., "5511987654321"),
// including the country code (55 for Brazil), area code (DDD), and the local number.
//
// The type validates the DDD, the number of digits for mobile vs. landline, and the mobile prefix.
//
// Examples:
//   - Input: "(11) 98765-4321"
//   - Stored as: "5511987654321"
//   - Formatted output: "+55 (11) 98765-4321"
type Phone string

// EmptyPhone represents the zero value for the Phone type.
var EmptyPhone Phone

// nonDigitRegex is used to remove all non-numeric characters from a phone number string.
var nonDigitRegex = regexp.MustCompile(`\\D+`)

// validDDDs is the set of all valid Brazilian area codes (DDD).
var validDDDs = map[string]struct{}{
	"11": {}, "12": {}, "13": {}, "14": {}, "15": {}, "16": {}, "17": {}, "18": {}, "19": {},
	"21": {}, "22": {}, "24": {}, "27": {}, "28": {},
	"31": {}, "32": {}, "33": {}, "34": {}, "35": {}, "37": {}, "38": {},
	"41": {}, "42": {}, "43": {}, "44": {}, "45": {}, "46": {}, "47": {}, "48": {}, "49": {},
	"51": {}, "53": {}, "54": {}, "55": {},
	"61": {}, "62": {}, "63": {}, "64": {}, "65": {}, "66": {}, "67": {}, "68": {}, "69": {},
	"71": {}, "73": {}, "74": {}, "75": {}, "77": {}, "79": {},
	"81": {}, "82": {}, "83": {}, "84": {}, "85": {}, "86": {}, "87": {}, "88": {}, "89": {},
	"91": {}, "92": {}, "93": {}, "94": {}, "95": {}, "96": {}, "97": {}, "98": {}, "99": {},
}

// parsePhone contains the core logic for validating and normalizing a Brazilian phone number.
func parsePhone(input string) (Phone, error) {
	if input == "" {
		return EmptyPhone, nil
	}

	sanitized := nonDigitRegex.ReplaceAllString(input, "")

	if len(sanitized) < 10 {
		return EmptyPhone, fault.New("phone number is too short", fault.WithCode(fault.Invalid), fault.WithContext("input", input))
	}

	if len(sanitized) > 13 {
		return EmptyPhone, fault.New("phone number is too long", fault.WithCode(fault.Invalid), fault.WithContext("input", input))
	}

	if !strings.HasPrefix(sanitized, "55") {
		sanitized = "55" + sanitized
	}

	if len(sanitized) != 12 && len(sanitized) != 13 {
		return EmptyPhone, fault.New("invalid phone number length after normalization", fault.WithCode(fault.Invalid), fault.WithContext("normalized_number", sanitized))
	}

	areaCode := sanitized[2:4]
	if _, ok := validDDDs[areaCode]; !ok {
		return EmptyPhone, fault.New("invalid area code (DDD)", fault.WithCode(fault.Invalid), fault.WithContext("area_code", areaCode))
	}

	numberPart := sanitized[4:]
	if len(numberPart) == 9 && numberPart[0] != '9' {
		return EmptyPhone, fault.New("mobile number must start with digit 9", fault.WithCode(fault.Invalid), fault.WithContext("number", numberPart))
	}
	if len(numberPart) == 8 && (numberPart[0] < '2' || numberPart[0] > '5') {
		return EmptyPhone, fault.New("landline number has an invalid prefix", fault.WithCode(fault.Invalid), fault.WithContext("number", numberPart))
	}

	return Phone(sanitized), nil
}

// NewPhone creates a new Phone from a string.
// It sanitizes the input by removing non-digit characters, validates the length,
// ensures the Brazilian country code (55) is present, and validates the area code (DDD) and number format.
// It returns an error if the phone number is invalid in any of these ways.
func NewPhone(input string) (Phone, error) {
	return parsePhone(input)
}

// String returns the normalized phone number as a string (e.g., "5511987654321").
func (p Phone) String() string {
	return string(p)
}

// CountryCode returns the country code part of the number.
func (p Phone) CountryCode() string {
	if p.IsZero() || len(p) < 2 {
		return ""
	}
	return string(p[0:2])
}

// AreaCode returns the area code (DDD) part of the number.
func (p Phone) AreaCode() string {
	if p.IsZero() || len(p) < 4 {
		return ""
	}
	return string(p[2:4])
}

// Number returns the local number part (without country or area code).
func (p Phone) Number() string {
	if p.IsZero() || len(p) < 4 {
		return ""
	}
	return string(p[4:])
}

// IsZero returns true if the Phone is the zero value.
func (p Phone) IsZero() bool {
	return p == EmptyPhone
}

// IsMobile returns true if the phone number is identified as a mobile number (9 digits).
func (p Phone) IsMobile() bool {
	return !p.IsZero() && len(p.Number()) == 9
}

// IsLandline returns true if the phone number is identified as a landline number (8 digits).
func (p Phone) IsLandline() bool {
	return !p.IsZero() && len(p.Number()) == 8
}

// Formatted returns the phone number in a human-readable format.
// Mobile: "+55 (11) 98765-4321"
// Landline: "+55 (11) 4321-5432"
func (p Phone) Formatted() string {
	if p.IsZero() {
		return ""
	}
	number := p.Number()
	if p.IsMobile() {
		return fmt.Sprintf("+%s (%s) %s-%s", p.CountryCode(), p.AreaCode(), number[:5], number[5:])
	}
	return fmt.Sprintf("+%s (%s) %s-%s", p.CountryCode(), p.AreaCode(), number[:4], number[4:])
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Phone to its normalized string representation.
func (p Phone) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into a Phone, with validation.
func (p *Phone) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "phone must be a valid JSON string", fault.WithCode(fault.Invalid))
	}
	phone, err := NewPhone(s)
	if err != nil {
		return err
	}
	*p = phone
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the normalized phone number as a string.
func (p Phone) Value() (driver.Value, error) {
	if p.IsZero() {
		return nil, nil
	}
	return p.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a string or byte slice from the database and converts it into a Phone, with validation.
func (p *Phone) Scan(src interface{}) error {
	if src == nil {
		*p = EmptyPhone
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for Phone", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}

	phone, err := NewPhone(s)
	if err != nil {
		return err
	}
	*p = phone
	return nil
}
