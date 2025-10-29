package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Scan holds the schema definition for the Scan entity.
type Scan struct {
	ent.Schema
}

func (Scan) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (Scan) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(func() uuid.UUID { return uuid.New() }),
		field.UUID("tenant_id", uuid.UUID{}).Optional(),
		field.UUID("initiated_machine_id", uuid.UUID{}).Optional(),
		field.String("name"),
		field.String("description").Optional(),
		field.Time("started_at").Default(time.Now),
		field.Time("completed_at").Optional(),
		field.Int("duplicate_group_count").Default(0),
	}
}

func (Scan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("scans").
			Field("tenant_id").
			Unique(),
		edge.From("initiated_machine", Machine.Type).
			Ref("initiated_scans").
			Field("initiated_machine_id").
			Unique(),
		edge.To("duplicate_groups", DuplicateGroup.Type),
	}
}
