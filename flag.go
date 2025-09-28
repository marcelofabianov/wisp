package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/marcelofabianov/fault"
)

type Flag[T comparable] struct {
	value          T
	primaryValue   T
	secondaryValue T
}

func NewFlag[T comparable](currentValue, primaryValue, secondaryValue T) (Flag[T], error) {
	if currentValue != primaryValue && currentValue != secondaryValue {
		return Flag[T]{}, fault.New(
			"current value is not one of the allowed flag values",
			fault.WithCode(fault.Invalid),
			fault.WithContext("current_value", currentValue),
			fault.WithContext("allowed_values", []T{primaryValue, secondaryValue}),
		)
	}

	return Flag[T]{
		value:          currentValue,
		primaryValue:   primaryValue,
		secondaryValue: secondaryValue,
	}, nil
}

func (f Flag[T]) Get() T {
	return f.value
}

func (f Flag[T]) IsPrimary() bool {
	return f.value == f.primaryValue
}

func (f Flag[T]) IsSecondary() bool {
	return f.value == f.secondaryValue
}

func (f Flag[T]) Is(value T) bool {
	return f.value == value
}

func (f Flag[T]) String() string {
	return fmt.Sprintf("%v", f.value)
}

func (f Flag[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.value)
}

func (f Flag[T]) Value() (driver.Value, error) {
	return f.value, nil
}

func (f *Flag[T]) Scan(src interface{}) error {
	if src == nil {
		var zero T
		f.value = zero
		return nil
	}

	val, ok := src.(T)
	if !ok {
		return fault.New("unsupported scan type for Flag", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}
	f.value = val

	return nil
}
