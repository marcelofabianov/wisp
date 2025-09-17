package wisp_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/marcelofabianov/wisp"
)

type AuditSuite struct {
	suite.Suite
	user1  wisp.AuditUser
	user2  wisp.AuditUser
	system wisp.AuditUser
}

func (s *AuditSuite) SetupSuite() {
	s.user1, _ = wisp.NewAuditUser("user1@example.com")
	s.user2, _ = wisp.NewAuditUser("user2@example.com")
	s.system, _ = wisp.NewAuditUser("system")
}

func TestAuditSuite(t *testing.T) {
	suite.Run(t, new(AuditSuite))
}

func (s *AuditSuite) TestNewAudit() {
	s.Run("should correctly initialize a new audit record", func() {
		audit := wisp.NewAudit(s.user1)

		s.False(audit.CreatedAt.Time().IsZero(), "CreatedAt should be set")
		s.Equal(s.user1, audit.CreatedBy, "CreatedBy should be the provided user")
		s.Equal(audit.CreatedAt.Time(), audit.UpdatedAt.Time(), "UpdatedAt should initially be the same as CreatedAt")
		s.Equal(s.user1, audit.UpdatedBy, "UpdatedBy should initially be the same as CreatedBy")
		s.Equal(wisp.Version(1), audit.Version, "Version should start at 1")
		s.True(audit.ArchivedAt.IsZero(), "ArchivedAt should be zero")
		s.True(audit.DeletedAt.IsZero(), "DeletedAt should be zero")
	})
}

func (s *AuditSuite) TestAudit_Touch() {
	audit := wisp.NewAudit(s.user1)
	originalUpdatedAt := audit.UpdatedAt
	originalVersion := audit.Version

	time.Sleep(10 * time.Millisecond)

	audit.Touch(s.user2)

	s.True(audit.UpdatedAt.Time().After(originalUpdatedAt.Time()), "UpdatedAt should be updated")
	s.Equal(s.user2, audit.UpdatedBy, "UpdatedBy should be the new user")
	s.Equal(originalVersion+1, audit.Version, "Version should be incremented")
	s.Equal(s.user1, audit.CreatedBy)
}

func (s *AuditSuite) TestAudit_Archive() {
	audit := wisp.NewAudit(s.user1)
	originalUpdatedAt := audit.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	audit.Archive(s.user2)

	s.False(audit.ArchivedAt.IsZero(), "ArchivedAt should be set")
	s.True(audit.ArchivedAt.Time.After(originalUpdatedAt.Time()))
	s.True(audit.UpdatedAt.Time().After(originalUpdatedAt.Time()), "UpdatedAt should be touched on archive")
	s.Equal(s.user2, audit.UpdatedBy, "UpdatedBy should be the user who archived")
	s.Equal(wisp.Version(2), audit.Version, "Version should be incremented on archive")
}

func (s *AuditSuite) TestAudit_Unarchive() {
	audit := wisp.NewAudit(s.user1)
	audit.Archive(s.user2)
	originalUpdatedAt := audit.UpdatedAt
	originalVersion := audit.Version

	time.Sleep(10 * time.Millisecond)

	audit.Unarchive(s.system)

	s.True(audit.ArchivedAt.IsZero(), "ArchivedAt should be zeroed out")
	s.True(audit.UpdatedAt.Time().After(originalUpdatedAt.Time()), "UpdatedAt should be touched on unarchive")
	s.Equal(s.system, audit.UpdatedBy, "UpdatedBy should be the user who unarchived")
	s.Equal(originalVersion+1, audit.Version, "Version should be incremented on unarchive")
}

func (s *AuditSuite) TestAudit_Delete() {
	audit := wisp.NewAudit(s.user1)
	originalUpdatedAt := audit.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	audit.Delete(s.system)

	s.False(audit.DeletedAt.IsZero(), "DeletedAt should be set")
	s.True(audit.DeletedAt.Time.After(originalUpdatedAt.Time()))
	s.True(audit.UpdatedAt.Time().After(originalUpdatedAt.Time()), "UpdatedAt should be touched on delete")
	s.Equal(s.system, audit.UpdatedBy, "UpdatedBy should be the user who deleted")
	s.Equal(wisp.Version(2), audit.Version, "Version should be incremented on delete")
}

func (s *AuditSuite) TestAudit_Undelete() {
	audit := wisp.NewAudit(s.user1)
	audit.Delete(s.user2)
	originalUpdatedAt := audit.UpdatedAt
	originalVersion := audit.Version

	time.Sleep(10 * time.Millisecond)

	audit.Undelete(s.system)

	s.True(audit.DeletedAt.IsZero(), "DeletedAt should be zeroed out")
	s.True(audit.UpdatedAt.Time().After(originalUpdatedAt.Time()), "UpdatedAt should be touched on undelete")
	s.Equal(s.system, audit.UpdatedBy, "UpdatedBy should be the user who undeleted")
	s.Equal(originalVersion+1, audit.Version, "Version should be incremented on undelete")
}

func (s *AuditSuite) TestAudit_States() {
	activeAudit := wisp.NewAudit(s.user1)
	archivedAudit := wisp.NewAudit(s.user1)
	archivedAudit.Archive(s.user2)
	deletedAudit := wisp.NewAudit(s.user1)
	deletedAudit.Delete(s.user2)
	undeletedAudit := wisp.NewAudit(s.user1)
	undeletedAudit.Delete(s.user2)
	undeletedAudit.Undelete(s.user1)
	unarchivedAudit := wisp.NewAudit(s.user1)
	unarchivedAudit.Archive(s.user2)
	unarchivedAudit.Unarchive(s.user1)

	s.Run("IsActive", func() {
		s.True(activeAudit.IsActive())
		s.False(archivedAudit.IsActive())
		s.False(deletedAudit.IsActive())
	})

	s.Run("IsArchived", func() {
		s.False(activeAudit.IsArchived())
		s.True(archivedAudit.IsArchived())
		s.False(deletedAudit.IsArchived())
		s.False(unarchivedAudit.IsArchived())
	})

	s.Run("IsDeleted", func() {
		s.False(activeAudit.IsDeleted())
		s.False(archivedAudit.IsDeleted())
		s.True(deletedAudit.IsDeleted())
		s.False(undeletedAudit.IsDeleted())
	})
}
