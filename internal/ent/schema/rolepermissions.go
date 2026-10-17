package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type RolePermission struct {
	ent.Schema
}

func (RolePermission) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("tenant_id", uuid.UUID{}).Optional(),
		field.UUID("role_id", uuid.UUID{}),
		field.UUID("permission_id", uuid.UUID{}),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (RolePermission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role",Role.Type).Field("role_id").Ref("permissions").Unique().Required(),
		edge.From("permission",Permission.Type).Field("permission_id").Ref("rolepermissions").Unique().Required(),
		edge.From("tenant",Tenant.Type).Field("tenant_id").Ref("rolepermissions").Unique(),
	}
}

func (RolePermission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id","permission_id").Unique(),
	}
}