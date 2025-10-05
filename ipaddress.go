package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"net"

	"github.com/marcelofabianov/fault"
)

// IPAddress is a value object representing a validated IPv4 or IPv6 address.
// It wraps Go's native `net.IP` to ensure that only valid IP address strings are used.
// This provides type safety for network-related operations.
//
// The zero value is ZeroIPAddress.
//
// Examples:
//   ipv4, err := NewIPAddress("192.168.1.1")
//   ipv6, err := NewIPAddress("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
type IPAddress struct {
	ip net.IP
}

// ZeroIPAddress represents the zero value for the IPAddress type.
var ZeroIPAddress = IPAddress{}

// NewIPAddress creates a new IPAddress from a string.
// It uses `net.ParseIP` to validate the address format.
// Returns an error if the string is not a valid IPv4 or IPv6 address.
func NewIPAddress(value string) (IPAddress, error) {
	if value == "" {
		return ZeroIPAddress, fault.New("ip address cannot be empty", fault.WithCode(fault.Invalid))
	}

	ip := net.ParseIP(value)
	if ip == nil {
		return ZeroIPAddress, fault.New(
			"invalid ip address format",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", value),
		)
	}

	return IPAddress{ip: ip}, nil
}

// IsZero returns true if the IPAddress is the zero value.
func (ip IPAddress) IsZero() bool {
	return ip.ip == nil
}

// IsV4 returns true if the IP address is an IPv4 address.
func (ip IPAddress) IsV4() bool {
	return !ip.IsZero() && ip.ip.To4() != nil
}

// IsV6 returns true if the IP address is an IPv6 address.
func (ip IPAddress) IsV6() bool {
	return !ip.IsZero() && ip.ip.To4() == nil
}

// String returns the string representation of the IP address.
func (ip IPAddress) String() string {
	if ip.IsZero() {
		return ""
	}
	return ip.ip.String()
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the IPAddress to its string representation, or null if zero.
func (ip IPAddress) MarshalJSON() ([]byte, error) {
	if ip.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(ip.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON string into an IPAddress, with validation.
func (ip *IPAddress) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*ip = ZeroIPAddress
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fault.Wrap(err, "IPAddress must be a valid JSON string", fault.WithCode(fault.Invalid))
	}

	newIP, err := NewIPAddress(s)
	if err != nil {
		return err
	}
	*ip = newIP
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the IPAddress as a string.
func (ip IPAddress) Value() (driver.Value, error) {
	if ip.IsZero() {
		return nil, nil
	}
	return ip.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a string or byte slice from the database and converts it into an IPAddress, with validation.
func (ip *IPAddress) Scan(src interface{}) error {
	if src == nil {
		*ip = ZeroIPAddress
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fault.New("unsupported scan type for IPAddress", fault.WithCode(fault.Invalid))
	}

	newIP, err := NewIPAddress(s)
	if err != nil {
		return err
	}
	*ip = newIP
	return nil
}
