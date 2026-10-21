package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Role struct {
	ent.Schema
}

func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id",uuid.UUID{}).Default(uuid.New),
		field.UUID("tenant_id",uuid.UUID{}).Optional().Nillable(),
		field.UUID("node_id",uuid.UUID{}).Optional().Nillable(),
		field.String("name").NotEmpty(),
		field.Enum("type").Values("tenant","project").Default("tenant"),
		field.UUID("created_by",uuid.UUID{}).Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant",Tenant.Type).Field("tenant_id").Ref("roles").Unique(),
		edge.From("node",OrgNode.Type).Field("node_id").Ref("roles").Unique(),
		edge.To("permissions",RolePermission.Type),
		edge.To("assignments",UserNodeAssignment.Type),
	}
}

func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("node_id"),
	}
}
