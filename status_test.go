package wisp_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type StatusSuite struct {
	suite.Suite
}

func TestStatusSuite(t *testing.T) {
	suite.Run(t, new(StatusSuite))
}

func (s *StatusSuite) SetupTest() {
	wisp.ClearRegisteredStatuses()
}

func (s *StatusSuite) TestRegisterAndValidateStatus() {
	const (
		StatusActive   wisp.Status = "ACTIVE"
		StatusInactive wisp.Status = "INACTIVE"
		StatusPending  wisp.Status = "PENDING"
	)
	wisp.RegisterStatuses(StatusActive, StatusInactive, StatusPending)

	s.Run("NewStatus should create a valid status that is registered", func() {
		status, err := wisp.NewStatus("active")
		s.Require().NoError(err)
		s.Equal(StatusActive, status)
	})

	s.Run("NewStatus should fail for a status that is not registered", func() {
		_, err := wisp.NewStatus("BLOCKED")
		s.Require().Error(err)
	})

	s.Run("IsValid should work correctly", func() {
		s.True(wisp.Status("ACTIVE").IsValid())
		s.False(wisp.Status("ARCHIVED").IsValid())
	})
}
