package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// ActionAudit records user driven actions for duplicates.
type ActionAudit struct {
	ent.Schema
}

func (ActionAudit) Mixin() []ent.Mixin {
	return []ent.Mixin{mixin.Time{}}
}

func (ActionAudit) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(func() uuid.UUID { return uuid.New() }),
		field.UUID("tenant_id", uuid.UUID{}).Immutable(),
		field.UUID("duplicate_group_id", uuid.UUID{}).Immutable(),
		field.String("actor").Default("system"),
		field.Enum("action_type").
			Values("assign_keeper", "delete_copies", "create_hardlinks", "quarantine", "retry", "note"),
		field.JSON("payload", map[string]any{}).Optional(),
		field.Time("performed_at").Default(time.Now),
		field.Bool("stubbed").Default(false),
	}
}

func (ActionAudit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("action_audits").
			Field("tenant_id").
			Required().
			Unique(),
		edge.From("duplicate_group", DuplicateGroup.Type).
			Ref("action_audits").
			Field("duplicate_group_id").
			Required().
			Unique(),
	}
}
