package test

import (
	"auth/db"
	"auth/db/access"
	"github.com/hashicorp/consul/api"
	"log"
	"sync"
)

func init() {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	dbc, _, err = db.ConnectWithConsul(cli, "db/auth/local_test")
	if err != nil {
		log.Fatal(err)
	}
	db.Migrate(dbc)

	manager, err = db.NewAccessorManage(access.Default(dbc))
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