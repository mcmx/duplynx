package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// FileInstance holds the schema definition for the FileInstance entity.
type FileInstance struct {
	ent.Schema
}

func (FileInstance) Mixin() []ent.Mixin {
	return []ent.Mixin{mixin.Time{}}
}

func (FileInstance) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(func() uuid.UUID { return uuid.New() }),
		field.UUID("duplicate_group_id", uuid.UUID{}),
		field.UUID("machine_id", uuid.UUID{}),
		field.String("path"),
		field.Int64("size_bytes").NonNegative(),
		field.String("checksum"),
		field.Time("last_seen_at").Default(time.Now),
		field.Bool("quarantined").Default(false),
	}
}

func (FileInstance) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("duplicate_group", DuplicateGroup.Type).
			Ref("file_instances").
			Field("duplicate_group_id").
			Required().
			Unique(),
		edge.From("machine", Machine.Type).
			Ref("file_instances").
			Field("machine_id").
			Required().
			Unique(),
	}
}
