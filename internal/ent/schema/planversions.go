package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type PlanVersion struct {
	ent.Schema
}

func (PlanVersion) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id",uuid.UUID{}).Default(uuid.New),
		field.UUID("plan_id",uuid.UUID{}).Immutable(),
		field.Int("version").Default(1).NonNegative().Immutable(),
		field.Float("cost").Immutable(),
		field.String("currency").Default("usd"),
		field.Int("interval").Default(30).NonNegative().Immutable(),
		field.JSON("perks",map[string]int{}).Immutable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (PlanVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("plan",Plan.Type).Field("plan_id").Ref("versions").Unique().Required().Immutable(),
		edge.To("subscription",Subscription.Type),
	}
}