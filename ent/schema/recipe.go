package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Recipe holds the schema definition for the Recipe entity.
type Recipe struct {
	ent.Schema
}

// Fields of the Recipe.
func (Recipe) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.String("title"),
	}
}

// Edges of the Recipe.
func (Recipe) Edges() []ent.Edge {
	return nil
}
