package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/marcelofabianov/fault"
)

type Version int

var ZeroVersion Version

func NewVersion(v int) (Version, error) {
	if v < 0 {
		return ZeroVersion, fault.New(
			"version cannot be negative",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", v),
		)
	}
	return Version(v), nil
}

func InitialVersion() Version {
	return Version(1)
}

func (v Version) Increment() Version {
	return v + 1
}

func (v Version) Previous() Version {
	if v <= 1 {
		return ZeroVersion
	}
	return v - 1
}

func (v Version) IsZero() bool {
	return v == ZeroVersion
}

func (v Version) Equals(other Version) bool {
	return v == other
}

func (v Version) IsGreaterThan(other Version) bool {
	return v > other
}

func (v Version) IsLessThan(other Version) bool {
	return v < other
}

func (v Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(v))
}

func (v *Version) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*v = ZeroVersion
		return nil
	}

	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return fault.Wrap(err,
			"version must be a valid JSON number",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_json", string(data)),
		)
	}

	if i < 0 {
		return fault.New(
			"version cannot be negative",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", i),
		)
	}

	*v = Version(i)
	return nil
}

func (v Version) Value() (driver.Value, error) {
	return int64(v), nil
}

func (v *Version) Scan(src interface{}) error {
	if src == nil {
		*v = ZeroVersion
		return nil
	}

	var intVal int64
	switch s := src.(type) {
	case int64:
		intVal = s
	case []byte:
		parsed, err := strconv.ParseInt(string(s), 10, 64)
		if err != nil {
			return fault.Wrap(err,
				"failed to convert bytes to version number",
				fault.WithCode(fault.Invalid),
				fault.WithContext("input_bytes", string(s)),
			)
		}
		intVal = parsed
	default:
		return fault.New(
			"incompatible type for Version scan",
			fault.WithCode(fault.Invalid),
			fault.WithContext("received_type", fmt.Sprintf("%T", src)),
		)
	}

	if intVal < 0 {
		return fault.New(
			"version from database cannot be negative",
			fault.WithCode(fault.Invalid),
			fault.WithContext("source_value", intVal),
		)
	}

	*v = Version(intVal)
	return nil
}

func (v Version) Int() int {
	return int(v)
}
