package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Subscription struct {
	ent.Schema
}

func (Subscription) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id",uuid.UUID{}).Default(uuid.New),
		field.UUID("tenant_id",uuid.UUID{}),
		field.UUID("plan_id",uuid.UUID{}),
		field.UUID("plan_version_id",uuid.UUID{}),
		field.UUID("status_id",uuid.UUID{}),
		field.Time("started_at").Default(time.Now).Immutable(),
		field.Time("current_period_end"),
		field.Bool("cancel_at_period_end").Default(false),
		field.Time("canceled_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Subscription) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant",Tenant.Type).Field("tenant_id").Ref("subscription").Unique().Required(),
		edge.From("plan",Plan.Type).Field("plan_id").Ref("subscription").Unique().Required(),
		edge.From("planversion",PlanVersion.Type).Field("plan_version_id").Ref("subscription").Unique().Required(),
		edge.From("statedefinition",StateDefinitions.Type).Field("status_id").Ref("subscription").Unique().Required(),
	}
}