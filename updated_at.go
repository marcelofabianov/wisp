package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marcelofabianov/fault"
)

// UpdatedAt is a value object that represents the timestamp when an entity was last updated.
// It is an alias for time.Time and provides methods to easily update the timestamp.
// This is essential for audit trails and optimistic concurrency control.
//
// Example:
//   myObject.UpdatedAt.Touch() // Updates the timestamp to the current time
type UpdatedAt time.Time

// NewUpdatedAt creates a new UpdatedAt timestamp, capturing the current time in UTC.
func NewUpdatedAt() UpdatedAt {
	return UpdatedAt(time.Now().UTC())
}

// Touch updates the UpdatedAt timestamp to the current time in UTC.
// This method should be called whenever the associated entity is modified.
func (u *UpdatedAt) Touch() {
	*u = UpdatedAt(time.Now().UTC())
}

// Time returns the underlying time.Time value.
func (u UpdatedAt) Time() time.Time {
	return time.Time(u)
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the UpdatedAt timestamp into a standard JSON time format.
func (u UpdatedAt) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Time())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON time string into an UpdatedAt timestamp.
func (u *UpdatedAt) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return fault.Wrap(err, "UpdatedAt must be a valid JSON time string", fault.WithCode(fault.Invalid))
	}
	*u = UpdatedAt(t)
	return nil
}

// Value implements the driver.Valuer interface for database storage.
// It returns the UpdatedAt timestamp as a time.Time value.
func (u UpdatedAt) Value() (driver.Value, error) {
	return u.Time(), nil
}

// Scan implements the sql.Scanner interface for database retrieval.
// It accepts a time.Time from the database and converts it into an UpdatedAt timestamp.
func (u *UpdatedAt) Scan(src interface{}) error {
	switch v := src.(type) {
	case time.Time:
		*u = UpdatedAt(v)
		return nil
	default:
		return fault.New("unsupported scan type for UpdatedAt", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}
}
