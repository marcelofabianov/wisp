package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

// ErrValueSubceedsMin is a standard error returned when an operation on a MinValue
// would cause its current value to fall below its minimum allowed value.
var ErrValueSubceedsMin = fault.New("operation would fall below the minimum value", fault.WithCode(fault.Conflict))

// MinValue is a value object representing a value that must stay at or above a specified minimum.
// It is useful for concepts like a minimum order quantity, a score that cannot go below a certain floor,
// or any value with a lower bound.
//
// All operations are immutable, returning a new MinValue instance.
//
// Example:
//   stock, _ := NewMinValue(10, 5) // 10 items, with a minimum stock of 5
//   sold, _ := stock.Subtract(3) // 7 items
//   isLow := sold.IsAtMin() // false
type MinValue struct {
	current int64
	min     int64
}

// ZeroMinValue represents the zero value for MinValue.
var ZeroMinValue = MinValue{}

// NewMinValue creates a new MinValue.
// It returns an error if the current value is less than the specified minimum.
func NewMinValue(current, min int64) (MinValue, error) {
	if current < min {
		return ZeroMinValue, fault.New(
			"current value cannot be less than the minimum value",
			fault.WithCode(fault.Invalid),
			fault.WithContext("current", current),
			fault.WithContext("min", min),
		)
	}
	return MinValue{current: current, min: min}, nil
}

// Current returns the current value.
func (mv MinValue) Current() int64 {
	return mv.current
}

// Min returns the minimum allowed value.
func (mv MinValue) Min() int64 {
	return mv.min
}

// IsAtMin returns true if the current value is equal to the minimum value.
func (mv MinValue) IsAtMin() bool {
	return mv.current == mv.min
}

// Add returns a new MinValue with the amount added to the current value.
// It returns an error if the amount to add is negative.
func (mv MinValue) Add(amount int64) (MinValue, error) {
	if amount < 0 {
		return ZeroMinValue, fault.New("amount to add must be non-negative", fault.WithCode(fault.Invalid))
	}
	return MinValue{current: mv.current + amount, min: mv.min}, nil
}

// Subtract returns a new MinValue with the amount subtracted from the current value.
// It returns an error if the amount is negative or if the operation would fall below the minimum value.
func (mv MinValue) Subtract(amount int64) (MinValue, error) {
	if amount < 0 {
		return ZeroMinValue, fault.New("amount to subtract must be non-negative", fault.WithCode(fault.Invalid))
	}
	if mv.current-mv.min < amount {
		return ZeroMinValue, ErrValueSubceedsMin
	}
	return MinValue{current: mv.current - amount, min: mv.min}, nil
}

// Set returns a new MinValue with the current value set to a new value.
// It returns an error if the new value is less than the minimum.
func (mv MinValue) Set(newValue int64) (MinValue, error) {
	if newValue < mv.min {
		return ZeroMinValue, fault.New("new value cannot be less than the minimum", fault.WithCode(fault.Invalid))
	}
	return MinValue{current: newValue, min: mv.min}, nil
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the MinValue to a JSON object with "current" and "min" fields.
func (mv MinValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Current int64 `json:"current"`
		Min     int64 `json:"min"`
	}{
		Current: mv.current,
		Min:     mv.min,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object into a MinValue, with validation.
func (mv *MinValue) UnmarshalJSON(data []byte) error {
	dto := &struct {
		Current int64 `json:"current"`
		Min     int64 `json:"min"`
	}{}
	if err := json.Unmarshal(data, &dto); err != nil {
		return fault.Wrap(err, "invalid JSON for MinValue", fault.WithCode(fault.Invalid))
	}
	newMv, err := NewMinValue(dto.Current, dto.Min)
	if err != nil {
		return err
	}
	*mv = newMv
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the MinValue as a JSON string or nil if it's the zero value.
func (mv MinValue) Value() (driver.Value, error) {
	if mv.current == 0 && mv.min == 0 {
		return nil, nil
	}

	data, err := mv.MarshalJSON()
	if err != nil {
		return nil, fault.Wrap(err,
			"failed to marshal min value for database storage",
			fault.WithCode(fault.Internal),
		)
	}

	return string(data), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts string or []byte values containing JSON and validates them as MinValue.
func (mv *MinValue) Scan(src interface{}) error {
	if src == nil {
		*mv = ZeroMinValue
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
			"unsupported scan type for MinValue",
			fault.WithCode(fault.Invalid),
			fault.WithContext("received_type", fmt.Sprintf("%T", src)),
		)
	}

	if err := mv.UnmarshalJSON(data); err != nil {
		return err
	}

	return nil
}
