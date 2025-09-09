package wisp

import "time"

type Audit struct {
	CreatedAt  CreatedAt    `db:"audit_created_at" json:"created_at"`
	CreatedBy  AuditUser    `db:"audit_created_by" json:"created_by"`
	UpdatedAt  UpdatedAt    `db:"audit_updated_at" json:"updated_at"`
	UpdatedBy  AuditUser    `db:"audit_updated_by" json:"updated_by"`
	ArchivedAt NullableTime `db:"audit_archived_at" json:"archived_at,omitempty"`
	DeletedAt  NullableTime `db:"audit_deleted_at" json:"deleted_at,omitempty"`
	Version    Version      `db:"audit_version" json:"version"`
}

func NewAudit(createdBy AuditUser) Audit {
	now := NewCreatedAt()
	return Audit{
		CreatedAt: now,
		CreatedBy: createdBy,
		UpdatedAt: UpdatedAt(now.Time()),
		UpdatedBy: createdBy,
		Version:   InitialVersion(),
	}
}

func (a *Audit) Touch(updatedBy AuditUser) {
	a.UpdatedAt.Touch()
	a.UpdatedBy = updatedBy
	a.Version.Increment()
}

func (a *Audit) Archive(archivedBy AuditUser) {
	a.ArchivedAt = NewNullableTime(time.Now().UTC())
	a.Touch(archivedBy)
}

func (a *Audit) Delete(deletedBy AuditUser) {
	a.DeletedAt = NewNullableTime(time.Now().UTC())
	a.Touch(deletedBy)
}
