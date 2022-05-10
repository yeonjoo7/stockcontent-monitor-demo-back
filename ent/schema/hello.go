package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Hello holds the schema definition for the Hello entity.
type Hello struct {
	ent.Schema
}

// Fields of the Hello.
func (Hello) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("value").MaxLen(300),
	}
}

// Edges of the Hello.
func (Hello) Edges() []ent.Edge {
	return nil
}
