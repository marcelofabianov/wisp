package wisp

import (
	"encoding/json"

	"github.com/marcelofabianov/fault"
)

type DiscountType string

const (
	FixedDiscount      DiscountType = "fixed"
	PercentageDiscount DiscountType = "percentage"
)

type Discount struct {
	discountType    DiscountType
	fixedValue      Money
	percentageValue Percentage
}

var ZeroDiscount = Discount{}

func NewFixedDiscount(value Money) (Discount, error) {
	if value.IsNegative() {
		return ZeroDiscount, fault.New("fixed discount value cannot be negative", fault.WithCode(fault.Invalid))
	}
	return Discount{
		discountType: FixedDiscount,
		fixedValue:   value,
	}, nil
}

func NewPercentageDiscount(value Percentage) (Discount, error) {
	if value.IsNegative() || value.Float64() > 1.0 {
		return ZeroDiscount, fault.New("percentage discount must be between 0% and 100%", fault.WithCode(fault.Invalid))
	}
	return Discount{
		discountType:    PercentageDiscount,
		percentageValue: value,
	}, nil
}

func (d Discount) ApplyTo(m Money) (Money, error) {
	if d.IsZero() {
		return m, nil
	}

	var discountAmount Money
	var err error

	switch d.discountType {
	case FixedDiscount:
		if m.Currency() != d.fixedValue.Currency() {
			return ZeroMoney, fault.New("cannot apply fixed discount with different currency", fault.WithCode(fault.DomainViolation))
		}
		discountAmount = d.fixedValue
	case PercentageDiscount:
		discountAmount = d.percentageValue.ApplyTo(m)
	default:
		return m, nil
	}

	result, err := m.Subtract(discountAmount)
	if err != nil {
		return ZeroMoney, err
	}

	if result.IsNegative() {
		return NewMoney(0, m.Currency())
	}

	return result, nil
}

func (d Discount) String() string {
	if d.IsZero() {
		return "No Discount"
	}
	switch d.discountType {
	case FixedDiscount:
		return d.fixedValue.String()
	case PercentageDiscount:
		return d.percentageValue.String()
	}
	return ""
}

func (d Discount) IsZero() bool {
	return d.discountType == ""
}

func (d Discount) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return json.Marshal(nil)
	}

	var value any
	if d.discountType == FixedDiscount {
		value = d.fixedValue
	} else {
		value = d.percentageValue.Float64()
	}

	return json.Marshal(&struct {
		Type  DiscountType `json:"type"`
		Value any          `json:"value"`
	}{
		Type:  d.discountType,
		Value: value,
	})
}

func (d *Discount) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*d = ZeroDiscount
		return nil
	}

	dto := &struct {
		Type  DiscountType    `json:"type"`
		Value json.RawMessage `json:"value"`
	}{}

	if err := json.Unmarshal(data, dto); err != nil {
		return fault.Wrap(err, "invalid JSON format for Discount", fault.WithCode(fault.Invalid))
	}

	var newDiscount Discount
	var err error

	switch dto.Type {
	case FixedDiscount:
		var m Money
		if err = json.Unmarshal(dto.Value, &m); err != nil {
			return fault.Wrap(err, "invalid money format for fixed discount value", fault.WithCode(fault.Invalid))
		}
		newDiscount, err = NewFixedDiscount(m)
	case PercentageDiscount:
		var p float64
		if err = json.Unmarshal(dto.Value, &p); err != nil {
			return fault.Wrap(err, "invalid number format for percentage discount value", fault.WithCode(fault.Invalid))
		}
		perc, pErr := NewPercentageFromFloat(p)
		if pErr != nil {
			return pErr
		}
		newDiscount, err = NewPercentageDiscount(perc)
	default:
		err = fault.New("invalid discount type in JSON", fault.WithCode(fault.Invalid), fault.WithContext("type", dto.Type))
	}

	if err != nil {
		return err
	}
	*d = newDiscount
	return nil
}
