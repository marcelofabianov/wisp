package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/marcelofabianov/fault"
)

// Color is a value object representing a color, parsed from a hex string.
// It stores the color as an `image/color.RGBA` value, ensuring it is always valid.
// It supports both 3-digit (#F00) and 6-digit (#FF0000) hex formats.
//
// The zero value is ZeroColor, which is transparent black.
//
// Example:
//   c, err := ParseColor("#FF0000")
//   r, g, b, a := c.RGBA() // 255, 0, 0, 255
//   hex := c.Hex() // "#ff0000"
type Color struct {
	rgba color.RGBA
}

// ZeroColor represents the zero value for the Color type (transparent black).
var ZeroColor = Color{}

// ParseColor creates a new Color object from a hex string (e.g., "#FF0000" or "#F00").
// It validates the format and returns an error if the hex code is invalid.
func ParseColor(hex string) (Color, error) {
	s := strings.ToLower(strings.TrimSpace(hex))

	if !strings.HasPrefix(s, "#") {
		return ZeroColor, fault.New("hex color must start with '#'", fault.WithCode(fault.Invalid))
	}

	s = strings.TrimPrefix(s, "#")

	var r, g, b uint8
	var a uint8 = 255 // Default alpha is fully opaque

	var err error
	switch len(s) {
	case 3: // #RGB format
		r, err = parseHexComponent(s[0:1] + s[0:1])
		if err == nil {
			g, err = parseHexComponent(s[1:2] + s[1:2])
		}
		if err == nil {
			b, err = parseHexComponent(s[2:3] + s[2:3])
		}
	case 6: // #RRGGBB format
		r, err = parseHexComponent(s[0:2])
		if err == nil {
			g, err = parseHexComponent(s[2:4])
		}
		if err == nil {
			b, err = parseHexComponent(s[4:6])
		}
	default:
		return ZeroColor, fault.New("hex color must have 3 or 6 characters after '#'", fault.WithCode(fault.Invalid))
	}

	if err != nil {
		return ZeroColor, fault.Wrap(err, "invalid hex value in color code", fault.WithCode(fault.Invalid))
	}

	return Color{rgba: color.RGBA{R: r, G: g, B: b, A: a}}, nil
}

// parseHexComponent converts a two-character hex string into a uint8.
func parseHexComponent(s string) (uint8, error) {
	val, err := strconv.ParseUint(s, 16, 8)
	return uint8(val), err
}

// RGBA returns the red, green, blue, and alpha components of the color.
func (c Color) RGBA() (r, g, b, a uint8) {
	return c.rgba.R, c.rgba.G, c.rgba.B, c.rgba.A
}

// Hex returns the 6-digit hex string representation of the color (e.g., "#ff0000").
func (c Color) Hex() string {
	if c.IsZero() {
		return ""
	}
	return fmt.Sprintf("#%02x%02x%02x", c.rgba.R, c.rgba.G, c.rgba.B)
}

// IsZero returns true if the Color is the zero value (transparent black).
func (c Color) IsZero() bool {
	return c.rgba.R == 0 && c.rgba.G == 0 && c.rgba.B == 0 && c.rgba.A == 0
}

// String returns the hex string representation of the color.
func (c Color) String() string {
	return c.Hex()
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Color to its hex string representation.
func (c Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Hex())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string (hex code) into a Color, with validation.
func (c *Color) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "Color must be a valid JSON string", fault.WithCode(fault.Invalid))
	}

	if s == "" {
		*c = ZeroColor
		return nil
	}

	color, err := ParseColor(s)
	if err != nil {
		return err
	}
	*c = color
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the Color as its hex string representation.
func (c Color) Value() (driver.Value, error) {
	if c.IsZero() {
		return nil, nil
	}
	return c.Hex(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a string (hex code) from the database and converts it into a Color.
func (c *Color) Scan(src interface{}) error {
	if src == nil {
		*c = ZeroColor
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for Color", fault.WithCode(fault.Invalid))
	}

	return c.UnmarshalJSON([]byte(`"` + s + `"`))
}
