package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// CookingHistory holds the schema definition for the CookingHistory entity.
type CookingHistory struct {
	ent.Schema
}

// Fields of the CookingHistory.
func (CookingHistory) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.Time("cookedAt"),
	}
}

// Edges of the CookingHistory.
func (CookingHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("recipe", Recipe.Type).
			Ref("history").
			Required().
			Unique(),
	}
}
