package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)


type User struct {
	ent.Schema
}

func (User) Fields () []ent.Field {
	return []ent.Field {
		field.UUID("id",uuid.UUID{}).Default(uuid.New),
		field.String("email").NotEmpty().Unique(),
		field.UUID("status_id",uuid.UUID{}),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (User) Edges () []ent.Edge {
	return []ent.Edge {
		edge.To("credentials",UserCredential.Type).Unique(),
		edge.To("sessions",UserSession.Type),
		edge.From("state",StateDefinitions.Type).Field("status_id").Ref("user").Unique().Required(),
		edge.To("created_nodes",OrgNode.Type),
		edge.To("node_assignments",UserNodeAssignment.Type),
		edge.To("created_assignments",UserNodeAssignment.Type),
	}
}