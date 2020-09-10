package test

import (
	"auth/adapter"
	"auth/db"
	dbAccess "auth/db/access"
	"github.com/hashicorp/consul/api"
	"github.com/jinzhu/gorm"
	"log"
	"sync"
)

var (
	manager db.AccessorManage
	dbc *gorm.DB
	waitForFinish sync.WaitGroup
)

const numberOfTestFunc = 15

func init() {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	dbc, _, err = adapter.ConnectDBWithConsul(cli, "db/auth/local_test")
	if err != nil {
		log.Fatal(err)
	}
	db.Migrate(dbc)

	manager, err = db.NewAccessorManage(dbAccess.DefaultReflectType(), dbc)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		waitForFinish = sync.WaitGroup{}
		waitForFinish.Add(numberOfTestFunc)
		waitForFinish.Wait()
		_ = dbc.Close()
	}()
}

var passwords = map[string]string{
	"testPW1": "$2a$10$POwSnghOjkriuQ4w1Bj3zeHIGA7fXv8UI/UFXEhnnO5YrcwkUDcXq",
	"testPW2": "$2a$10$XxGXTboHZxhoqzKcBVqkJOiNSy6narAvIQ/ljfTJ4m93jAt8GyX.e",
	"testPW3": "$2a$10$sfZLOR8iVyhXI0y8nXcKIuKseahKu4NLSlocUWqoBdGrpLIZzxJ2S",
}
