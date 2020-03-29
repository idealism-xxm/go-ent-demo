package main

import (
	"context"
	"github.com/facebookincubator/ent/dialect"
	entgen "go-ent-demo/ent/gen"
	"go-ent-demo/ent/gen/migrate"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

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
