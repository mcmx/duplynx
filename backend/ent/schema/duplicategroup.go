package schema

import (
	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// DuplicateGroup holds the schema definition for the DuplicateGroup entity.
type DuplicateGroup struct {
	ent.Schema
}

func (DuplicateGroup) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (DuplicateGroup) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(func() uuid.UUID { return uuid.New() }),
		field.UUID("tenant_id", uuid.UUID{}).Optional(),
		field.UUID("scan_id", uuid.UUID{}).Optional(),
		field.UUID("keeper_machine_id", uuid.UUID{}).Optional(),
		field.String("hash"),
		field.Enum("status").
			Values("review", "action_needed", "resolved", "archived").
			Default("review"),
		field.Int("file_count").Positive(),
		field.Int64("total_size_bytes"),
	}
}

func (DuplicateGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("scans").
			Field("tenant_id").
			Unique(),
		edge.From("scan", Scan.Type).
			Ref("duplicate_groups").
			Field("scan_id").
			Unique(),
		edge.From("keeper_machine", Machine.Type).
			Ref("keeper_groups").
			Field("keeper_machine_id").
			Unique().
			Optional(),
		edge.To("file_instances", FileInstance.Type),
		edge.To("action_audits", ActionAudit.Type),
	}
}
