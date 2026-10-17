package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type ProjectMember struct {
	ent.Schema
}

func (ProjectMember) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("tenant_id", uuid.UUID{}),
		field.UUID("project_id", uuid.UUID{}),
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("role_id", uuid.UUID{}),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (ProjectMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Field("user_id").Ref("projectrole").Unique().Required(),
		edge.From("tenant", Tenant.Type).Field("tenant_id").Ref("projectrole").Unique().Required(),
		edge.From("role", Role.Type).Field("role_id").Ref("projectrole").Unique().Required(),
		edge.From("project", Project.Type).Field("project_id").Ref("projectrole").Unique().Required(),
	}
}

func (ProjectMember) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "project_id", "user_id").
			Unique(),
	}
}
