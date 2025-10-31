package schema

import (
	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Machine holds the schema definition for the Machine entity.
type Machine struct {
	ent.Schema
}

func (Machine) Mixin() []ent.Mixin {
	return []ent.Mixin{mixin.Time{}}
}

func (Machine) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(func() uuid.UUID { return uuid.New() }),
		field.UUID("tenant_id", uuid.UUID{}).Immutable(),
		field.String("name"),
		field.Enum("category").Values("personal_laptop", "server"),
		field.String("hostname").Optional(),
		field.String("role").Optional(),
		field.Time("last_scan_at").Optional(),
	}
}

func (Machine) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("machines").
			Field("tenant_id").
			Required().
			Unique(),
	}
}
