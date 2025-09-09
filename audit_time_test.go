package atomic_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/marcelofabianov/atomic"
	"github.com/stretchr/testify/suite"
)

type AuditTimeSuite struct {
	suite.Suite
}

func TestAuditTimeSuite(t *testing.T) {
	suite.Run(t, new(AuditTimeSuite))
}

func (s *AuditTimeSuite) TestCreatedAt() {
	s.Run("should create a non-zero timestamp", func() {
		ca := atomic.NewCreatedAt()
		s.False(ca.Time().IsZero())
	})

	s.Run("should marshal and unmarshal correctly", func() {
		now := time.Now().UTC().Truncate(time.Second)
		ca := atomic.CreatedAt(now)

		data, err := json.Marshal(ca)
		s.Require().NoError(err)

		expectedData, _ := json.Marshal(now)
		s.JSONEq(string(expectedData), string(data))

		var unmarshaledCA atomic.CreatedAt
		err = json.Unmarshal(data, &unmarshaledCA)
		s.Require().NoError(err)
		s.Equal(ca, unmarshaledCA)
	})

	s.Run("should handle database Value and Scan", func() {
		now := time.Now().UTC().Truncate(time.Second)
		ca := atomic.CreatedAt(now)

		val, err := ca.Value()
		s.Require().NoError(err)
		s.Equal(now, val)

		var scannedCA atomic.CreatedAt
		err = scannedCA.Scan(now)
		s.Require().NoError(err)
		s.Equal(ca, scannedCA)
	})
}

func (s *AuditTimeSuite) TestUpdatedAt() {
	s.Run("should create a non-zero timestamp", func() {
		ua := atomic.NewUpdatedAt()
		s.False(ua.Time().IsZero())
	})

	s.Run("Touch method should update the time", func() {
		ua := atomic.NewUpdatedAt()
		originalTime := ua.Time()
		time.Sleep(10 * time.Millisecond) // Ensure time moves forward
		ua.Touch()
		s.True(ua.Time().After(originalTime))
	})

	s.Run("should marshal and unmarshal correctly", func() {
		now := time.Now().UTC().Truncate(time.Second)
		ua := atomic.UpdatedAt(now)

		data, err := json.Marshal(ua)
		s.Require().NoError(err)

		expectedData, _ := json.Marshal(now)
		s.JSONEq(string(expectedData), string(data))

		var unmarshaledUA atomic.UpdatedAt
		err = json.Unmarshal(data, &unmarshaledUA)
		s.Require().NoError(err)
		s.Equal(ua, unmarshaledUA)
	})

	s.Run("should handle database Value and Scan", func() {
		now := time.Now().UTC().Truncate(time.Second)
		ua := atomic.UpdatedAt(now)

		val, err := ua.Value()
		s.Require().NoError(err)
		s.Equal(now, val)

		var scannedUA atomic.UpdatedAt
		err = scannedUA.Scan(now)
		s.Require().NoError(err)
		s.Equal(ua, scannedUA)
	})
}
