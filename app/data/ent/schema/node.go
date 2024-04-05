package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/dstgo/maxwell/app/pkg/ts"
	"github.com/oklog/ulid/v2"
)

// Node holds the schema definition for the Node entity.
type Node struct {
	ent.Schema
}

// Fields of the Node.
func (Node) Fields() []ent.Field {
	return []ent.Field{
		field.String("uid").DefaultFunc(ulid.Make).Unique(),
		field.String("name").Unique(),
		field.String("address").Unique(),
		field.String("note"),
		field.Int64("created_at").DefaultFunc(ts.Ts),
		field.Int64("updated_at").DefaultFunc(ts.Ts).UpdateDefault(ts.Ts),
	}
}

// Edges of the Node.
func (Node) Edges() []ent.Edge {
	return nil
}
