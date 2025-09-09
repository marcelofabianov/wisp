package wisp

type Audit struct {
	CreatedAt  CreatedAt    `db:"audit_created_at" json:"created_at"`
	CreatedBy  AuditUser    `db:"audit_created_by" json:"created_by"`
	UpdatedAt  UpdatedAt    `db:"audit_updated_at" json:"updated_at"`
	UpdatedBy  AuditUser    `db:"audit_updated_by" json:"updated_by"`
	ArchivedAt NullableTime `db:"audit_archived_at" json:"archived_at,omitempty"`
	DeletedAt  NullableTime `db:"audit_deleted_at" json:"deleted_at,omitempty"`
	Version    Version      `db:"audit_version" json:"version"`
}
