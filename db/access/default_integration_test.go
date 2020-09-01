package access

import (
	"auth/adapter"
	"auth/db"
	"auth/model"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
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

func TestDefault_CreateStudentAuth(t *testing.T) {
	// Tx 시작
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}

	// StudentAuth.ParentUUID에 설정할 값을 위한 학부모 계정 생성
	inits := []struct {
		UUID, ParentID, ParentPW string
	} {
		{
			UUID:     "parent-432143214321",
			ParentID: "jinhong07191",
			ParentPW: passwords["testPW1"],
		}, {
			UUID:     "parent-123412341234",
			ParentID: "jinhong07192",
			ParentPW: passwords["testPW2"],
		},
	}

	for _, init := range inits {
		auth := &model.ParentAuth{UUID: init.UUID, ParentId: init.ParentID, ParentPw: init.ParentPW}
		if _, err := access.CreateParentAuth(auth); err != nil {
			access.Rollback()
			log.Fatal(fmt.Sprintf("error occurs while creating parent auth, err: %v", err))
		}
	}

	tests := []struct {
		UUID, ParentUUID, StudentID, StudentPW string
		ExpectAuth *model.StudentAuth
		ExpectError error
	} {
		{ // 정상적인 input test case
			UUID: "student-123412341234",
			ParentUUID: "parent-123412341234",
			StudentID: "jinhong07191",
			StudentPW: passwords["testPW1"],
			ExpectError: nil,
		},
	}

	for _, test := range tests {
		paramAuth := &model.StudentAuth{
			UUID:       test.UUID,
			StudentId:  test.StudentID,
			StudentPw:  test.StudentPW,
			ParentUUID: test.ParentUUID,
		}

		createdAuth, err := access.CreateStudentAuth(paramAuth)
		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, paramAuth, createdAuth, "result model assertion (test case: %v)", test)
	}

	access.Rollback()
}