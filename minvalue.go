package wisp

import (
	"encoding/json"

	"github.com/marcelofabianov/fault"
)

var ErrValueSubceedsMin = fault.New("operation would fall below the minimum value", fault.WithCode(fault.Conflict))

type MinValue struct {
	current int64
	min     int64
}

var ZeroMinValue = MinValue{}

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

func (mv MinValue) Current() int64 {
	return mv.current
}

func (mv MinValue) Min() int64 {
	return mv.min
}

func (mv MinValue) IsAtMin() bool {
	return mv.current == mv.min
}

func (mv MinValue) Add(amount int64) (MinValue, error) {
	if amount < 0 {
		return ZeroMinValue, fault.New("amount to add must be non-negative", fault.WithCode(fault.Invalid))
	}
	return MinValue{current: mv.current + amount, min: mv.min}, nil
}

func (mv MinValue) Subtract(amount int64) (MinValue, error) {
	if amount < 0 {
		return ZeroMinValue, fault.New("amount to subtract must be non-negative", fault.WithCode(fault.Invalid))
	}
	if mv.current-mv.min < amount {
		return ZeroMinValue, ErrValueSubceedsMin
	}
	return MinValue{current: mv.current - amount, min: mv.min}, nil
}

func (mv MinValue) Set(newValue int64) (MinValue, error) {
	if newValue < mv.min {
		return ZeroMinValue, fault.New("new value cannot be less than the minimum", fault.WithCode(fault.Invalid))
	}
	return MinValue{current: newValue, min: mv.min}, nil
}

func (mv MinValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Current int64 `json:"current"`
		Min     int64 `json:"min"`
	}{
		Current: mv.current,
		Min:     mv.min,
	})
}

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
