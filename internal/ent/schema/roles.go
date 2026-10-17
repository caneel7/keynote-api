package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Role struct {
	ent.Schema
}

func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id",uuid.UUID{}).Default(uuid.New),
		field.UUID("tenant_id",uuid.UUID{}).Nillable(),
		field.UUID("project_id",uuid.UUID{}).Nillable(),
		field.String("name").NotEmpty(),
		field.Enum("type").Values("tenant","project").Default("tenant"),
		field.UUID("created_by",uuid.UUID{}).Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant",Tenant.Type).Field("tenant_id").Ref("roles").Unique().Required(),
		edge.From("project",Project.Type).Field("project_id").Ref("roles").Unique().Required(),
		edge.To("permissions",RolePermission.Type),
		edge.To("tenantrole",TenantMember.Type),
		edge.To("projectrole",ProjectMember.Type),
	}
}