package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

// ErrValueExceedsMax is a standard error returned when an operation on a BoundedValue
// would cause its current value to exceed its maximum allowed value.
var ErrValueExceedsMax = fault.New("operation would exceed the maximum value", fault.WithCode(fault.Conflict))

// BoundedValue is a value object representing a value that must stay within a `[0, max]` range.
// It is useful for modeling concepts like health points, energy levels, or inventory stock
// where a value cannot exceed a maximum capacity or fall below zero.
//
// All operations are immutable, returning a new BoundedValue instance.
//
// Example:
//   health, _ := NewBoundedValue(80, 100) // 80/100 HP
//   healed, _ := health.Add(15) // 95/100 HP
//   isFull := healed.IsFull() // false
type BoundedValue struct {
	current int64
	max     int64
}

// ZeroBoundedValue represents the zero value for BoundedValue.
var ZeroBoundedValue = BoundedValue{}

// NewBoundedValue creates a new BoundedValue.
// It returns an error if the max value is negative, or if the current value is negative or exceeds the max.
func NewBoundedValue(current, max int64) (BoundedValue, error) {
	if max < 0 {
		return ZeroBoundedValue, fault.New("maximum value cannot be negative", fault.WithCode(fault.Invalid))
	}
	if current < 0 {
		return ZeroBoundedValue, fault.New("current value cannot be negative", fault.WithCode(fault.Invalid))
	}
	if current > max {
		return ZeroBoundedValue, fault.New(
			"current value cannot be greater than the maximum value",
			fault.WithCode(fault.Invalid),
			fault.WithContext("current", current),
			fault.WithContext("max", max),
		)
	}
	return BoundedValue{current: current, max: max}, nil
}

// Current returns the current value.
func (bv BoundedValue) Current() int64 {
	return bv.current
}

// Max returns the maximum allowed value.
func (bv BoundedValue) Max() int64 {
	return bv.max
}

// AvailableSpace returns the difference between the max and current values.
func (bv BoundedValue) AvailableSpace() int64 {
	return bv.max - bv.current
}

// IsFull returns true if the current value is equal to the maximum value.
func (bv BoundedValue) IsFull() bool {
	return bv.current == bv.max
}

// IsZero returns true if the BoundedValue is the zero value.
func (bv BoundedValue) IsZero() bool {
	return bv.current == 0 && bv.max == 0
}

// Add returns a new BoundedValue with the amount added to the current value.
// It returns an error if the amount is negative or if the operation would exceed the max value.
func (bv BoundedValue) Add(amount int64) (BoundedValue, error) {
	if amount < 0 {
		return ZeroBoundedValue, fault.New("amount to add must be non-negative", fault.WithCode(fault.Invalid))
	}
	if bv.AvailableSpace() < amount {
		return ZeroBoundedValue, ErrValueExceedsMax
	}
	return BoundedValue{current: bv.current + amount, max: bv.max}, nil
}

// Subtract returns a new BoundedValue with the amount subtracted from the current value.
// It returns an error if the amount is negative or if the operation would result in a value below zero.
func (bv BoundedValue) Subtract(amount int64) (BoundedValue, error) {
	if amount < 0 {
		return ZeroBoundedValue, fault.New("amount to subtract must be non-negative", fault.WithCode(fault.Invalid))
	}
	if bv.current < amount {
		return ZeroBoundedValue, fault.New("cannot subtract more than the current value (would be negative)", fault.WithCode(fault.Invalid))
	}
	return BoundedValue{current: bv.current - amount, max: bv.max}, nil
}

// Set returns a new BoundedValue with the current value set to a new value.
// It returns an error if the new value is outside the valid [0, max] range.
func (bv BoundedValue) Set(newValue int64) (BoundedValue, error) {
	if newValue < 0 || newValue > bv.max {
		return ZeroBoundedValue, fault.New("new value is outside the allowed [0, max] range", fault.WithCode(fault.Invalid))
	}
	return BoundedValue{current: newValue, max: bv.max}, nil
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the BoundedValue to a JSON object with "current" and "max" fields.
func (bv BoundedValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Current int64 `json:"current"`
		Max     int64 `json:"max"`
	}{
		Current: bv.current,
		Max:     bv.max,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object into a BoundedValue, with validation.
func (bv *BoundedValue) UnmarshalJSON(data []byte) error {
	dto := &struct {
		Current int64 `json:"current"`
		Max     int64 `json:"max"`
	}{}
	if err := json.Unmarshal(data, &dto); err != nil {
		return fault.Wrap(err, "invalid JSON for BoundedValue", fault.WithCode(fault.Invalid))
	}
	newBv, err := NewBoundedValue(dto.Current, dto.Max)
	if err != nil {
		return err
	}
	*bv = newBv
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the BoundedValue as a JSON string or nil if it's the zero value.
func (bv BoundedValue) Value() (driver.Value, error) {
	if bv.IsZero() {
		return nil, nil
	}

	data, err := bv.MarshalJSON()
	if err != nil {
		return nil, fault.Wrap(err,
			"failed to marshal bounded value for database storage",
			fault.WithCode(fault.Internal),
		)
	}

	return string(data), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts string or []byte values containing JSON and validates them as BoundedValue.
func (bv *BoundedValue) Scan(src interface{}) error {
	if src == nil {
		*bv = ZeroBoundedValue
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
			"unsupported scan type for BoundedValue",
			fault.WithCode(fault.Invalid),
			fault.WithContext("received_type", fmt.Sprintf("%T", src)),
		)
	}

	if err := bv.UnmarshalJSON(data); err != nil {
		return err
	}

	return nil
}
