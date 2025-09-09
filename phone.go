package atomic

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/marcelofabianov/fault"
)

type Phone string

var EmptyPhone Phone

var nonDigitRegex = regexp.MustCompile(`\D+`)

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

func NewPhone(input string) (Phone, error) {
	return parsePhone(input)
}

func (p Phone) String() string {
	return string(p)
}

func (p Phone) CountryCode() string {
	if p.IsZero() || len(p) < 2 {
		return ""
	}
	return string(p[0:2])
}

func (p Phone) AreaCode() string {
	if p.IsZero() || len(p) < 4 {
		return ""
	}
	return string(p[2:4])
}

func (p Phone) Number() string {
	if p.IsZero() || len(p) < 4 {
		return ""
	}
	return string(p[4:])
}

func (p Phone) IsZero() bool {
	return p == EmptyPhone
}

func (p Phone) IsMobile() bool {
	return !p.IsZero() && len(p.Number()) == 9
}

func (p Phone) IsLandline() bool {
	return !p.IsZero() && len(p.Number()) == 8
}

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

func (p Phone) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

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

func (p Phone) Value() (driver.Value, error) {
	if p.IsZero() {
		return nil, nil
	}
	return p.String(), nil
}

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
