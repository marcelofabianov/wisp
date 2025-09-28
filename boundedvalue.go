package wisp

import (
	"encoding/json"

	"github.com/marcelofabianov/fault"
)

var ErrValueExceedsMax = fault.New("operation would exceed the maximum value", fault.WithCode(fault.Conflict))

type BoundedValue struct {
	current int64
	max     int64
}

var ZeroBoundedValue = BoundedValue{}

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

func (bv BoundedValue) Current() int64 {
	return bv.current
}

func (bv BoundedValue) Max() int64 {
	return bv.max
}

func (bv BoundedValue) AvailableSpace() int64 {
	return bv.max - bv.current
}

func (bv BoundedValue) IsFull() bool {
	return bv.current == bv.max
}

func (bv BoundedValue) IsZero() bool {
	return bv.current == 0 && bv.max == 0
}

func (bv BoundedValue) Add(amount int64) (BoundedValue, error) {
	if amount < 0 {
		return ZeroBoundedValue, fault.New("amount to add must be non-negative", fault.WithCode(fault.Invalid))
	}
	if bv.AvailableSpace() < amount {
		return ZeroBoundedValue, ErrValueExceedsMax
	}
	return BoundedValue{current: bv.current + amount, max: bv.max}, nil
}

func (bv BoundedValue) Subtract(amount int64) (BoundedValue, error) {
	if amount < 0 {
		return ZeroBoundedValue, fault.New("amount to subtract must be non-negative", fault.WithCode(fault.Invalid))
	}
	if bv.current < amount {
		return ZeroBoundedValue, fault.New("cannot subtract more than the current value (would be negative)", fault.WithCode(fault.Invalid))
	}
	return BoundedValue{current: bv.current - amount, max: bv.max}, nil
}

func (bv BoundedValue) Set(newValue int64) (BoundedValue, error) {
	if newValue < 0 || newValue > bv.max {
		return ZeroBoundedValue, fault.New("new value is outside the allowed [0, max] range", fault.WithCode(fault.Invalid))
	}
	return BoundedValue{current: newValue, max: bv.max}, nil
}

func (bv BoundedValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Current int64 `json:"current"`
		Max     int64 `json:"max"`
	}{
		Current: bv.current,
		Max:     bv.max,
	})
}

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
