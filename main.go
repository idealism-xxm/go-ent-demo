package main

import (
	"context"
	"fmt"
	"github.com/facebookincubator/ent/dialect"
	entgen "go-ent-demo/ent/gen"
	"go-ent-demo/ent/gen/article"
	"go-ent-demo/ent/gen/group"
	"go-ent-demo/ent/gen/user"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	client, err := entgen.Open(dialect.MySQL, "root:root@tcp(localhost:3306)/ent?parseTime=true")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	//createData(ctx, client)
	queryData(ctx, client)
}

func createData(ctx context.Context, client *entgen.Client) {
	// 1. 初始化 Group
	goGroup, _ := client.Group. // Group client
					Create().      // Group create builder
					SetName("Go"). // 设置组名为 Go
					Save(ctx)      // 创建并返回
	pythonGroup := client.Group. // Group client
					Create().          // Group create builder
					SetName("Python"). // 设置组名为 Python
					SaveX(ctx)         // 创建并返回，如果出错会直接 panic

	// 2. 初始化 User
	idealismUser, _ := client.User. // User client
					Create().                          // User create builder
					SetNickname("idealism").           // 设置昵称为 idealism
					SetUsername("idealism-xxm").       // 设置用户名为 idealism-xxm
					SetEmail("idealism-xxm@fake.com"). // 设置邮箱为 idealism-xxm@fake.com
					AddGroups(goGroup).                // 使用 goGroup 实例添加到 goGroup 组中
					AddGroupIDs(pythonGroup.ID).       // 使用 pythonGroup 的 id 添加到 pythonGroup 组中
					Save(ctx)                          // 创建并返回
	anonymityUser, _ := client.User. // User client
						Create().                       // User create builder
						SetNickname("anonymity").       // 设置昵称为 idealism
						SetUsername("anonymity").       // 设置用户名为 idealism-xxm
						SetEmail("anonymity@fake.com"). // 设置邮箱为 anonymity@fake.com
						AddGroups(goGroup).             // 使用 goGroup 实例添加到 goGroup 组中
						Save(ctx)                       // 创建并返回

	// 3. 初始化 Article
	_, _ = client.Article. // Article client
				Create().                         // Article create builder
				SetTitle("Go 中三种 GraphQL 库的简").   // 设置文章并标题
				SetContent("Go 中 GraphQL 库 ..."). // 设置文章内容
				SetAuthor(idealismUser).          // 使用 idealismUser 实例指定作者
				Save(ctx)                         // 创建并返回
	_, _ = client.Article. // Article client
				Create().                               // Article create builder
				SetTitle("title for test article").     // 设置文章并标题
				SetContent("content for test article"). // 设置文章内容
				SetAuthorID(anonymityUser.ID).          // 使用 anonymityUser 的 id 指定作者
				Save(ctx)                               // 创建并返回
}

func queryData(ctx context.Context, client *entgen.Client) {
	// 总共 1 条 SQL
	articleTitles, _ := client.Group. // Group client
						Query().                                    // Group query builder
						Where(group.NameIn("Go", "Python")).        // 过滤名称为 Go 或者 Python 的组
						QueryUsers().                               // 查询这些组中的用户（联表查询）
						QueryArticles().                            // 查询这些用户的文章（联表查询）
						Order(entgen.Desc(article.FieldCreatedAt)). // 按照文章创建时间倒序排序
						Offset(1).                                  // 从第二条开始获取
						Limit(1).                                   // 只获取一条
						Select("title").                            // 只查询文章的标题
						Strings(ctx)                                // 查询所有结果
	for _, articleTitle := range articleTitles {
		fmt.Printf("article.Title = %v\n", articleTitle)
	}

	// 总共 4 条 SQL
	usersWithArticlesAndGroups, _ := client.User. // User client
							Query().                                        // User query builder
							Where(user.Not(user.Username("anonymity"))).    // 过滤用户名不为 anonymity 的用户
							WithArticles(func(query *entgen.ArticleQuery) { // 用查询出来的 userIds 再查 article 表
			query.Limit(2) // 只查所有文章的前 2 篇文章（然后再将这 2 篇文章放入对应的用户中）
		}).           // 将每个用户的文章都查出来，并放在 Edges.Articles 中
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
