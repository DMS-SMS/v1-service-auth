package access

import (
	"auth/adapter"
	"auth/db"
	"github.com/hashicorp/consul/api"
	"log"
)

var manager db.AccessorManage

func init() {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	dbc, _, err := adapter.ConnectDBWithConsul(cli, "db/auth/local_test")
	if err != nil {
		log.Fatal(err)
	}
	db.Migrate(dbc)
	dbc.LogMode(true)

	manager, err = db.NewAccessorManage(DefaultReflectType(), dbc)
	if err != nil {
		log.Fatal(err)
	}
}

var passwords = map[string]string{
	"testPW1": "$2a$10$POwSnghOjkriuQ4w1Bj3zeHIGA7fXv8UI/UFXEhnnO5YrcwkUDcXq",
	"testPW2": "$2a$10$XxGXTboHZxhoqzKcBVqkJOiNSy6narAvIQ/ljfTJ4m93jAt8GyX.e",
	"testPW3": "$2a$10$sfZLOR8iVyhXI0y8nXcKIuKseahKu4NLSlocUWqoBdGrpLIZzxJ2S",
}

