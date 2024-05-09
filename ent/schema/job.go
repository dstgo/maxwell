package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/dstgo/maxwell/pkg/ts"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

func (Job) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("container cron job tables"),
	}
}

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("description"),
		field.Int("job_type"),
		field.Int64("created_at").DefaultFunc(ts.UnixMicro),
		field.Int64("updated_at").DefaultFunc(ts.UnixMicro).UpdateDefault(ts.UnixMicro),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Container.Type).Ref("jobs"),
	}
}
