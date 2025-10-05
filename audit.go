package wisp

import "time"

// Audit is an embeddable struct that provides a standard set of fields for tracking
// the lifecycle of an entity. It includes timestamps and user identifiers for creation,
// updates, archival, and deletion, as well as a version number for optimistic locking.
//
// By embedding this struct into other domain models, you can easily add comprehensive
// auditing capabilities.
//
// Example:
//   type Product struct {
//       ID wisp.UUID
//       Name string
//       wisp.Audit
//   }
//
//   // Creating a new product with audit trail
//   prod := Product{
//       ID: wisp.MustNewUUID(),
//       Name: "New Gadget",
//       Audit: wisp.NewAudit(adminUser.ID),
//   }
//
//   // Updating the product
//   prod.Audit.Touch(editorUser.ID)
type Audit struct {
	CreatedAt  CreatedAt    `db:"audit_created_at" json:"created_at"`
	CreatedBy  AuditUser    `db:"audit_created_by" json:"created_by"`
	UpdatedAt  UpdatedAt    `db:"audit_updated_at" json:"updated_at"`
	UpdatedBy  AuditUser    `db:"audit_updated_by" json:"updated_by"`
	ArchivedAt NullableTime `db:"audit_archived_at" json:"archived_at,omitempty"`
	DeletedAt  NullableTime `db:"audit_deleted_at" json:"deleted_at,omitempty"`
	Version    Version      `db:"audit_version" json:"version"`
}

// NewAudit creates a new Audit instance, initializing it for a newly created entity.
// It sets the creation and update timestamps to the current time, and the created/updated user
// to the provided user ID. The version is initialized to 1.
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

// Touch updates the audit trail for a modification.
// It sets the `UpdatedAt` timestamp to the current time, records the user who made the change,
// and increments the version number.
func (a *Audit) Touch(updatedBy AuditUser) {
	a.UpdatedAt.Touch()
	a.UpdatedBy = updatedBy
	a.Version = a.Version.Increment()
}

// Archive marks the entity as archived.
// It sets the `ArchivedAt` timestamp and calls `Touch` to update the modification trail.
func (a *Audit) Archive(archivedBy AuditUser) {
	a.ArchivedAt = NewNullableTime(time.Now().UTC())
	a.Touch(archivedBy)
}

// Unarchive removes the archived status from the entity.
// It nullifies the `ArchivedAt` timestamp and calls `Touch`.
func (a *Audit) Unarchive(updatedBy AuditUser) {
	a.ArchivedAt = NullableTime{}
	a.Touch(updatedBy)
}

// Delete marks the entity as deleted (soft delete).
// It sets the `DeletedAt` timestamp and calls `Touch`.
func (a *Audit) Delete(deletedBy AuditUser) {
	a.DeletedAt = NewNullableTime(time.Now().UTC())
	a.Touch(deletedBy)
}

// Undelete removes the deleted status from the entity.
// It nullifies the `DeletedAt` timestamp and calls `Touch`.
func (a *Audit) Undelete(updatedBy AuditUser) {
	a.DeletedAt = NullableTime{}
	a.Touch(updatedBy)
}

// IsArchived returns true if the entity has been archived.
func (a *Audit) IsArchived() bool {
	return !a.ArchivedAt.IsZero()
}

// IsDeleted returns true if the entity has been soft-deleted.
func (a *Audit) IsDeleted() bool {
	return !a.DeletedAt.IsZero()
}

// IsActive returns true if the entity is neither archived nor deleted.
func (a *Audit) IsActive() bool {
	return !a.IsArchived() && !a.IsDeleted()
}
