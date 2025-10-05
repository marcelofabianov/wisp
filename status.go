package wisp

import (
	"strings"

	"github.com/marcelofabianov/fault"
)

type Status string

var validStatuses = make(map[Status]struct{})
var EmptyStatus Status

func RegisterStatuses(statuses ...Status) {
	for _, s := range statuses {
		normalized := Status(strings.ToUpper(strings.TrimSpace(string(s))))
		if normalized != "" {
			validStatuses[normalized] = struct{}{}
		}
	}
}

func NewStatus(value string) (Status, error) {
	normalized := Status(strings.ToUpper(strings.TrimSpace(value)))
	if normalized == EmptyStatus {
		return EmptyStatus, nil
	}

	if !normalized.IsValid() {
		return EmptyStatus, fault.New(
			"status is not registered as a valid status",
			fault.WithCode(fault.Invalid),
			fault.WithContext("input_status", value),
		)
	}
	return normalized, nil
}

func ClearRegisteredStatuses() {
	validStatuses = make(map[Status]struct{})
}

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	_, ok := validStatuses[s]
	return ok
}

func (s Status) IsZero() bool {
	return s == EmptyStatus
}
