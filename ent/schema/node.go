package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/dstgo/maxwell/pkg/ids"
	"github.com/dstgo/maxwell/pkg/ts"
)

// Node holds the schema definition for the Node entity.
type Node struct {
	ent.Schema
}

// Fields of the Node.
func (Node) Fields() []ent.Field {
	return []ent.Field{
		field.String("uid").DefaultFunc(ids.ULID).Unique(),
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
