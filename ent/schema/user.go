package schema

import (
	"errors"
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/index"
	"regexp"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	usernameReg := regexp.MustCompile("[a-zA-Z][a-zA-Z0-9-]{6,14}[a-zA-Z0-9]")
	return []ent.Field{
		field.String("nickname").MaxLen(16).Comment("昵称"),
		// 如果不需要自定义报错信息，可以直接使用 .Match(usernameReg) 即可
		field.String("username").MaxLen(16).Unique().Validate(func(s string) error {
			// 用户名只允许 英文字母、数字和 - ，且必须以英文字母开始
			if !usernameReg.MatchString(s) {
				return errors.New(
					"username may only contain alphanumeric characters or single hyphens " +
						"with length 8-16, and must start with alphabet and cannot end with a hyphen",
				)
			}
			return nil
		}).Comment("用户名"),
		// 使用 Unique 创建单字段唯一索引
		field.String("email").MaxLen(32).Unique().Comment("邮箱（一个邮箱只能注册一次）"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// 根据 User.author 建立反向边，不额外创建边
		edge.From("articles", Article.Type).Ref("author").Comment("文章列表"),
		// 根据 Group.users 建立反向边，不额外创建关联表
		edge.From("groups", Group.Type).Ref("users").Comment("所在的组列表"),
	}
}

func (User) Config() ent.Config {
	return ent.Config{
		Table: "user",
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		// 添加联合索引，使用 .Unique() 会变成联合唯一索引
		index.Fields("nickname", "created_at"),
		// 单字段唯一索引直接在 Fields() 中进行设置接口
	}
}
