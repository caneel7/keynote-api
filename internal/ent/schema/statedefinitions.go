package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type StateDefinitions struct {
	ent.Schema
}

func (StateDefinitions) Fields () []ent.Field {
	return []ent.Field{
		field.UUID("id",uuid.UUID{}).Default(uuid.New),
		field.UUID("entity_id",uuid.UUID{}),
		field.String("key").NotEmpty(),
		field.String("name").NotEmpty(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (StateDefinitions) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("entity_kind",EntityKind.Type).Ref("state_definitions").Field("entity_id").Unique().Required(),
		edge.To("tenants",Tenant.Type),
		edge.To("user",User.Type),
		edge.To("workspace",Workspace.Type),
		edge.To("projects",Project.Type),
		edge.To("subscription",Subscription.Type),
	}
}

func (StateDefinitions) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key").Edges("entity_kind").Unique(),
	}
}