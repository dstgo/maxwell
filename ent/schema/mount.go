package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/dstgo/maxwell/pkg/ts"
)

// Mount holds the schema definition for the Mount entity.
type Mount struct {
	ent.Schema
}

func (Mount) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("container mount record table"),
	}
}

// Fields of the Mount.
func (Mount) Fields() []ent.Field {
	return []ent.Field{
		field.String("host"),
		field.String("bind"),
		field.Bool("is_dir"),
		field.Int64("created_at").DefaultFunc(ts.UnixMicro),
		field.Int64("updated_at").DefaultFunc(ts.UnixMicro).UpdateDefault(ts.UnixMicro),
	}
}

// Edges of the Mount.
func (Mount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Container.Type).Ref("mounts"),
	}
}
