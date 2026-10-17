package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Permission struct {
	ent.Schema
}

func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("module_id", uuid.UUID{}),
		field.String("name").NotEmpty(),
		field.String("key").NotEmpty().Unique(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("modules",Module.Type).Field("module_id").Ref("permissions").Required().Unique(),
		edge.To("rolepermissions",RolePermission.Type),
	}
}

func (Permission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("module_id"),
	}
}