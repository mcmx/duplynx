package schema

import (
	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Machine holds the schema definition for the Machine entity.
type Machine struct {
	ent.Schema
}

func (Machine) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

// Fields of the Machine.
func (Machine) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(func() uuid.UUID { return uuid.New() }),
		field.UUID("tenant_id", uuid.UUID{}).Optional(),
		field.String("name"),
		field.Enum("category").Values("personal_laptop", "server"),
		field.String("hostname").Optional(),
		field.String("role").Optional(),
		field.Time("last_scan_at").Optional(),
	}
}

// Edges of the Machine.
func (Machine) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("machines").
			Field("tenant_id").
			Unique(),
		edge.To("initiated_scans", Scan.Type),
		edge.To("keeper_groups", DuplicateGroup.Type),
		edge.To("file_instances", FileInstance.Type),
	}
}
