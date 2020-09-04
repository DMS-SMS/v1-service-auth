package access

import (
	"auth/adapter"
	"auth/db"
	"auth/model"
	"auth/tool/mysqlerr"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var manager db.AccessorManage
var dbc *gorm.DB

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

	manager, err = db.NewAccessorManage(DefaultReflectType(), dbc)
	if err != nil {
		log.Fatal(err)
	}
}

func Test_default_CreateStudentAuth(t *testing.T) {
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
			UUID:     "parent-111111111111",
			ParentID: "jinhong07191",
			ParentPW: passwords["testPW1"],
		}, {
			UUID:     "parent-222222222222",
			ParentID: "jinhong07192",
			ParentPW: passwords["testPW2"],
		},
	}

	for _, init := range inits {
		auth := &model.ParentAuth{
			UUID:     model.UUID(init.UUID),
			ParentID: model.ParentID(init.ParentID),
			ParentPW: model.ParentPW(init.ParentPW),
		}
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
		{ // success case
			UUID: "student-111111111111",
			ParentUUID: "parent-111111111111",
			StudentID: "jinhong0719",
			StudentPW: passwords["testPW1"],
			ExpectError: nil,
		}, { // UUID duplicate
			UUID: "student-111111111111",
			ParentUUID: "parent-111111111111",
			StudentID: "jinhong07191",
			StudentPW: passwords["testPW1"],
			ExpectError: mysqlerr.DuplicateEntry("uuid", "student-111111111111"),
		}, { // StudentID duplicate
			UUID: "student-222222222222",
			ParentUUID: "parent-111111111111",
			StudentID: "jinhong0719",
			StudentPW: passwords["testPW1"],
			ExpectError: mysqlerr.DuplicateEntry("student_id", "jinhong0719"),
		}, { // ParentUUID(foreign key) reference constraint X
			UUID: "student-222222222222",
			ParentUUID: "parent-123412341234", // not exist parent uuid
			StudentID: "jinhong07192",
			StudentPW: passwords["testPW1"],
			ExpectError: studentAuthParentUUIDFKConstraintFailError,
		},
	}

	for _, test := range tests {
		auth := &model.StudentAuth{
			UUID:       model.UUID(test.UUID),
			StudentID:  model.StudentID(test.StudentID),
			StudentPW:  model.StudentPW(test.StudentPW),
			ParentUUID: model.ParentUUID(test.ParentUUID),
		}

		test.ExpectAuth = auth.DeepCopy()
		auth, err := access.CreateStudentAuth(auth)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectAuth, auth.ExceptGormModel(), "result model assertion (test case: %v)", test)
	}

	access.Rollback()
}

func Test_default_CreateParentAuth(t *testing.T) {
	tests := []struct{
		UUID, ParentId, ParentPw string
		ExpectAuth *model.ParentAuth
		ExpectError error
	} {
		{ // success case
			UUID: "parent-111111111111",
			ParentId: "parent1",
			ParentPw: passwords["testPW1"],
			ExpectError: nil,
		}, { // UUID duplicate
			UUID: "parent-111111111111",
			ParentId: "parent2",
			ParentPw: passwords["testPW2"],
			ExpectError: mysqlerr.DuplicateEntry("uuid", "parent-111111111111"),
		}, { // ParentId duplicate
			UUID: "parent-222222222222",
			ParentId: "parent1",
			ParentPw: passwords["testPW2"],
			ExpectError: mysqlerr.DuplicateEntry("parent_id", "parent1"),
		},
	}

	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}

	for _, test := range tests {
		auth := &model.ParentAuth{
			UUID:     model.UUID(test.UUID),
			ParentID: model.ParentID(test.ParentId),
			ParentPW: model.ParentPW(test.ParentPw),
		}

		test.ExpectAuth = auth.DeepCopy()
		auth, err := access.CreateParentAuth(auth)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectAuth, auth.ExceptGormModel(), "result model assertion (test case: %v)", test)
	}

	access.Rollback()
}

func TestDBClose(t *testing.T) {
	_ = dbc.Close()
}

var passwords = map[string]string{
	"testPW1": "$2a$10$POwSnghOjkriuQ4w1Bj3zeHIGA7fXv8UI/UFXEhnnO5YrcwkUDcXq",
	"testPW2": "$2a$10$XxGXTboHZxhoqzKcBVqkJOiNSy6narAvIQ/ljfTJ4m93jAt8GyX.e",
	"testPW3": "$2a$10$sfZLOR8iVyhXI0y8nXcKIuKseahKu4NLSlocUWqoBdGrpLIZzxJ2S",
}

var (
	studentAuthParentUUIDFKConstraintFailError = mysqlerr.FKConstraintFail("sms_auth_test_db",
		"student_auths", (&model.StudentAuth{}).ParentUUIDConstraintName(), "parent_uuid",
		mysqlerr.Reference{
			TableName: "parent_auths",
			AttrName:  "uuid",
		})
)