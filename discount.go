package wisp

import (
	"encoding/json"

	"github.com/marcelofabianov/fault"
)

// DiscountType defines the type of discount, which can be either fixed or percentage-based.
type DiscountType string

// Defines the supported types of discounts.
const (
	FixedDiscount      DiscountType = "fixed"      // A fixed monetary amount (e.g., R$10.00 off).
	PercentageDiscount DiscountType = "percentage" // A percentage of the total amount (e.g., 15% off).
)

// Discount represents a value object for a discount, which can be either a fixed monetary amount
// or a percentage. This structure ensures that discounts are applied correctly and safely,
// handling different currencies and preventing invalid states.
//
// A Discount is immutable. Operations like applying it to a monetary value return a new result.
//
// Examples:
//   fixed, _ := NewFixedDiscount(NewMoney(1000, BRL)) // R$10.00 discount
//   percent, _ := NewPercentageDiscount(NewPercentageFromFloat(0.15)) // 15% discount
type Discount struct {
	discountType    DiscountType
	fixedValue      Money
	percentageValue Percentage
}

// ZeroDiscount represents the zero value for the Discount type (no discount).
var ZeroDiscount = Discount{}

// NewFixedDiscount creates a new discount with a fixed monetary value.
// Returns an error if the provided Money value is negative.
func NewFixedDiscount(value Money) (Discount, error) {
	if value.IsNegative() {
		return ZeroDiscount, fault.New("fixed discount value cannot be negative", fault.WithCode(fault.Invalid))
	}
	return Discount{
		discountType: FixedDiscount,
		fixedValue:   value,
	}, nil
}

// NewPercentageDiscount creates a new discount with a percentage value.
// The percentage must be between 0.0 (0%) and 1.0 (100%), inclusive.
// Returns an error if the percentage is outside this valid range.
func NewPercentageDiscount(value Percentage) (Discount, error) {
	if value.IsNegative() || value.Float64() > 1.0 {
		return ZeroDiscount, fault.New("percentage discount must be between 0% and 100%", fault.WithCode(fault.Invalid))
	}
	return Discount{
		discountType:    PercentageDiscount,
		percentageValue: value,
	}, nil
}

// ApplyTo applies the discount to a given Money value and returns the new amount.
// - For a fixed discount, it subtracts the fixed amount. Currencies must match.
// - For a percentage discount, it calculates and subtracts the percentage amount.
// If the resulting amount is negative, it is floored at zero.
// Returns an error if a fixed discount is applied to a different currency.
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

// String returns a string representation of the discount.
// For a fixed discount, it returns the formatted money string (e.g., "BRL 10.00").
// For a percentage discount, it returns the formatted percentage string (e.g., "15.00%").
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

// IsZero returns true if the discount is the zero value (no discount).
func (d Discount) IsZero() bool {
	return d.discountType == ""
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the Discount into a JSON object with "type" and "value" fields.
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

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object into a Discount, validating its type and value.
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
