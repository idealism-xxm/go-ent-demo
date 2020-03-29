package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Article holds the schema definition for the Article entity.
type Article struct {
	ent.Schema
}

func (Article) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Article.
func (Article) Fields() []ent.Field {
	return []ent.Field{
		// .MaxLen(...) 限制的长度是字节长度，一个汉字占 3 个字节
		field.String("title").MaxLen(32).Comment("标题"),
		field.Text("content").Comment("内容"),
	}
}

// Edges of the Article.
func (Article) Edges() []ent.Edge {
	return []ent.Edge{
		// 正向边，
		// 不使用 Unique 时会使用一张表存储对应关系，
		// 使用 Unique 时会成为外键
		edge.To("author", User.Type).Unique().Comment("作者"),
	}
}

func (Article) Config() ent.Config {
	return ent.Config{
		Table: "article",
	}
}
