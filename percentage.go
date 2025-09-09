package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"

	"github.com/marcelofabianov/fault"
)

type Percentage int64

var ZeroPercentage Percentage

const percentageFactor = 10000.0

func NewPercentageFromFloat(value float64) (Percentage, error) {
	if value < 0 {
		return ZeroPercentage, fault.New(
			"percentage value cannot be negative",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_value", value),
		)
	}

	scaledValue := math.RoundToEven(value * percentageFactor)
	return Percentage(scaledValue), nil
}

func (p Percentage) Float64() float64 {
	return float64(p) / percentageFactor
}

func (p Percentage) String() string {
	return fmt.Sprintf("%.2f%%", p.Float64()*100)
}

func (p Percentage) IsZero() bool {
	return p == ZeroPercentage
}

func (p Percentage) ApplyTo(m Money) Money {
	if m.IsZero() || p.IsZero() {
		return Money{amount: 0, currency: m.Currency()}
	}

	result := float64(m.Amount()) * p.Float64()
	roundedAmount := int64(math.RoundToEven(result))

	return Money{
		amount:   roundedAmount,
		currency: m.Currency(),
	}
}

func (p Percentage) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Float64())
}

func (p *Percentage) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return fault.Wrap(err, "Percentage must be a valid JSON number", fault.WithCode(fault.Invalid))
	}
	perc, err := NewPercentageFromFloat(f)
	if err != nil {
		return err
	}
	*p = perc
	return nil
}

func (p Percentage) Value() (driver.Value, error) {
	return int64(p), nil
}

func (p *Percentage) Scan(src interface{}) error {
	if src == nil {
		*p = ZeroPercentage
		return nil
	}

	var intVal int64
	switch v := src.(type) {
	case int64:
		intVal = v
	default:
		return fault.New("unsupported scan type for Percentage", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}

	if intVal < 0 {
		return fault.New("percentage from database cannot be negative", fault.WithCode(fault.Invalid), fault.WithContext("source_value", intVal))
	}

	*p = Percentage(intVal)
	return nil
}
