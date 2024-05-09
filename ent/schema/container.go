package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/dstgo/maxwell/pkg/ts"
)

// Container holds the schema definition for the Container entity.
type Container struct {
	ent.Schema
}

func (Container) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("container info table"),
	}
}

// Fields of the Container.
func (Container) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("cid").Unique(),
		field.String("image"),
		field.Int64("cpu"),
		field.Int64("memory"),
		field.Int64("created_at").DefaultFunc(ts.UnixMicro),
		field.Int64("updated_at").DefaultFunc(ts.UnixMicro).UpdateDefault(ts.UnixMicro),
	}
}

// Edges of the Container.
func (Container) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("containers"),
		edge.From("node", Node.Type).Ref("containers"),
		edge.To("mounts", Mount.Type),
		edge.To("ports", Port.Type),
		edge.To("jobs", Job.Type),
	}
}
