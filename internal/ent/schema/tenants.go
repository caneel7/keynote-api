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
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("name").SchemaType(map[string]string{
			dialect.Postgres: "citext",
		}).NotEmpty(),
		field.String("slug").SchemaType(map[string]string{
			dialect.Postgres: "citext",
		}).NotEmpty().Unique(),
		field.UUID("status_id", uuid.UUID{}),
		field.UUID("created_by", uuid.UUID{}),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (Tenant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("state", StateDefinitions.Type).Ref("tenants").Field("status_id").Unique().Required(),
		edge.From("user", User.Type).Ref("tenant").Field("created_by").Unique().Required(),
		edge.To("workspace", Workspace.Type).Unique(),
		edge.To("projects", Project.Type),
		edge.To("roles",Role.Type),
		edge.To("rolepermissions",RolePermission.Type),
		edge.To("tenantrole",TenantMember.Type),
		edge.To("projectrole",ProjectMember.Type),
		edge.To("subscription",Subscription.Type),
	}
}
