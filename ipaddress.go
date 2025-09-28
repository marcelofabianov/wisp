package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"net"

	"github.com/marcelofabianov/fault"
)

type IPAddress struct {
	ip net.IP
}

var ZeroIPAddress = IPAddress{}

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

func (ip IPAddress) IsZero() bool {
	return ip.ip == nil
}

func (ip IPAddress) IsV4() bool {
	return !ip.IsZero() && ip.ip.To4() != nil
}

func (ip IPAddress) IsV6() bool {
	return !ip.IsZero() && ip.ip.To4() == nil
}

func (ip IPAddress) String() string {
	if ip.IsZero() {
		return ""
	}
	return ip.ip.String()
}

func (ip IPAddress) MarshalJSON() ([]byte, error) {
	if ip.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(ip.String())
}

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

func (ip IPAddress) Value() (driver.Value, error) {
	if ip.IsZero() {
		return nil, nil
	}
	return ip.String(), nil
}

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
