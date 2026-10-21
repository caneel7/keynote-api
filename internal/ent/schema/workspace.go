package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Workspace struct {
	ent.Schema
}

func (Workspace) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Immutable(),
		field.UUID("tenant_id", uuid.UUID{}).Immutable(),
		field.String("slug").SchemaType(map[string]string{
			dialect.Postgres: "citext",
		}).NotEmpty(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Workspace) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("workspaces").Field("tenant_id").Unique().Required().Immutable(),
	}
}

func (Workspace) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "id").Unique(),
		index.Fields("tenant_id", "slug").Unique(),
	}
}
