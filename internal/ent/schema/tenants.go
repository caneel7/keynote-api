package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Tenant struct {
	ent.Schema
}

func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Immutable(),
		field.String("slug").SchemaType(map[string]string{
			dialect.Postgres: "citext",
		}).NotEmpty().Unique(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Tenant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("workspaces", Workspace.Type),
		edge.To("projects", Project.Type),
		edge.To("roles", Role.Type),
		edge.To("rolepermissions", RolePermission.Type),
		edge.To("assignments", UserNodeAssignment.Type),
		edge.To("subscription", Subscription.Type),
	}
}
