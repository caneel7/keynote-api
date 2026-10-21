package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type UserNodeAssignment struct {
	ent.Schema
}

func (UserNodeAssignment) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("tenant_id", uuid.UUID{}).Immutable(),
		field.UUID("node_id", uuid.UUID{}).Immutable(),
		field.UUID("user_id", uuid.UUID{}).Immutable(),
		field.UUID("role_id", uuid.UUID{}),
		field.UUID("created_by", uuid.UUID{}).Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (UserNodeAssignment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("assignments").Field("tenant_id").Unique().Required().Immutable(),
		edge.From("node", OrgNode.Type).Ref("assignments").Field("node_id").Unique().Required().Immutable(),
		edge.From("user", User.Type).Ref("node_assignments").Field("user_id").Unique().Required().Immutable(),
		edge.From("role", Role.Type).Ref("assignments").Field("role_id").Unique().Required(),
		edge.From("creator", User.Type).Ref("created_assignments").Field("created_by").Unique(),
	}
}

func (UserNodeAssignment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "node_id").Unique().
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
		index.Fields("user_id", "tenant_id"),
		index.Fields("node_id"),
	}
}
