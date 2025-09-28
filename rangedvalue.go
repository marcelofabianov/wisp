package wisp

import (
	"encoding/json"

	"github.com/marcelofabianov/fault"
)

type RangedValue struct {
	current int64
	min     int64
	max     int64
}

var ZeroRangedValue = RangedValue{}

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

func (rv RangedValue) Current() int64 {
	return rv.current
}

func (rv RangedValue) Min() int64 {
	return rv.min
}

func (rv RangedValue) Max() int64 {
	return rv.max
}

func (rv RangedValue) IsAtMin() bool {
	return rv.current == rv.min
}

func (rv RangedValue) IsAtMax() bool {
	return rv.current == rv.max
}

func (rv RangedValue) Add(amount int64) (RangedValue, error) {
	if amount < 0 {
		return ZeroRangedValue, fault.New("amount to add must be non-negative", fault.WithCode(fault.Invalid))
	}
	if rv.max-rv.current < amount {
		return ZeroRangedValue, ErrValueExceedsMax
	}
	return RangedValue{current: rv.current + amount, min: rv.min, max: rv.max}, nil
}

func (rv RangedValue) Subtract(amount int64) (RangedValue, error) {
	if amount < 0 {
		return ZeroRangedValue, fault.New("amount to subtract must be non-negative", fault.WithCode(fault.Invalid))
	}
	if rv.current-rv.min < amount {
		return ZeroRangedValue, ErrValueSubceedsMin
	}
	return RangedValue{current: rv.current - amount, min: rv.min, max: rv.max}, nil
}

func (rv RangedValue) Set(newValue int64) (RangedValue, error) {
	if newValue < rv.min || newValue > rv.max {
		return ZeroRangedValue, fault.New("new value is outside the allowed range", fault.WithCode(fault.Invalid))
	}
	return RangedValue{current: newValue, min: rv.min, max: rv.max}, nil
}

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
