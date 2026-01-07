package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

// RangedValue is a value object representing a value that must stay within a `[min, max]` range.
// It combines the concepts of MinValue and BoundedValue, enforcing both a lower and an upper bound.
// This is useful for any value that has a defined, inclusive range, such as a rating from 1 to 5.
//
// All operations are immutable, returning a new RangedValue instance.
//
// Example:
//   rating, _ := NewRangedValue(4, 1, 5) // A rating of 4 on a scale of 1-5
//   newRating, _ := rating.Add(1) // 5/5
//   isMax := newRating.IsAtMax() // true
type RangedValue struct {
	current int64
	min     int64
	max     int64
}

// ZeroRangedValue represents the zero value for RangedValue.
var ZeroRangedValue = RangedValue{}

// NewRangedValue creates a new RangedValue.
// It returns an error if min > max, or if the current value is outside the [min, max] range.
func NewRangedValue(current, min, max int64) (RangedValue, error) {
	if min > max {
		return ZeroRangedValue, fault.New("min value cannot be greater than max value", fault.WithCode(fault.Invalid))
	}
	if current < min || current > max {
		return ZeroRangedValue, fault.New(
			"current value is outside the allowed range [min, max]",
			fault.WithCode(fault.Invalid),
			fault.WithContext("current", current),
			fault.WithContext("min", min),
			fault.WithContext("max", max),
		)
	}
	return RangedValue{current: current, min: min, max: max}, nil
}

// Current returns the current value.
func (rv RangedValue) Current() int64 {
	return rv.current
}

// Min returns the minimum allowed value.
func (rv RangedValue) Min() int64 {
	return rv.min
}

// Max returns the maximum allowed value.
func (rv RangedValue) Max() int64 {
	return rv.max
}

// IsAtMin returns true if the current value is equal to the minimum value.
func (rv RangedValue) IsAtMin() bool {
	return rv.current == rv.min
}

// IsAtMax returns true if the current value is equal to the maximum value.
func (rv RangedValue) IsAtMax() bool {
	return rv.current == rv.max
}

// Add returns a new RangedValue with the amount added to the current value.
// It returns an error if the amount is negative or if the operation would exceed the max value.
func (rv RangedValue) Add(amount int64) (RangedValue, error) {
	if amount < 0 {
		return ZeroRangedValue, fault.New("amount to add must be non-negative", fault.WithCode(fault.Invalid))
	}
	if rv.max-rv.current < amount {
		return ZeroRangedValue, ErrValueExceedsMax
	}
	return RangedValue{current: rv.current + amount, min: rv.min, max: rv.max}, nil
}

// Subtract returns a new RangedValue with the amount subtracted from the current value.
// It returns an error if the amount is negative or if the operation would fall below the min value.
func (rv RangedValue) Subtract(amount int64) (RangedValue, error) {
	if amount < 0 {
		return ZeroRangedValue, fault.New("amount to subtract must be non-negative", fault.WithCode(fault.Invalid))
	}
	if rv.current-rv.min < amount {
		return ZeroRangedValue, ErrValueSubceedsMin
	}
	return RangedValue{current: rv.current - amount, min: rv.min, max: rv.max}, nil
}

// Set returns a new RangedValue with the current value set to a new value.
// It returns an error if the new value is outside the allowed [min, max] range.
func (rv RangedValue) Set(newValue int64) (RangedValue, error) {
	if newValue < rv.min || newValue > rv.max {
		return ZeroRangedValue, fault.New("new value is outside the allowed range", fault.WithCode(fault.Invalid))
	}
	return RangedValue{current: newValue, min: rv.min, max: rv.max}, nil
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the RangedValue to a JSON object with "current", "min", and "max" fields.
func (rv RangedValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Current int64 `json:"current"`
		Min     int64 `json:"min"`
		Max     int64 `json:"max"`
	}{
		Current: rv.current,
		Min:     rv.min,
		Max:     rv.max,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object into a RangedValue, with validation.
func (rv *RangedValue) UnmarshalJSON(data []byte) error {
	dto := &struct {
		Current int64 `json:"current"`
		Min     int64 `json:"min"`
		Max     int64 `json:"max"`
	}{}
	if err := json.Unmarshal(data, &dto); err != nil {
		return fault.Wrap(err, "invalid JSON for RangedValue", fault.WithCode(fault.Invalid))
	}
	newRv, err := NewRangedValue(dto.Current, dto.Min, dto.Max)
	if err != nil {
		return err
	}
	*rv = newRv
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the RangedValue as a JSON string or nil if it's the zero value.
func (rv RangedValue) Value() (driver.Value, error) {
	if rv.current == 0 && rv.min == 0 && rv.max == 0 {
		return nil, nil
	}

	data, err := rv.MarshalJSON()
	if err != nil {
		return nil, fault.Wrap(err,
			"failed to marshal ranged value for database storage",
			fault.WithCode(fault.Internal),
		)
	}

	return string(data), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts string or []byte values containing JSON and validates them as RangedValue.
func (rv *RangedValue) Scan(src interface{}) error {
	if src == nil {
		*rv = ZeroRangedValue
		return nil
	}

	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fault.New(
			"unsupported scan type for RangedValue",
			fault.WithCode(fault.Invalid),
			fault.WithContext("received_type", fmt.Sprintf("%T", src)),
		)
	}

	if err := rv.UnmarshalJSON(data); err != nil {
		return err
	}

	return nil
}
