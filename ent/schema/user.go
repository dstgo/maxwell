package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/dstgo/maxwell/pkg/ids"
	"github.com/dstgo/maxwell/pkg/ts"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("uid").DefaultFunc(ids.ULID).Unique(),
		field.String("username").Unique(),
		field.String("email").Unique(),
		field.String("password"),
		field.Int64("created_at").DefaultFunc(ts.Ts),
		field.Int64("updated_at").DefaultFunc(ts.Ts).UpdateDefault(ts.Ts),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
