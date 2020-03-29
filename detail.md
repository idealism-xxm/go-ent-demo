最近发现一个 `Go` 的 schema-first 的 ORM 库 [ent](https://github.com/facebookincubator/ent)，它主打处理图模式的数据，和 `GraphQL` 一样出自 Facebook ，目前还是比较早期的版本 (`v0.1.4`) ，已经实现了 [Roadmap for v1](https://github.com/facebookincubator/ent/issues/46) 中的 40% 的功能。看到另一个 [issue](https://github.com/facebookincubator/ent/issues/357) 提到正在整合 [glqgen](https://github.com/99designs/gqlgen) ，感觉如果使用这个 ORM 的话会和 `gqlgen` 绑定在一起（[go-graphql-comparison](https://github.com/idealism-xxm/go-graphql-comparison) 中实现了 `Go` 中三个 `GraphQL` 库的 Demo 并进行了使用体验上的对比）。

下面使用类似[官网的例子](https://entgo.io/docs/getting-started)体验一下各种功能。

#### 定义 Schema

本样例包含三个实体及相应的关系： 
- Article: 存储文章信息，一篇文章的作者唯一
- User: 存储用户信息，一个用户可以有多篇文章，一个用户可以在多个群组中
- Group: 存储群组信息，一个群组可以有多个用户

接下来需要定义一下各个实体（一般对应数据库的表）的字段等信息 (schema) ，在此之前先介绍一下相关的一些概念：
- fields(or properties): 对应表中的非外键字段，不仅可以定义相关的字段，还能实现应用层的校验限制等逻辑
- edges(or relations): 对应关联关系
    - 关系为 O2O/O2M 时对应表中的外键字段
    - 关系为 M2M 时对应关联表
- indexes: 对应联合（唯一）索引和单字段索引，单字段唯一索引可以直接在 fields 中定义
- config: 实体的设置，目前仅能自定义表的名字
- mixin: 定义一些通用的非外键字段，方便这些字段供其他实体重复使用

首先运行 `entc init Article User Group` 命令，在 `ent/schema` 下面生成三个实体对应的 schema 文件，然后可以就可以按照自己的需求定义相应的信息即可。

为了方便使用，我还在 `ent/schema` 下面新建了 [mixins.go](ent/schema/mixins.go) 文件，用于定义各种 `Mixin` ，样例仅添加 `TimeMixin` 用于把创建时间和修改时间抽成公用代码。

```go
type TimeMixin struct{}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Immutable().Default(time.Now).Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
	}
}
```

接下来就是在各个实体对应的文件中定义各自相关的信息即可：

[ent/schema/article.go](ent/schema/article.go) 中：
- Mixin(): `Article` 仅使用 `TimeMixin` ，会在生成时自动与 Fields 中定义的字段合并
- Fields():  `Article` 仅有两个字段， `.Comment(...)` 用于定义该字段的注释（只在 schema 中，不会在生成的代码和数据库中）, `.StorageKey(...)` 用于自定义数据库中字段名
- Edges(): `Article` 有一条唯一正向边（存储在 article 表中），代表文章的作者
- Config(): `Article` 自定义表名为 `article` ，默认是复数形式
- Indexes(): 由于 `author` 是边中定义的，所以没法和其他字段组成联合索引

```go
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
```

[ent/schema/user.go](ent/schema/user.go) 中：
- Mixin(): `User` 仅使用 `TimeMixin` ，会在生成时自动与 Fields 中定义的字段合并
- Fields():  `User` 有三个字段， `.Validate(...)` 可以在应用层添加自定义的校验逻辑， `.Unique()` 可以指定当前字段有唯一索引
- Edges(): `User` 有两条反向边（反向边不会新建字段和表，会使用对应的正向边），分别代表文章列表和所在组的列表
- Config(): `User` 自定义表名为 `user` ，默认是复数形式
- Indexes(): `User` 有一个联合索引

```go
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
```

[ent/schema/group.go](ent/schema/group.go) 中：
- Fields():  `Group` 有一个字段
- Edges(): `Group` 有一条正向边，代表组拥有的成员列表，会新建一张表存储对应的关系
- Config(): `Group` 自定义表名为 `group` ，默认是复数形式

```go
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
```

#### 生成代码

定义完每个实体的信息后，就需要生成代码了，可以运行 `entc generate ./ent/schema --target ./ent/gen` 将代码生成在 `ent/gen` 文件夹下（需要提前建立文件夹，并新建一个包含 `package gen` 的 `.go` 文件，否则会生成失败）。这里我不太希望生成的代码和自己的代码在一个文件夹下，所以就指定了另一个文件夹专门存放生成的代码，如果不介意可以直接运行 `entc generate ./ent/schema` 即可。为了后期方便，可以在 `ent` 文件夹下建立 [ent.go](/ent/ent.go) ，方便其他使用者生成代码。

```go
package ent

//go:generate go run github.com/facebookincubator/ent/cmd/entc generate ./schema --target ./gen
```

#### Migrate

接下来就是将生成的实体的信息更新到数据库中，直接使用官方样例中的代码即可， `migrate.WithDropIndex(true)` 选项将移除本次信息中没有的索引， `migrate.WithDropColumn(true)` 选项将移除本次信息中没有的字段，目前好像还没有删除没有的关联表的选项。为了后期方便，可以在 `ent` 文件夹下建立 [migrate/main.go](ent/migrate/main.go) ，方便其他使用者生成代码。

```go
func main() {
	client, err := entgen.Open(dialect.MySQL, "root:root@tcp(localhost:3306)/ent")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()

	// run the auto migration tool.
	// 使用 .Debug() 将打印所有的 SQL queries
	err = client.Debug().Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
```

#### CRUD

最后就是 CRUD 操作了：

创建实体好像只能用暴露出来的 api 进行单个创建，没法通过示例直接创建，实际使用起来可能很痛。

```go
func createData(ctx context.Context, client *entgen.Client) {
	// 1. 初始化 Group
	goGroup, _ := client.Group. // Group client
		Create(). // Group create builder
		SetName("Go"). // 设置组名为 Go
		Save(ctx) // 创建并返回
	pythonGroup := client.Group. // Group client
		Create(). // Group create builder
		SetName("Python"). // 设置组名为 Python
		SaveX(ctx) // 创建并返回，如果出错会直接 panic

	// 2. 初始化 User
	idealismUser, _ := client.User. // User client
		Create(). // User create builder
		SetNickname("idealism"). // 设置昵称为 idealism
		SetUsername("idealism-xxm"). // 设置用户名为 idealism-xxm
		SetEmail("idealism-xxm@fake.com").  // 设置邮箱为 idealism-xxm@fake.com
		AddGroups(goGroup).	// 使用 goGroup 实例添加到 goGroup 组中
		AddGroupIDs(pythonGroup.ID). // 使用 pythonGroup 的 id 添加到 pythonGroup 组中
		Save(ctx) // 创建并返回
	anonymityUser, _ := client.User. // User client
		Create(). // User create builder
		SetNickname("anonymity"). // 设置昵称为 idealism
		SetUsername("anonymity"). // 设置用户名为 idealism-xxm
		SetEmail("anonymity@fake.com").  // 设置邮箱为 anonymity@fake.com
		AddGroups(goGroup).	// 使用 goGroup 实例添加到 goGroup 组中
		Save(ctx) // 创建并返回

	// 3. 初始化 Article
	_, _ = client.Article.  // Article client
		Create(). // Article create builder
		SetTitle("Go 中三种 GraphQL 库的简易 Demo 及使用对比"). // 设置文章并标题
		SetContent("Go 中 GraphQL 库 ..."). // 设置文章内容
		SetAuthor(idealismUser). // 使用 idealismUser 实例指定作者
		Save(ctx) // 创建并返回
	_, _ = client.Article.  // Article client
		Create(). // Article create builder
		SetTitle("title for test article"). // 设置文章并标题
		SetContent("content for test article"). // 设置文章内容
		SetAuthorID(anonymityUser.ID). // 使用 anonymityUser 的 id 指定作者
		Save(ctx) // 创建并返回
}
```

查询实体的时候比较方便，提供类型安全的接口，并且支持排序、分页、 `Eager Loading` 和聚合。 

如果查询的字段有 `time.Time` 类型，则需要对链接添加上 `parseTime=true` 参数。

```go
func queryData(ctx context.Context, client *entgen.Client) {
	// 总共 1 条 SQL
	articleTitles, _ := client.Group. // Group client
		Query(). // Group query builder
		Where(group.NameIn("Go", "Python")). // 过滤名称为 Go 或者 Python 的组
		QueryUsers(). // 查询这些组中的用户（联表查询）
		QueryArticles(). // 查询这些用户的文章（联表查询）
		Order(entgen.Desc(article.FieldCreatedAt)). // 按照文章创建时间倒序排序
		Offset(1). // 从第二条开始获取
		Limit(1). // 只获取一条
		Select("title"). // 只查询文章的标题
		Strings(ctx) // 查询所有结果
	for _, articleTitle := range articleTitles {
		fmt.Printf("article.Title = %v\n", articleTitle)
	}

	// 总共 4 条 SQL
	usersWithArticlesAndGroups, err := client.User. // User client
		Query(). // User query builder
		Where(user.Not(user.Username("anonymity"))). // 过滤用户名不为 anonymity 的用户
		WithArticles(func(query *entgen.ArticleQuery) { // 用查询出来的 userIds 再查 article 表
			query.Limit(2) // 只查所有文章的前 2 篇文章（然后再将这 2 篇文章放入对应的用户中）
		}). // 将每个用户的文章都查出来，并放在 Edges.Articles 中
		WithGroups(). // 将每个用户的组都查出来，并放在 Edges.Groups 中（用查询出来的 userIds 先查关联表，然后查 group 表）
		All(ctx)
	for _, userItem := range usersWithArticlesAndGroups {
		for _, articleItem := range userItem.Edges.Articles {
			fmt.Printf("user.Username = %v, article.Title = %v\n", userItem.Username, articleItem.Title)
		}
		for _, groupItem := range userItem.Edges.Groups {
			fmt.Printf("user.Username = %v, group.Name = %v\n", userItem.Username, groupItem.Name)
		}
	}
}
```

更新操作综合了创建与查询的特性，既可以使用 `.Where(...)` 找到待更新的实体集合，也可以使用 `.UpdateOneID(id)` 找到 id 对应的实体，还可以直接在已找到的实体实例上使用 `.Update()` 执行后续的更新操作。更新时的接口与创建时相同，最后使用 `.Save(ctx)` 让修改生效。

删除操作同样如此，，既可以使用 `.Where(...)` 找到待删除的实体集合，也可以使用 `.DeleteOne(...)` 找到实例对应的实体，还可以使用 `.UpdateOneID(id)` 找到 id 对应的实体，然后执行 `.Exec(ctx)` 让删除生效。

如果需要使用事务，则先调用 `client.Tx(ctx)` 获取到 `tx` ，后续的操作都用 `tx` 代替 `client` 即可，最后执行 `tx.Commit()` 进行提交。

同时， `ent` 还支持 [Hooks](https://entgo.io/docs/hooks/) 方便在执行操作前后添加自定义逻辑，既可在 schema 中定义，也能在运行时动态添加（支持动态添加全局的 hooks ，方便支持 traces, metrics, logs, ...）。

#### 总结

使用起来可以发现创建时最痛，暂时没有找到批量创建的方法，其他操作基本能覆盖常用的场景，毕竟这个库还处于早期版本，还在完成核心功能，后续可能会优化使用上的体验。

优点：
- 类似 Django 的 ORM ，支持反向查询
- 类型安全
- 支持 `Eager Loading` 和 `Hooks`

缺点：
- migrate 没有版本概念，无法进行回滚操作
- 非 `unique` 的边删除后，对应的关联表不会进行删除
- 边就算使用的是外键，也无法获取到对应的 id （当然它主打处理图模式的数据，与 GraphQL 相性较好）
- 边就算使用的是外键，也无法作为联合索引中的一个字段
- 创建实体好像只能用暴露出来的 api 单个创建，没法批量创建和通过实例直接创建，实际使用起来比较痛
- `.MaxLen(...)` 限制的长度是字节长度，一个汉字占 3 个字节

以上相关内容可以在 [go-ent-demo](https://github.com/idealism-xxm/go-ent-demo) 中找到。
