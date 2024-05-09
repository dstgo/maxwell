package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/dstgo/maxwell/pkg/ts"
)

// Port holds the schema definition for the Port entity.
type Port struct {
	ent.Schema
}

func (Port) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("container ports bind info table"),
	}
}

// Fields of the Port.
func (Port) Fields() []ent.Field {
	return []ent.Field{
		field.Int("host").Unique(),
		field.Int("bind").Comment("bind port in container"),
		field.Int("protocol").Comment("network protocol"),
		field.Int64("created_at").DefaultFunc(ts.UnixMicro),
		field.Int64("updated_at").DefaultFunc(ts.UnixMicro).UpdateDefault(ts.UnixMicro),
	}
}

// Edges of the Port.
func (Port) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Container.Type).Ref("ports"),
	}
}
