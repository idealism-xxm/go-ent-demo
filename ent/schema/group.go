package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Group holds the schema definition for the Group entity.
type Group struct {
	ent.Schema
}

// Fields of the Group.
func (Group) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().Comment("组名"),
	}
}

// Edges of the Group.
func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		// group 和 user 是多对多，使用关联表
		edge.To("users", User.Type).Comment("拥有的成员列表"),
	}
}

func (Group) Config() ent.Config {
	return ent.Config{
		Table: "group",
	}
}
