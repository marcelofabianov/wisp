package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marcelofabianov/fault"
)

// CreatedAt is a value object that represents the timestamp when an entity was created.
// It is an alias for time.Time, automatically set to the current UTC time upon creation.
// This type is typically used in audit trails or for tracking record creation.
//
// Example:
//   myObject.CreatedAt = wisp.NewCreatedAt()
type CreatedAt time.Time

// NewCreatedAt creates a new CreatedAt timestamp, capturing the current time in UTC.
func NewCreatedAt() CreatedAt {
	return CreatedAt(time.Now().UTC())
}

// Time returns the underlying time.Time value.
func (c CreatedAt) Time() time.Time {
	return time.Time(c)
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the CreatedAt timestamp into a standard JSON time format.
func (c CreatedAt) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Time())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON time string into a CreatedAt timestamp.
func (c *CreatedAt) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return fault.Wrap(err, "CreatedAt must be a valid JSON time string", fault.WithCode(fault.Invalid))
	}
	*c = CreatedAt(t)
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the CreatedAt timestamp as a time.Time value.
func (c CreatedAt) Value() (driver.Value, error) {
	return c.Time(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a time.Time from the database and converts it into a CreatedAt timestamp.
func (c *CreatedAt) Scan(src interface{}) error {
	switch v := src.(type) {
	case time.Time:
		*c = CreatedAt(v)
		return nil
	default:
		return fault.New("unsupported scan type for CreatedAt", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}
}
