package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type OrgNode struct {
	ent.Schema
}

func (OrgNode) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("tenant_id", uuid.UUID{}).Immutable(),
		field.UUID("parent_id", uuid.UUID{}).Optional().Nillable(),
		field.Enum("type").Values("tenant", "workspace", "project").Immutable(),
		field.String("path").SchemaType(map[string]string{
			dialect.Postgres: "ltree",
		}).NotEmpty().Unique(),
		field.String("name").SchemaType(map[string]string{
			dialect.Postgres: "citext",
		}).NotEmpty(),
		field.UUID("status_id", uuid.UUID{}),
		field.UUID("created_by", uuid.UUID{}),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (OrgNode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tenant_nodes", OrgNode.Type).
			From("tenant_node").Field("tenant_id").Unique().Required().Immutable(),
		edge.To("children", OrgNode.Type).
			From("parent").Field("parent_id").Unique(),
		edge.From("state", StateDefinitions.Type).Ref("org_nodes").Field("status_id").Unique().Required(),
		edge.From("creator", User.Type).Ref("created_nodes").Field("created_by").Unique().Required(),
		edge.To("roles", Role.Type),
		edge.To("assignments", UserNodeAssignment.Type),
	}
}

func (OrgNode) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "id").Unique(),
		index.Fields("tenant_id", "type"),
		index.Fields("parent_id"),
		index.Fields("parent_id", "name").Unique().
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
		index.Fields("path").
			Annotations(entsql.IndexType("GIST")).
			StorageKey("orgnode_path_gist"),
	}
}
