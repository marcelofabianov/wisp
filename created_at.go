package wisp

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marcelofabianov/fault"
)

type CreatedAt time.Time

func NewCreatedAt() CreatedAt {
	return CreatedAt(time.Now().UTC())
}

func (c CreatedAt) Time() time.Time {
	return time.Time(c)
}

func (c CreatedAt) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Time())
}

func (c *CreatedAt) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return fault.Wrap(err, "CreatedAt must be a valid JSON time string", fault.WithCode(fault.Invalid))
	}
	*c = CreatedAt(t)
	return nil
}

func (c CreatedAt) Value() (driver.Value, error) {
	return c.Time(), nil
}

func (c *CreatedAt) Scan(src interface{}) error {
	switch v := src.(type) {
	case time.Time:
		*c = CreatedAt(v)
		return nil
	default:
		return fault.New("unsupported scan type for CreatedAt", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}
}
