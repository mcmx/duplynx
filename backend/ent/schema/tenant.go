package schema

import (
	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Tenant holds the schema definition for the Tenant entity.
type Tenant struct {
	ent.Schema
}

func (Tenant) Mixin() []ent.Mixin {
	return []ent.Mixin{mixin.Time{}}
}

func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(func() uuid.UUID { return uuid.New() }),
		field.String("slug").Unique().Immutable(),
		field.String("name"),
		field.String("description").Optional(),
		field.String("primary_contact").Optional(),
	}
}

func (Tenant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("machines", Machine.Type),
		edge.To("scans", Scan.Type),
		edge.To("duplicate_groups", DuplicateGroup.Type),
		edge.To("action_audits", ActionAudit.Type),
	}
}
