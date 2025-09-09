package atomic

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marcelofabianov/fault"
)

type UpdatedAt time.Time

func NewUpdatedAt() UpdatedAt {
	return UpdatedAt(time.Now().UTC())
}

func (u *UpdatedAt) Touch() {
	*u = UpdatedAt(time.Now().UTC())
}

func (u UpdatedAt) Time() time.Time {
	return time.Time(u)
}

func (u UpdatedAt) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Time())
}

func (u *UpdatedAt) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return fault.Wrap(err, "UpdatedAt must be a valid JSON time string", fault.WithCode(fault.Invalid))
	}
	*u = UpdatedAt(t)
	return nil
}

func (u UpdatedAt) Value() (driver.Value, error) {
	return u.Time(), nil
}

func (u *UpdatedAt) Scan(src interface{}) error {
	switch v := src.(type) {
	case time.Time:
		*u = UpdatedAt(v)
		return nil
	default:
		return fault.New("unsupported scan type for UpdatedAt", fault.WithCode(fault.Invalid), fault.WithContext("received_type", fmt.Sprintf("%T", src)))
	}
}
